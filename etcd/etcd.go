package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"logAgent/tailf"
	"net"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
	Keys   []string
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
			if ipnet.IP.To4() != nil {
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
		etcdClient.Keys = append(etcdClient.Keys, etcdKey)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, etcdKey)
		if err != nil {
			continue
		}
		cancel()
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				err := json.Unmarshal(v.Value, &collectconf)
				if err != nil {
					logs.Error("unmarshal failed, ", err)
					continue
				}
				logs.Info("log config is %v", collectconf)
			}
		}
	}
	initEtcdWatch(addr)
	return
}

func initEtcdWatch(addr string) {
	for _, key := range etcdClient.Keys {
		go WatchKey(addr, key)
	}
}

func WatchKey(addr, key string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("new etcd failed, ", err)
	}
	for {
		rch := cli.Watch(context.Background(), key)
		var collectConf []tailf.CollectConf
		var getConfSuccess = true
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == clientv3.EventTypePut{
					if ev.Type == clientv3.EventTypePut && string(ev.Kv.Key) == key{
						err := json.Unmarshal(ev.Kv.Value, &collectConf)
						if err != nil{
							logs.Error("key[%s], unmasharl failed, [%s]", ev.Kv.Key,err)
							getConfSuccess = false
							continue
						}
					}
				}
				if ev.Type == clientv3.EventTypeDelete{
					logs.Warn("key[%s] config deleted", key)
					continue
				}
				logs.Debug("get config from etcd, %s %q: %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSuccess{
				logs.Debug("get config from etcd succ, %v", collectConf)
				tailf.UpdateConf(collectConf)
			}
		}
	}
}

