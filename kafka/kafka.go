package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	client sarama.SyncProducer
)

func InitKafka(addr string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	client, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logs.Debug("new kafka client error, ", err)
		return
	}

	logs.Debug("init kafka sucess")
	return
}

func SendToKafka(data, topic string) (err error) {
	msg := &sarama.ProducerMessage{
		Topic:topic,
		Value:sarama.StringEncoder(data),
	}
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send mesage error, ", err)
		return
	}
	logs.Debug("pid: %v, offset %v, topic: %s\n", pid, offset, topic)
	return
}

