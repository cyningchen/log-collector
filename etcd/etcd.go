package etcd

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"net"
	"time"
	"logAgent/tailf"
	"encoding/json"
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
			if ipnet.IP.To4() != nil{
				localIp = append(localIp, ipnet.IP.String())
			}
		}
	}
	fmt.Println("localIp: ", localIp)
}

func InitEtcd(addr string, key string) (collectconf []tailf.CollectConf, err error) {
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
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey{
				err := json.Unmarshal(v.Value, &collectconf)
				if err != nil{
					logs.Error("unmarshal failed, ", err)
					continue
				}
				logs.Info("log config is %v", collectconf)
			}
		}
	}
	return
}

