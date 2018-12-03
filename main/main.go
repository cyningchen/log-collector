package main

import (
	"fmt"
	"logAgent/config"
	 log "logAgent/logs"
	"github.com/astaxie/beego/logs"
	"logAgent/tailf"
	"logAgent/kafka"
)

func main() {
	filename := "./conf/logagent.conf"
	err := config.LoadConf(filename)
	if err != nil{
		fmt.Println("load conf failed, ", err)
		return
	}

	log.InitLogger()
	logs.Info("init logger sucess")
	fmt.Println(config.Global)

	err = tailf.InitTail(config.Global.CollectConf)
	if err != nil{
		logs.Error("init tail failed, ", err)
		return
	}
	err = kafka.InitKafka(config.Global.KafkaAddr)
	if err != nil{
		logs.Error("init kafka failed, ", err)
		return
	}
	logs.Info("init all sucess")
	err = ServerRun()
	if err != nil{
		logs.Error("server run error, ", err)
	}
}
