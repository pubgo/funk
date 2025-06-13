package bbolt

import (
	"context"
	"path/filepath"

	bolt "go.etcd.io/bbolt"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/merge"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/strutil"
	"github.com/pubgo/funk/v2/result"
)

func New(cfg *Config, log log.Logger) *Client {
	cfg = merge.Copy(DefaultConfig(), cfg).Unwrap()

	path := filepath.Join(config.GetConfigDir(), cfg.Path)
	assert.Must(pathutil.IsNotExistMkDir(filepath.Dir(path)))
	db := assert.Must1(bolt.Open(path, cfg.FileMode, cfg.getOpts()))

	return &Client{DB: db, log: log}
}

type Client struct {
	*bolt.DB
	log log.Logger
}

func (t *Client) bucket(name string, tx *bolt.Tx) (*bolt.Bucket, error) {
	_, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, errors.WrapKV(err, "bucket_name", name)
	}

	return tx.Bucket([]byte(name)), nil
}

func (t *Client) Set(ctx context.Context, key string, val []byte, names ...string) error {
	return t.Update(ctx, func(bucket *bolt.Bucket) error {
		return bucket.Put([]byte(key), val)
	}, names...)
}

func (t *Client) Get(ctx context.Context, key string, names ...string) result.Result[[]byte] {
	var (
		val []byte
		err = t.View(ctx, func(bucket *bolt.Bucket) error {
			val = bucket.Get([]byte(key))
			return nil
		}, names...)
	)

	return result.Wrap(val, err)
}

func (t *Client) List(ctx context.Context, fn func(k, v []byte) error, names ...string) error {
	return t.View(ctx, func(bucket *bolt.Bucket) error { return bucket.ForEach(fn) }, names...)
}

func (t *Client) Delete(ctx context.Context, key string, names ...string) error {
	return t.Update(ctx, func(bucket *bolt.Bucket) error {
		return bucket.Delete([]byte(key))
	}, names...)
}

func (t *Client) View(ctx context.Context, fn func(*bolt.Bucket) error, names ...string) error {
	name := strutil.GetDefault(names...)

	return t.DB.View(func(tx *bolt.Tx) error {
		bucket, err := t.bucket(name, tx)
		if err != nil {
			return err
		}
		return fn(bucket)
	})
}

func (t *Client) Update(ctx context.Context, fn func(*bolt.Bucket) error, names ...string) error {
	name := strutil.GetDefault(names...)

	return t.DB.Update(func(tx *bolt.Tx) (err error) {
		bucket, err := t.bucket(name, tx)
		if err != nil {
			return err
		}
		return fn(bucket)
	})
}
