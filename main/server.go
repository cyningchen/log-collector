package main

import (
	"github.com/astaxie/beego/logs"
	"logAgent/kafka"
	"logAgent/tailf"
)

func ServerRun() (err error) {
	for {
		msg := tailf.GetOneLine()
		err := sendToKafka(msg)
		if err != nil {
			logs.Error("send to kafka failed, ", err)
		}
	}
}

func sendToKafka(msg *tailf.TextMsg) (err error) {
	err = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
