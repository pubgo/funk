package etcdv3

import (
	client3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/merge"
	"github.com/pubgo/funk/retry"
)

func New(conf *Config) *Client {
	conf = config.MergeR(DefaultCfg(), *conf).Unwrap()
	cfg := merge.Struct(new(client3.Config), conf).Unwrap()
	cfg.DialOptions = append(
		cfg.DialOptions,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	// 创建etcd client对象
	return &Client{Client: assert.Must1(retry.Default().DoVal(func(i int) (interface{}, error) {
		return client3.New(*cfg)
	})).(*client3.Client)}
}

type Client struct {
	*client3.Client
}
