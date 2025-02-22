package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nspcc-dev/neofs-api-go/pkg/acl/eacl"
	"github.com/nspcc-dev/neofs-s3-gw/api"
	"github.com/nspcc-dev/neofs-s3-gw/api/auth"
	"github.com/nspcc-dev/neofs-s3-gw/api/errors"
	"github.com/nspcc-dev/neofs-s3-gw/api/layer"
	"go.uber.org/zap"
)

type postPolicy struct {
	Expiration time.Time          `json:"expiration"`
	Conditions []*policyCondition `json:"conditions"`
	empty      bool
}

func (p *postPolicy) condition(key string) *policyCondition {
	for _, condition := range p.Conditions {
		if condition.Key == key {
			return condition
		}
	}
	return nil
}

func (p *postPolicy) CheckContentLength(size int64) bool {
	if p.empty {
		return true
	}
	for _, condition := range p.Conditions {
		if condition.Matching == "content-length-range" {
			length := strconv.FormatInt(size, 10)
			return condition.Key <= length && length <= condition.Value
		}
	}
	return true
}

func (p *policyCondition) match(value string) bool {
	switch p.Matching {
	case "eq":
		p.Matched = p.Value == value
	case "starts-with":
		if p.Key == api.ContentType {
			p.Matched = true
			for _, contentType := range strings.Split(value, ",") {
				if !strings.HasPrefix(contentType, p.Value) {
					p.Matched = false
				}
			}
		} else {
			p.Matched = strings.HasPrefix(value, p.Value)
		}
	}
	return p.Matched
}

func (p *postPolicy) CheckField(key string, value string) error {
	if p.empty {
		return nil
	}
	cond := p.condition(key)
	if cond == nil {
		return errors.GetAPIError(errors.ErrPostPolicyConditionInvalidFormat)
	}

	if !cond.match(value) {
		return errors.GetAPIError(errors.ErrPostPolicyConditionInvalidFormat)
	}

	return nil
}

func (p *postPolicy) AllConditionMatched() bool {
	for _, condition := range p.Conditions {
		if !condition.Matched {
			return false
		}
	}
	return true
}

type policyCondition struct {
	Matching string
	Key      string
	Value    string
	Matched  bool
}

var errInvalidCondition = fmt.Errorf("invalid condition")

func (p *policyCondition) UnmarshalJSON(data []byte) error {
	var (
		ok bool
		v  interface{}
	)

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v := v.(type) {
	case []interface{}:
		if len(v) != 3 {
			return errInvalidCondition
		}
		if p.Matching, ok = v[0].(string); !ok {
			return errInvalidCondition
		}

		if p.Matching == "content-length-range" {
			min, ok := v[1].(float64)
			max, ok2 := v[2].(float64)
			if !ok || !ok2 {
				return errInvalidCondition
			}
			p.Key = strconv.FormatFloat(min, 'f', 0, 32)
			p.Value = strconv.FormatFloat(max, 'f', 0, 32)
		} else {
			key, ok2 := v[1].(string)
			p.Value, ok = v[2].(string)
			if !ok || !ok2 {
				return errInvalidCondition
			}
			p.Key = strings.ToLower(strings.TrimPrefix(key, "$"))
		}

	case map[string]interface{}:
		p.Matching = "eq"
		for key, val := range v {
			p.Key = strings.ToLower(key)
			if p.Value, ok = val.(string); !ok {
				return errInvalidCondition
			}
		}
	default:
		return fmt.Errorf("unknown condition type")
	}

	return nil
}

// keywords of predefined basic ACL values.
const (
	basicACLPrivate   = "private"
	basicACLReadOnly  = "public-read"
	basicACLPublic    = "public-read-write"
	cannedACLAuthRead = "authenticated-read"

	publicBasicRule = 0x0FFFFFFF
)

type createBucketParams struct {
	XMLName            xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CreateBucketConfiguration" json:"-"`
	LocationConstraint string
}

