package core

import (
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
)

var (
	env     = beego.AppConfig.String("env")
	Channel = beego.AppConfig.String("server_name")
	Topic   = "fx.ball.race"
	workDir = "work"
)

func IsProd() bool {
	return env == "" || env == "prod"
}

func IsDev() bool {
	return env == "dev"
}

func InitLog() {

	if IsProd() {
		beego.RunMode = "prod"
	} else {
		beego.RunMode = env
	}

	if !IsDev() {
		workdir_ := beego.AppConfig.String("work_dir")
		if workdir_ != "" {
			workDir, _ = filepath.Abs(workdir_)
		} else {
			workDir, _ := os.Getwd()
			workDir, _ = filepath.Abs(workDir)
			workDir = filepath.Join(workDir, "work")
		}
		logpath := filepath.Join(workDir, "logs", "app.log")
		logpath, _ = filepath.Abs(logpath)

		logpath = strings.Replace(logpath, "\\", "/", -1)
		beego.BeeLogger.SetLogger("file", `{"filename":"`+logpath+`","level":2}`)
		beego.BeeLogger.SetLogger("smtp", `{"username":"jxufeliujj@gmail.com","password":"liujj520*/*","host":"smtp.gmail.com:587","sendTos":["jxufeliujj@gmail.com"],"level":4}`)
		beego.BeeLogger.DelLogger("console")
		log.SetOutput(&systemLog{})
	}
}

type systemLog struct{}

func (*systemLog) Write(p []byte) (n int, err error) {
	beego.Info(string(p))
	return n, nil
}
func Writelog(loglevel int, calltype, url, userId, instructions, info string) {
	body := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s]", calltype, url, userId, instructions, info)
	switch loglevel {
	case LevelTrace:
		beego.Trace(body)
	case LevelDebug:
		beego.Debug(body)
	case LevelInfo:
		beego.Info(body)
	case LevelWarn:
		beego.Warn(body)
	case LevelError:
		beego.Error(body)
	case LevelCritical:
		beego.Critical(body)
	}
}
