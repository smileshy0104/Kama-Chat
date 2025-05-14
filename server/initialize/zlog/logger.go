package zlog

import (
	"Kama-Chat/global"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"runtime"
)

// logger 用于全局日志记录
var logger *zap.Logger

// logPath 存储日志文件的路径
var logPath string

// init 自动调用，用于初始化日志系统
func InitLogger() {
	// 创建一个生产环境的编码配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置日志记录中时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 日志encoder还是JSONEncoder，把日志行格式化成JSON格式
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	// 获取日志文件路径
	logPath = global.CONFIG.LogConfig.LogPath
	// 打开日志文件，如果不存在则创建
	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 644)
	// 创建一个同步写入日志文件的写同步器
	fileWriteSyncer := zapcore.AddSync(file)
	// 创建一个Tee，它可以把日志消息分发到多个核心
	core := zapcore.NewTee(
		// 创建一个核心，将DEBUG及以上级别的日志写入到标准输出
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		// 创建一个核心，将DEBUG及以上级别的日志写入到文件
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	// 创建一个Logger
	logger = zap.New(core)
}

// getFileLogWriter 返回一个写入日志文件的WriteSyncer
func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
	// 创建一个Lumberjack logger，用于滚动日志文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // 单个文件最大100M
		MaxBackups: 60,  // 多于60个日志文件后，清理较旧的日志
		MaxAge:     1,   // 一天一切割
		Compress:   false,
	}

	// 返回一个同步写入Lumberjack logger的写同步器
	return zapcore.AddSync(lumberJackLogger)
}

// getCallerInfoForLog 获得调用方的日志信息，包括函数名，文件名，行号
func getCallerInfoForLog() (callerFields []zap.Field) {
	// 回溯两层，拿到写日志的调用方的函数信息
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	// 获取函数名
	funcName := runtime.FuncForPC(pc).Name()
	// 只保留函数名，去掉包路径
	funcName = path.Base(funcName)

	// 创建包含函数名、文件名和行号的字段
	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}

// Info 记录一个INFO级别的日志消息
func Info(message string, fields ...zap.Field) {
	// 获取调用方的日志信息
	callerFields := getCallerInfoForLog()
	// 将调用方的信息添加到日志字段中
	fields = append(fields, callerFields...)
	// 记录日志消息
	logger.Info(message, fields...)
}

// Warn 记录一个WARN级别的日志消息
func Warn(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Warn(message, fields...)
}

// Error 记录一个ERROR级别的日志消息
func Error(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Error(message, fields...)
}

// Fatal 记录一个FATAL级别的日志消息，然后调用os.Exit(1)
func Fatal(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Fatal(message, fields...)
}

// Debug 记录一个DEBUG级别的日志消息
func Debug(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Debug(message, fields...)
}
