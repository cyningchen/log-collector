package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"logAgent/es"
	"sync"
)

var (
	wg sync.WaitGroup
)

func Run() (err error) {
	partitionList, err := kafkaClient.client.Partitions(Global.Topic)
	for partition := range partitionList {
		pc, err := kafkaClient.client.ConsumePartition(Global.Topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to consume partiton: %d, err: %v\n", partition, err)
		}
		defer pc.AsyncClose()
		go func(pc sarama.PartitionConsumer) {
			wg.Add(1)
			for msg := range pc.Messages() {
				logs.Debug("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				err = es.SendToEs(Global.Topic, string(msg.Value))
				if err != nil {
					logs.Error("send msg to es failed, err: %v", err)
				}
			}
			wg.Done()
		}(pc)
	}
	wg.Wait()
	return
}
