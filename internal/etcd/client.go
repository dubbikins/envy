package etcd

import clientv3 "go.etcd.io/etcd/client/v3"


type client struct {
	client *clientv3.Client
	cfg *clientv3.Config
}


func New(cfg *clientv3.Config) (*client, error) {
	c, err := clientv3.New(*cfg)
	if err != nil {
		return nil, err
	}
	return &client{
		client: c,
		cfg:    cfg,
	}, nil
}


type Test struct {
	DB_HOST string `etcd:"db_host;watch" required:"true"`
}


