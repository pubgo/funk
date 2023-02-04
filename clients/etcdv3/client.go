package etcdv3

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/merge"
	client3 "go.etcd.io/etcd/client/v3"
)

func New(cfg *Config) *Client {
	assert.If(cfg == nil, "cfg is nil")

	var cc = merge.Copy(DefaultCfg(), cfg).Unwrap()
	assert.Must(cc.Build())
	return &Client{Client: cc.Get()}
}

type Client struct {
	*client3.Client
}
