package logs

import (
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"fmt"
	"logAgent/config"
)


func convertLoglevel(level string) int {
	switch (level){
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

func InitLogger() {
	conf := make(map[string]interface{})
	conf["filename"] = config.Global.Path
	conf["level"] = convertLoglevel(config.Global.Loglevel)

	confStr, err := json.Marshal(conf)
	if err != nil{
		fmt.Println("init logger failed, err ", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(confStr))
}
