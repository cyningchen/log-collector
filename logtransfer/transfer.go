package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"logAgent/es"
)

var (
	Global      LogConfig
	kafkaClient KafkaConsumer
)

type LogConfig struct {
	KafkaAddr string
	ESAddr    string
	LogPath   string
	LogLevel  string
	Topic     string
}

type KafkaConsumer struct {
	client sarama.Consumer
	addr   string
	topic  string
}

func InitConf(filename string) (err error) {
	conf, err := config.NewConfig("ini", filename)
	if err != nil {
		fmt.Println("new config failed, ", err)
		return
	}

	Global.KafkaAddr = conf.String("kafka::server_addr")
	Global.Topic = conf.String("kafka::topic")
	Global.ESAddr = conf.String("es::server_addr")
	Global.LogPath = conf.String("log::path")
	Global.LogLevel = conf.String("log::log_level")
	return
}

func InitLogger() {
	conf := make(map[string]interface{})
	conf["filename"] = Global.LogPath
	conf["level"] = convertLoglevel(Global.LogLevel)

	confStr, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("init logger failed, err ", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(confStr))
	logs.Info("logger init success")
}

func convertLoglevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "error":
		return logs.LevelError
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func InitKafka(addr string) (err error) {
	consumer, err := sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		fmt.Println("new consumer failed,", err)
		return
	}
	kafkaClient.client = consumer
	logs.Info("kafka consumer init success")
	return
}

func main() {
	err := InitConf("transfer.conf")
	if err != nil {
		logs.Error(err)
	}

	InitLogger()

	err = InitKafka(Global.KafkaAddr)
	if err != nil {
		logs.Error(err)
	}

	err = es.InitEs(Global.ESAddr)
	if err != nil {
		logs.Error(err)
	}

	fmt.Printf("%v", Global)

	err = Run()
	if err != nil {
		logs.Error("run error, ", err)
	}

}
