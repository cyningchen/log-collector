package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logAgent/config"
	"logAgent/kafka"
	log "logAgent/logs"
	"logAgent/tailf"
	"logAgent/etcd"
)

func main() {
	filename := "./conf/logagent.conf"
	err := config.LoadConf(filename)
	if err != nil {
		fmt.Println("load conf failed, ", err)
		return
	}

	log.InitLogger()
	logs.Info("init logger sucess")
	fmt.Println(config.Global)

	collectConf, err := etcd.InitEtcd(config.Global.EtcdAddr, config.Global.EtcdKey)
        config.Global.CollectConf = collectConf
	
	if err != nil {
		logs.Error("init etcd failed, ", err)
		return
	}

	err = tailf.InitTail(collectConf)
	if err != nil {
		logs.Error("init tail failed, ", err)
		return
	}
	err = kafka.InitKafka(config.Global.KafkaAddr)
	if err != nil {
		logs.Error("init kafka failed, ", err)
		return
	}
	logs.Info("init all sucess")
	err = ServerRun()
	if err != nil {
		logs.Error("server run error, ", err)
	}
}

