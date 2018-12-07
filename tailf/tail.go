package tailf

import (
	"github.com/hpcloud/tail"
	"time"
	"github.com/astaxie/beego/logs"
)

var (
	tailObjMgr *TailObjMgr
)

type CollectConf struct {
	LogPath string
	Topic   string
}

type TailObj struct {
	tail *tail.Tail
	conf CollectConf
}

type TailObjMgr struct {
	tailObjs []*TailObj
	MsgChan chan *TextMsg
}

type TextMsg struct {
	Msg string
	Topic string
}

func GetOneLine() (msg *TextMsg){
	msg = <- tailObjMgr.MsgChan
	return
}

func InitTail(conf []CollectConf) (err error) {
	tailObjMgr = &TailObjMgr{
		MsgChan: make(chan *TextMsg, 100),
	}

	if len(conf) == 0 {
		logs.Warn("no any config here")
		return
	}

	for _, v := range conf{
		obj := &TailObj{
			conf:v,
		}
		tails, errTail := tail.TailFile(v.LogPath, tail.Config{
			ReOpen:    true,
			Follow:    true,
			Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
			MustExist: false,
			Poll:      true,
		})
		if errTail != nil {
			logs.Error("init tail file failed, conf: %v , err: %v, ", conf, errTail)
			continue
		}
		obj.tail = tails
		tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)
		go readFromTail(obj)
	}
	return
}

func readFromTail(tailobj *TailObj)  {
		for true {
			msg, ok := <- tailobj.tail.Lines
			if !ok {
				logs.Warning("tail file close reopen, filename: %s\n", tailobj.tail.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			textMsg := &TextMsg{
				Msg: msg.Text,
				Topic:tailobj.conf.Topic,
			}
			tailObjMgr.MsgChan <- textMsg
		}
		return
}

// 更新tail配置
func UpdateConf(confs []CollectConf) (err error)  {
	for _, oneConf := range confs{
		var isRunning = false
		for _, obj := range tailObjMgr.tailObjs{
			if oneConf == obj.conf{
				isRunning = true
				break
			}
		}
		if isRunning{
			continue
		}
		go createTailTask(oneConf)
	}
	return
}


func createTailTask(conf CollectConf){
	obj := &TailObj{
		conf:conf,
	}
	tails, errTail := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if errTail != nil {
		logs.Error("tail file err, ", errTail)
		return
	}
	obj.tail = tails
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)
	go readFromTail(obj)
}
