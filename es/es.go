package es

import (
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic"
)

var (
	esClient *elastic.Client
)

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

func InitEs(esaddr string) (err error) {
	esClient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esaddr))
	if err != nil {
		logs.Error("connect es failed, ", err)
		return
	}
	logs.Info("es init success")
	return
}

func SendToEs(topic, data string) (err error) {
	msg := &LogMessage{
		Topic:   topic,
		Message: data,
	}
	_, err = esClient.Index().
		Index(topic).
		Type(topic).
		BodyJson(msg).
		Do(context.Background())
	if err != nil {
		panic(err)
		return
	}
	return
}
