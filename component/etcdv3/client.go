package etcdv3

import (
	client3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/config"
	"github.com/pubgo/funk/v2/merge"
	"github.com/pubgo/funk/v2/retry"
)

func New(conf *Config) *Client {
	conf = config.MergeR(DefaultCfg(), *conf).Unwrap()
	cfg := merge.Struct(new(client3.Config), conf).Unwrap()
	cfg.DialOptions = append(
		cfg.DialOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	// 创建etcd client对象
	return &Client{Client: assert.Must1(retry.Default().DoVal(func(i int) (any, error) {
		return client3.New(*cfg)
	})).(*client3.Client)}
}

type Client struct {
	*client3.Client
}
