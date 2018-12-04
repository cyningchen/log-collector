package etcd

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"net"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
}

var (
	etcdClient *EtcdClient
	localIp    []string
)

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logs.Error("get localIp failed, ", err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			localIp = append(localIp, ipnet.IP.String())
		}
	}
	fmt.Println("localIp: ", localIp)
}

func InitEtcd(addr string, key string) (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("new etcd failed, ", err)
	}
	etcdClient = &EtcdClient{
		client: cli,
	}
	for _, ip := range localIp {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, etcdKey)
		if err != nil {
			continue
		}
		cancel()
		for k, v := range resp.Kvs {
			fmt.Println(k, v)
		}
	}
	return err
}
