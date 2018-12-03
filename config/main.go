package config

import (
	"github.com/astaxie/beego/config"
	"fmt"
	"logAgent/tailf"
)

var (
	Global Config
)

type Config struct {
	Loglevel string
	Path string
	CollectConf []tailf.CollectConf
	KafkaAddr string
}



func LoadConf(filename string) (err error) {
	conf, err := config.NewConfig("ini", filename)
	if err != nil{
		fmt.Println("new config failed, ", err)
		return
	}

	Global.Loglevel = conf.String("log::log_level")
	Global.Path = conf.String("log::path")
	Global.KafkaAddr = conf.String("kafka::server_addr")

	err = loadCollectConf(conf)
	if err != nil{
		fmt.Println("load collectconf failed, ", err)
		return
	}
	return
}

func loadCollectConf(conf config.Configer) (err error) {
	var cc tailf.CollectConf
	cc.LogPath = conf.String("collect::log_path")
	cc.Topic = conf.String("collect::topic")
	Global.CollectConf = append(Global.CollectConf, cc)
	return
}
