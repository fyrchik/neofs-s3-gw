package tokens

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neofs-api-go/pkg/client"
	cid "github.com/nspcc-dev/neofs-api-go/pkg/container/id"
	"github.com/nspcc-dev/neofs-api-go/pkg/object"
	"github.com/nspcc-dev/neofs-api-go/pkg/owner"
	"github.com/nspcc-dev/neofs-s3-gw/creds/accessbox"
	"github.com/nspcc-dev/neofs-sdk-go/pool"
)

type (
	// Credentials is a bearer token get/put interface.
	Credentials interface {
		GetBox(context.Context, *object.Address) (*accessbox.Box, error)
		Put(context.Context, *cid.ID, *owner.ID, *accessbox.AccessBox, ...*keys.PublicKey) (*object.Address, error)
	}

	cred struct {
		key  *keys.PrivateKey
		pool pool.Pool
	}
)

var (
	// ErrEmptyPublicKeys is returned when no HCS keys are provided.
	ErrEmptyPublicKeys = errors.New("HCS public keys could not be empty")
	// ErrEmptyBearerToken is returned when no bearer token is provided.
	ErrEmptyBearerToken = errors.New("Bearer token could not be empty")
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var _ = New

// New creates new Credentials instance using given cli and key.
func New(conns pool.Pool, key *keys.PrivateKey) Credentials {
	return &cred{pool: conns, key: key}
}

func (c *cred) acquireBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func (c *cred) releaseBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

func (c *cred) GetTokens(ctx context.Context, address *object.Address) (*accessbox.GateData, error) {
	box, err := c.getAccessBox(ctx, address)
	if err != nil {
		return nil, err
	}

	return box.GetTokens(c.key)
}

func (c *cred) GetBox(ctx context.Context, address *object.Address) (*accessbox.Box, error) {
	box, err := c.getAccessBox(ctx, address)
	if err != nil {
		return nil, err
	}

	return box.GetBox(c.key)
}

func (c *cred) getAccessBox(ctx context.Context, address *object.Address) (*accessbox.AccessBox, error) {
	var (
		box accessbox.AccessBox
		buf = c.acquireBuffer()
	)
	defer c.releaseBuffer(buf)

	ops := new(client.GetObjectParams).WithAddress(address).WithPayloadWriter(buf)

	_, err := c.pool.GetObject(
		ctx,
		ops,
	)
	if err != nil {
		return nil, err
	}

	if err = box.Unmarshal(buf.Bytes()); err != nil {
		return nil, err
	}
	return &box, nil
}

func (c *cred) Put(ctx context.Context, cid *cid.ID, issuer *owner.ID, box *accessbox.AccessBox, keys ...*keys.PublicKey) (*object.Address, error) {
	var (
		err     error
		created = strconv.FormatInt(time.Now().Unix(), 10)
	)

	if len(keys) == 0 {
		return nil, ErrEmptyPublicKeys
	} else if box == nil {
		return nil, ErrEmptyBearerToken
	}
	data, err := box.Marshal()
	if err != nil {
		return nil, err
	}

	timestamp := object.NewAttribute()
	timestamp.SetKey(object.AttributeTimestamp)
	timestamp.SetValue(created)

	filename := object.NewAttribute()
	filename.SetKey(object.AttributeFileName)
	filename.SetValue(created + "_access.box")

	raw := object.NewRaw()
	raw.SetContainerID(cid)
	raw.SetOwnerID(issuer)
	raw.SetAttributes(filename, timestamp)

	ops := new(client.PutObjectParams).WithObject(raw.Object()).WithPayloadReader(bytes.NewBuffer(data))
	oid, err := c.pool.PutObject(
		ctx,
		ops,
	)
	if err != nil {
		return nil, err
	}
	address := object.NewAddress()
	address.SetObjectID(oid)
	address.SetContainerID(cid)
	return address, nil
}
