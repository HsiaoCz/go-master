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
	// 配置日志轮转
	logFile := &lumberjack.Logger{
		Filename:   "app.log", // 日志文件路径
		MaxSize:    10,        // 每个日志文件的最大大小（MB）
		MaxBackups: 3,         // 保留的旧日志文件数量
		MaxAge:     28,        // 保留旧日志文件的最大天数
		Compress:   true,      // 是否压缩旧日志文件
	}

	// 可选：同时输出到控制台和文件
	L.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// 设置日志级别
	L.SetLevel(logrus.DebugLevel)
}
