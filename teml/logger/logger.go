package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var L *logrus.Logger

func init() {
	L = logrus.New()
	L.SetFormatter(&logrus.JSONFormatter{})

	// 设置日志轮换
	logfile := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	// 可选：同时输出到控制台和文件
	L.SetOutput(io.MultiWriter(os.Stdout, logfile))

	// 设置日志级别
	L.SetLevel(logrus.DebugLevel)
}