func (h *handler) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	var newEaclTable *eacl.Table
	reqInfo := api.GetReqInfo(r.Context())
	tagSet, err := parseTaggingHeader(r.Header)
	if err != nil {
		h.logAndSendError(w, "could not parse tagging header", reqInfo, err)
		return
	}

	bktInfo, err := h.obj.GetBucketInfo(r.Context(), reqInfo.BucketName)
	if err != nil {
		h.logAndSendError(w, "could not get bucket eacl", reqInfo, err)
		return
	}
	if err = checkOwner(bktInfo, r.Header.Get(api.AmzExpectedBucketOwner)); err != nil {
		h.logAndSendError(w, "expected owner doesn't match", reqInfo, err)
		return
	}

	metadata := parseMetadata(r)
	if contentType := r.Header.Get(api.ContentType); len(contentType) > 0 {
		metadata[api.ContentType] = contentType
	}

	params := &layer.PutObjectParams{
		Bucket: reqInfo.BucketName,
		Object: reqInfo.ObjectName,
		Reader: r.Body,
		Size:   r.ContentLength,
		Header: metadata,
	}

	info, err := h.obj.PutObject(r.Context(), params)
	if err != nil {
		h.logAndSendError(w, "could not upload object", reqInfo, err)
		return
	}

	if containsACLHeaders(r) {
		if newEaclTable, err = h.getNewEAclTable(r, info); err != nil {
			h.logAndSendError(w, "could not get new eacl table", reqInfo, err)
			return
		}
	}

	if tagSet != nil {
		if err = h.obj.PutObjectTagging(r.Context(), &layer.PutTaggingParams{ObjectInfo: info, TagSet: tagSet}); err != nil {
			h.logAndSendError(w, "could not upload object tagging", reqInfo, err)
			return
		}
	}

	if newEaclTable != nil {
		p := &layer.PutBucketACLParams{
			Name: reqInfo.BucketName,
			EACL: newEaclTable,
		}

		if err = h.obj.PutBucketACL(r.Context(), p); err != nil {
			h.logAndSendError(w, "could not put bucket acl", reqInfo, err)
			return
		}
	}

	if versioning, err := h.obj.GetBucketVersioning(r.Context(), reqInfo.BucketName); err != nil {
		h.log.Warn("couldn't get bucket versioning", zap.String("bucket name", reqInfo.BucketName), zap.Error(err))
	} else if versioning.VersioningEnabled {
		w.Header().Set(api.AmzVersionID, info.Version())
	}

	w.Header().Set(api.ETag, info.HashSum)
	api.WriteSuccessResponseHeadersOnly(w)
}

func (h *handler) PostObject(w http.ResponseWriter, r *http.Request) {
	var (
		newEaclTable *eacl.Table
		tagSet       map[string]string
		reqInfo      = api.GetReqInfo(r.Context())
		metadata     = make(map[string]string)
	)

	policy, err := checkPostPolicy(r, reqInfo, metadata)
	if err != nil {
		h.logAndSendError(w, "failed check policy", reqInfo, err)
		return
	}

	if tagging := auth.MultipartFormValue(r, "tagging"); tagging != "" {
		buffer := bytes.NewBufferString(tagging)
		tagSet, err = readTagSet(buffer)
		if err != nil {
			h.logAndSendError(w, "could not read tag set", reqInfo, err)
			return
		}
	}

	var contentReader io.Reader
	var size int64
	if content, ok := r.MultipartForm.Value["file"]; ok {
		contentReader = bytes.NewBufferString(content[0])
		size = int64(len(content[0]))
	} else {
		file, head, err := r.FormFile("file")
		if err != nil {
			h.logAndSendError(w, "could get uploading file", reqInfo, err)
			return
		}
		contentReader = file
		size = head.Size
		reqInfo.ObjectName = strings.ReplaceAll(reqInfo.ObjectName, "${filename}", head.Filename)
	}
	if !policy.CheckContentLength(size) {
		h.logAndSendError(w, "invalid content-length", reqInfo, errors.GetAPIError(errors.ErrInvalidArgument))
		return
	}

	params := &layer.PutObjectParams{
		Bucket: reqInfo.BucketName,
		Object: reqInfo.ObjectName,
		Reader: contentReader,
		Size:   size,
		Header: metadata,
	}

	info, err := h.obj.PutObject(r.Context(), params)
	if err != nil {
		h.logAndSendError(w, "could not upload object", reqInfo, err)
		return
	}

	if acl := auth.MultipartFormValue(r, "acl"); acl != "" {
		r.Header.Set(api.AmzACL, acl)
		r.Header.Set(api.AmzGrantFullControl, "")
		r.Header.Set(api.AmzGrantWrite, "")
		r.Header.Set(api.AmzGrantRead, "")

		if newEaclTable, err = h.getNewEAclTable(r, info); err != nil {
			h.logAndSendError(w, "could not get new eacl table", reqInfo, err)
			return
		}
	}

	if tagSet != nil {
		if err = h.obj.PutObjectTagging(r.Context(), &layer.PutTaggingParams{ObjectInfo: info, TagSet: tagSet}); err != nil {
			h.logAndSendError(w, "could not upload object tagging", reqInfo, err)
			return
		}
	}

	if newEaclTable != nil {
		p := &layer.PutBucketACLParams{
			Name: reqInfo.BucketName,
			EACL: newEaclTable,
		}

		if err = h.obj.PutBucketACL(r.Context(), p); err != nil {
			h.logAndSendError(w, "could not put bucket acl", reqInfo, err)
			return
		}
	}

	if versioning, err := h.obj.GetBucketVersioning(r.Context(), reqInfo.BucketName); err != nil {
		h.log.Warn("couldn't get bucket versioning", zap.String("bucket name", reqInfo.BucketName), zap.Error(err))
	} else if versioning.VersioningEnabled {
		w.Header().Set(api.AmzVersionID, info.Version())
	}

	if redirectURL := auth.MultipartFormValue(r, "success_action_redirect"); redirectURL != "" {
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}
	status := http.StatusNoContent
	if statusStr := auth.MultipartFormValue(r, "success_action_status"); statusStr != "" {
		switch statusStr {
		case "200":
			status = http.StatusOK
		case "201":
			status = http.StatusCreated
			resp := &PostResponse{
				Bucket: info.Bucket,
				Key:    info.Name,
				ETag:   info.HashSum,
			}
			w.WriteHeader(status)
			if _, err = w.Write(api.EncodeResponse(resp)); err != nil {
				h.logAndSendError(w, "something went wrong", reqInfo, err)
			}
			return
		}
	}

	w.Header().Set(api.ETag, info.HashSum)
	w.WriteHeader(status)
}

