package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Log *zap.SugaredLogger

func init() {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100, // megabytes
		MaxBackups: 1,
		MaxAge:     7,     //days
		Compress:   false, // disabled by default
	}
	defer lumberJackLogger.Close()

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeCaller = zapcore.ShortCallerEncoder // 显示完整文件路径
	config.EncodeTime = zapcore.ISO8601TimeEncoder   // 设置时间格式
	fileEncoder := zapcore.NewConsoleEncoder(config)

	coreConsole := zapcore.NewCore(
		fileEncoder,             //编码设置
		zapcore.Lock(os.Stdout), //输出到控制台
		zap.InfoLevel,           //日志等级
	)

	coreFile := zapcore.NewCore(
		fileEncoder,                       //编码设置
		zapcore.AddSync(lumberJackLogger), //输出到文件
		zap.InfoLevel,                     //日志等级
	)

	core := zapcore.NewTee(
		coreConsole,
		coreFile,
	)

	_log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Log = _log.Sugar()

}
