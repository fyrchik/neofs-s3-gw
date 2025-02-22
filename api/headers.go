package api

// Standard S3 HTTP request/response constants.
const (
	MetadataPrefix       = "X-Amz-Meta-"
	AmzMetadataDirective = "X-Amz-Metadata-Directive"
	AmzVersionID         = "X-Amz-Version-Id"
	AmzTaggingCount      = "X-Amz-Tagging-Count"
	AmzTagging           = "X-Amz-Tagging"

	LastModified       = "Last-Modified"
	Date               = "Date"
	ETag               = "ETag"
	ContentType        = "Content-Type"
	ContentMD5         = "Content-Md5"
	ContentEncoding    = "Content-Encoding"
	Expires            = "Expires"
	ContentLength      = "Content-Length"
	ContentLanguage    = "Content-Language"
	ContentRange       = "Content-Range"
	Connection         = "Connection"
	AcceptRanges       = "Accept-Ranges"
	AmzBucketRegion    = "X-Amz-Bucket-Region"
	ServerInfo         = "Server"
	RetryAfter         = "Retry-After"
	Location           = "Location"
	CacheControl       = "Cache-Control"
	ContentDisposition = "Content-Disposition"
	Authorization      = "Authorization"
	Action             = "Action"
	IfModifiedSince    = "If-Modified-Since"
	IfUnmodifiedSince  = "If-Unmodified-Since"
	IfMatch            = "If-Match"
	IfNoneMatch        = "If-None-Match"

	AmzCopyIfModifiedSince       = "X-Amz-Copy-Source-If-Modified-Since"
	AmzCopyIfUnmodifiedSince     = "X-Amz-Copy-Source-If-Unmodified-Since"
	AmzCopyIfMatch               = "X-Amz-Copy-Source-If-Match"
	AmzCopyIfNoneMatch           = "X-Amz-Copy-Source-If-None-Match"
	AmzACL                       = "X-Amz-Acl"
	AmzGrantFullControl          = "X-Amz-Grant-Full-Control"
	AmzGrantRead                 = "X-Amz-Grant-Read"
	AmzGrantWrite                = "X-Amz-Grant-Write"
	AmzExpectedBucketOwner       = "X-Amz-Expected-Bucket-Owner"
	AmzSourceExpectedBucketOwner = "X-Amz-Source-Expected-Bucket-Owner"

	ContainerID = "X-Container-Id"
)

// S3 request query params.
const (
	QueryVersionID = "versionId"
)