func checkPostPolicy(r *http.Request, reqInfo *api.ReqInfo, metadata map[string]string) (*postPolicy, error) {
	policy := &postPolicy{}
	if policyStr := auth.MultipartFormValue(r, "policy"); policyStr != "" {
		policyData, err := base64.StdEncoding.DecodeString(policyStr)
		if err != nil {
			return nil, fmt.Errorf("could not decode policy: %w", err)
		}
		if err = json.Unmarshal(policyData, policy); err != nil {
			return nil, fmt.Errorf("could not unmarshal policy: %w", err)
		}
		if policy.Expiration.Before(time.Now()) {
			return nil, fmt.Errorf("policy is expired: %w", errors.GetAPIError(errors.ErrInvalidArgument))
		}
	}

	for key, v := range r.MultipartForm.Value {
		value := v[0]
		if key == "file" || key == "policy" || key == "x-amz-signature" || strings.HasPrefix(key, "x-ignore-") {
			continue
		}
		if err := policy.CheckField(key, value); err != nil {
			return nil, fmt.Errorf("'%s' form field doesn't match the policy: %w", key, err)
		}

		prefix := strings.ToLower(api.MetadataPrefix)
		if strings.HasPrefix(key, prefix) {
			metadata[strings.TrimPrefix(key, prefix)] = value
		}

		if key == "content-type" {
			metadata[api.ContentType] = value
		}

		if key == "key" {
			reqInfo.ObjectName = value
		}
	}

	for _, cond := range policy.Conditions {
		if cond.Key == "bucket" {
			if !cond.match(reqInfo.BucketName) {
				return nil, errors.GetAPIError(errors.ErrPostPolicyConditionInvalidFormat)
			}
		}
	}

	return policy, nil
}

func containsACLHeaders(r *http.Request) bool {
	return r.Header.Get(api.AmzACL) != "" || r.Header.Get(api.AmzGrantRead) != "" ||
		r.Header.Get(api.AmzGrantFullControl) != "" || r.Header.Get(api.AmzGrantWrite) != ""
}

func (h *handler) getNewEAclTable(r *http.Request, objInfo *layer.ObjectInfo) (*eacl.Table, error) {
	var newEaclTable *eacl.Table
	objectACL, err := parseACLHeaders(r)
	if err != nil {
		return nil, fmt.Errorf("could not parse object acl: %w", err)
	}

	resInfo := &resourceInfo{
		Bucket:  objInfo.Bucket,
		Object:  objInfo.Name,
		Version: objInfo.Version(),
	}

	bktPolicy, err := aclToPolicy(objectACL, resInfo)
	if err != nil {
		return nil, fmt.Errorf("could not translate object acl to bucket policy: %w", err)
	}

	astChild, err := policyToAst(bktPolicy)
	if err != nil {
		return nil, fmt.Errorf("could not translate policy to ast: %w", err)
	}

	bacl, err := h.obj.GetBucketACL(r.Context(), objInfo.Bucket)
	if err != nil {
		return nil, fmt.Errorf("could not get bucket eacl: %w", err)
	}

	parentAst := tableToAst(bacl.EACL, objInfo.Bucket)
	for _, resource := range parentAst.Resources {
		if resource.Bucket == bacl.Info.CID.String() {
			resource.Bucket = objInfo.Bucket
		}
	}

	if resAst, updated := mergeAst(parentAst, astChild); updated {
		if newEaclTable, err = astToTable(resAst); err != nil {
			return nil, fmt.Errorf("could not translate ast to table: %w", err)
		}
	}

	return newEaclTable, nil
}

