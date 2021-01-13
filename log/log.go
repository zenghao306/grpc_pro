package log

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strings"
)

var Logger *log

type log struct {
	logPath   string
	logLevel  int
	beeLogger *logs.BeeLogger
}

func NewLogger(logPath1, logLevel1, adapter string) {
	fmt.Println("日志参数:",logPath1, logLevel1, adapter)
	Logger = &log{
		logPath:   logPath1,
		beeLogger: logs.NewLogger(1000),
	}

	//输出文件名和行号
	Logger.beeLogger.EnableFuncCallDepth(true)
	Logger.beeLogger.SetLogFuncCallDepth(3)

	//日志级别
	switch logLevel1 {
	case "debug":
		Logger.logLevel = logs.LevelDebug
	case "info":
		Logger.logLevel = logs.LevelInfo
	case "error":
		Logger.logLevel = logs.LevelError
	default:
		Logger.logLevel = logs.LevelDebug
	}
	if adapter == "" {
		adapter = logs.AdapterFile
	}
	Logger.beeLogger.SetLogger(adapter, fmt.Sprintf(`{"filename":"%s","level":%d,"maxlines":0,"maxsize":0,"daily":true,"maxdays":60}`, Logger.logPath, Logger.logLevel))
}

func (l *log) Debug(v ...interface{}) {
	l.beeLogger.Debug(l.generateFmtStr(len(v)), v...)
}

func (l *log) Info(v ...interface{}) {
	l.beeLogger.Info(l.generateFmtStr(len(v)), v...)
}

func (l *log) Error(v ...interface{}) {
	l.beeLogger.Error(l.generateFmtStr(len(v)), v...)
}

func (l *log) generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

func Debug(v ...interface{}) {
	Logger.beeLogger.Debug(Logger.generateFmtStr(len(v)), v...)
}

func Info(v ...interface{}) {
	Logger.beeLogger.Info(Logger.generateFmtStr(len(v)), v...)
}

func Error(v ...interface{}) {
	Logger.beeLogger.Error(Logger.generateFmtStr(len(v)), v...)
}