func parseTaggingHeader(header http.Header) (map[string]string, error) {
	var tagSet map[string]string
	if tagging := header.Get(api.AmzTagging); len(tagging) > 0 {
		queries, err := url.ParseQuery(tagging)
		if err != nil {
			return nil, errors.GetAPIError(errors.ErrInvalidArgument)
		}
		if len(queries) > maxTags {
			return nil, errors.GetAPIError(errors.ErrInvalidTagsSizeExceed)
		}
		tagSet = make(map[string]string, len(queries))
		for k, v := range queries {
			tag := Tag{Key: k, Value: v[0]}
			if err = checkTag(tag); err != nil {
				return nil, err
			}
			tagSet[tag.Key] = tag.Value
		}
	}
	return tagSet, nil
}

func parseMetadata(r *http.Request) map[string]string {
	res := make(map[string]string)
	for k, v := range r.Header {
		if strings.HasPrefix(k, api.MetadataPrefix) {
			key := strings.ToLower(strings.TrimPrefix(k, api.MetadataPrefix))
			res[key] = v[0]
		}
	}
	return res
}

func (h *handler) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var (
		reqInfo = api.GetReqInfo(r.Context())
		p       = layer.CreateBucketParams{Name: reqInfo.BucketName, ACL: publicBasicRule}
	)

	if err := checkBucketName(reqInfo.BucketName); err != nil {
		h.logAndSendError(w, "invalid bucket name", reqInfo, err)
		return
	}

	bktACL, err := parseACLHeaders(r)
	if err != nil {
		h.logAndSendError(w, "could not parse bucket acl", reqInfo, err)
		return
	}
	resInfo := &resourceInfo{Bucket: reqInfo.BucketName}

	p.EACL, err = bucketACLToTable(bktACL, resInfo)
	if err != nil {
		h.logAndSendError(w, "could translate bucket acl to eacl", reqInfo, err)
		return
	}

	createParams, err := parseLocationConstraint(r)
	if err != nil {
		h.logAndSendError(w, "could not parse body", reqInfo, err)
		return
	}

	p.BoxData, err = layer.GetBoxData(r.Context())
	if err != nil {
		h.logAndSendError(w, "could not get boxData", reqInfo, err)
		return
	}

	if createParams.LocationConstraint != "" {
		for _, placementPolicy := range p.BoxData.Policies {
			if placementPolicy.LocationConstraint == createParams.LocationConstraint {
				p.Policy = placementPolicy.Policy
				break
			}
		}
	}
	if p.Policy == nil {
		p.Policy = h.cfg.DefaultPolicy
	}

	cid, err := h.obj.CreateBucket(r.Context(), &p)
	if err != nil {
		h.logAndSendError(w, "could not create bucket", reqInfo, err)
		return
	}

	h.log.Info("bucket is created",
		zap.String("container_id", cid.String()))

	api.WriteSuccessResponseHeadersOnly(w)
}

func checkBucketName(bucketName string) error {
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return errors.GetAPIError(errors.ErrInvalidBucketName)
	}

	if strings.HasPrefix(bucketName, "xn--") || strings.HasSuffix(bucketName, "-s3alias") {
		return errors.GetAPIError(errors.ErrInvalidBucketName)
	}
	if net.ParseIP(bucketName) != nil {
		return errors.GetAPIError(errors.ErrInvalidBucketName)
	}

	labels := strings.Split(bucketName, ".")
	for _, label := range labels {
		if len(label) == 0 {
			return errors.GetAPIError(errors.ErrInvalidBucketName)
		}
		for i, r := range label {
			if !isAlphaNum(r) && r != '-' {
				return errors.GetAPIError(errors.ErrInvalidBucketName)
			}
			if (i == 0 || i == len(label)-1) && r == '-' {
				return errors.GetAPIError(errors.ErrInvalidBucketName)
			}
		}
	}

	return nil
}

func isAlphaNum(char int32) bool {
	return 'a' <= char && char <= 'z' || '0' <= char && char <= '9'
}

func parseLocationConstraint(r *http.Request) (*createBucketParams, error) {
	if r.ContentLength == 0 {
		return new(createBucketParams), nil
	}

	params := new(createBucketParams)
	if err := xml.NewDecoder(r.Body).Decode(params); err != nil {
		return nil, err
	}
	return params, nil
}
