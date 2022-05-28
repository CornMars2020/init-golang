package config

import (
	"fmt"
	"init-golang/libs/logger"
	"io"
	"os"
	"path"
	"sync"
)

var (
	// 默认日志
	defaultLogger *logger.Logger

	// apiLogger API请求返回日志
	apiLogger *logger.Logger
	// priceLogger 价格日志
	priceLogger *logger.Logger
)

var (
	// 默认日志
	defaultCFLogger LoggerMap

	// spotCFLogger
	spotCFLogger LoggerMap
	// swapCFLogger
	swapCFLogger LoggerMap
	// cacheCFLogger
	cacheCFLogger LoggerMap
	// statCFLogger
	statCFLogger LoggerMap
	// APICFLogger API请求返回日志
	APICFLogger LoggerMap
	// knockCFLogger 敲单日志
	knockCFLogger LoggerMap
	// PriceCFLogger 价格日志
	priceCFLogger LoggerMap
)

// LoggerMap 日志多例
type LoggerMap struct {
	sync.Map
}

// Len 获取包含元素个数
func (mp *LoggerMap) Len() int {
	length := 0

	mp.Range(func(k interface{}, v interface{}) bool {
		length = length + 1
		return true
	})

	return length
}

// Copy 拷贝数据
func (mp *LoggerMap) Copy(src *LoggerMap) {
	mp.Range(func(k interface{}, v interface{}) bool {
		mp.Delete(k)
		return true
	})
	src.Range(func(k interface{}, v interface{}) bool {
		mp.Store(k, v)
		return true
	})
}

// Read 读取
func (mp *LoggerMap) Read(key string) (*CFLogger, bool) {
	exist, ok := mp.Load(key)
	if !ok {
		return nil, false
	}
	return exist.(*CFLogger), true
}

// Loop 循环遍历数据
func (mp *LoggerMap) Loop(f func(key string, order *CFLogger)) {
	mp.Range(func(k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			return false
		}

		order, ok := v.(*CFLogger)
		if !ok {
			return false
		}

		f(key, order)

		return true
	})
}

// ToMap 转换成普通 Map
func (mp *LoggerMap) ToMap() map[string]*CFLogger {
	ordermap := make(map[string]*CFLogger)
	mp.Range(func(k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			return false
		}

		order, ok := v.(*CFLogger)
		if !ok {
			return false
		}

		ordermap[key] = order

		return true
	})
	return ordermap
}

// CFLogger 封装支持自定义等级的日志
type CFLogger struct {
	Level  int16
	Prefix string
	Logger *logger.Logger
}

// IsTraceEnabled 是否开启 Trace 日志
func (cf *CFLogger) IsTraceEnabled() bool {
	return cf.Level >= int16(logger.TraceLevel)
}

// IsDebugEnabled 是否开启 Debug 日志
func (cf *CFLogger) IsDebugEnabled() bool {
	return cf.Level >= int16(logger.DebugLevel)
}

// IsInfoEnabled 是否开启 Info 日志
func (cf *CFLogger) IsInfoEnabled() bool {
	return cf.Level >= int16(logger.InfoLevel)
}

// IsWarnEnabled 是否开启 Warn 日志
func (cf *CFLogger) IsWarnEnabled() bool {
	return cf.Level >= int16(logger.WarnLevel)
}

// IsErrorEnabled 是否开启 Error 日志
func (cf *CFLogger) IsErrorEnabled() bool {
	return cf.Level >= int16(logger.ErrorLevel)
}

// IsFatalEnabled 是否开启 Fatal 日志
func (cf *CFLogger) IsFatalEnabled() bool {
	return cf.Level >= int16(logger.FatalLevel)
}

// IsPanicEnabled 是否开启 Panic 日志
func (cf *CFLogger) IsPanicEnabled() bool {
	return cf.Level >= int16(logger.PanicLevel)
}

// Trace Trace级别日志
func (cf *CFLogger) Trace(format string, args ...interface{}) {
	if cf.IsTraceEnabled() {
		format = fmt.Sprintf("%s %s", cf.Prefix, format)
		cf.Logger.Tracef(format, args...)
	}
}

// Debug Debug 级别日志
func (cf *CFLogger) Debug(format string, args ...interface{}) {
	if cf.IsDebugEnabled() {
		format = fmt.Sprintf("%s %s", cf.Prefix, format)
		cf.Logger.Debugf(format, args...)
	}
}

// Info Info级别日志
func (cf *CFLogger) Info(format string, args ...interface{}) {
	if cf.IsInfoEnabled() {
		format = fmt.Sprintf("%s %s", cf.Prefix, format)
		cf.Logger.Infof(format, args...)
	}
}

// Warn Warn级别日志
func (cf *CFLogger) Warn(format string, args ...interface{}) {
	if cf.IsWarnEnabled() {
		format = fmt.Sprintf("%s %s", cf.Prefix, format)
		cf.Logger.Warnf(format, args...)
	}
}

// Warning Warn级别日志
func (cf *CFLogger) Warning(format string, args ...interface{}) {
	if cf.IsWarnEnabled() {
		format = fmt.Sprintf("%s %s", cf.Prefix, format)
		cf.Logger.Warnf(format, args...)
	}
}

// Error Error级别日志
func (cf *CFLogger) Error(format string, args ...interface{}) {
	if cf.IsErrorEnabled() {
		format = fmt.Sprintf("%s %s", cf.Prefix, format)
		cf.Logger.Errorf(format, args...)
	}
}

// Fatal Fatal级别日志
func (cf *CFLogger) Fatal(format string, args ...interface{}) {
	if cf.IsFatalEnabled() {
		format = fmt.Sprintf("%s %s", cf.Prefix, format)
		cf.Logger.Fatalf(format, args...)
	}
}

// Panic Panic级别日志
func (cf *CFLogger) Panic(format string, args ...interface{}) {
	if cf.IsPanicEnabled() {
		format = fmt.Sprintf("%s %s", cf.Prefix, format)
		cf.Logger.Panicf(format, args...)
	}
}

// SetLevel 设置日志等级
func (cf *CFLogger) SetLevel(level int16) {
	cf.Level = level
}

// DefaultLogger 获取 Default 日志实例
func DefaultLogger(key string) *CFLogger {
	cfLogger, ok := defaultCFLogger.Read(key)
	if defaultLogger == nil {
		initDefaultLogger()
	}
	if cfLogger == nil || !ok {
		cfLogger = &CFLogger{
			Level:  int16(logger.TraceLevel),
			Prefix: key,
			Logger: defaultLogger,
		}
		defaultCFLogger.Store(key, cfLogger)
	}
	return cfLogger
}

// APILogger 获取 API 请求返回日志
func APILogger(key string) *CFLogger {
	cfLogger, ok := APICFLogger.Read(key)
	if apiLogger == nil {
		initAPILogger()
	}
	if cfLogger == nil || !ok {
		cfLogger = &CFLogger{
			Level:  int16(logger.TraceLevel),
			Prefix: key,
			Logger: apiLogger,
		}
		APICFLogger.Store(key, cfLogger)
	}
	return cfLogger
}

// PriceLogger 对标价格数据
func PriceLogger(key string) *CFLogger {
	cfLogger, ok := priceCFLogger.Read(key)
	if priceLogger == nil {
		initPriceLogger()
	}
	if cfLogger == nil || !ok {
		cfLogger = &CFLogger{
			Level:  int16(logger.TraceLevel),
			Prefix: key,
			Logger: priceLogger,
		}
		priceCFLogger.Store(key, cfLogger)
	}
	return cfLogger
}

var loggerMaps map[string]*logger.Logger = make(map[string]*logger.Logger)
var loggerFormatter = logger.FormatterNginx{}
var loggerFilePath = "."  // 默认当前文件夹
var loggerRotateMode = "" // hour:小时分割 day:天分割 "":不分割
var namePrefix = ""       // 日志文件名前缀

// GetMMLogger 基于交易对存储日志工厂方法
func GetMMLogger(symbol string) *logger.Logger {
	if symbol == "" {
		symbol = "default"
	}

	ins, ok := loggerMaps[symbol]
	if !ok {
		// TODO 此处加锁, 或者使用 sync map
		ins = logger.New()
		ins.SetFormatter(&loggerFormatter)
		ins.SetOutput(io.MultiWriter(ins.NewLogFile(&logger.WriterFile{Filename: loggerFilePath + "/" + symbol + ".log", RotateMode: loggerRotateMode})))
		ins.SetLevel(logger.TraceLevel)

		loggerMaps[symbol] = ins
	}
	return ins
}

// InitLog 初始化日志配置
func InitLog(prefix string, loggerCfg Logger) error {

	// logger.SetFormatter(&logger.TextFormatter{DisableTimestamp: true})
	// logger.SetFormatter(&logger.FormatterJSON{})
	loggerFormatter = logger.FormatterNginx{
		HideKeys:        loggerCfg.IsHideKey,
		NoColors:        loggerCfg.IsColor,
		TimestampFormat: loggerCfg.TimeFormat,
	}
	loggerRotateMode = loggerCfg.FileRotateMode

	runpath, _ := os.Getwd()
	if path.IsAbs(loggerCfg.Path) {
		loggerFilePath = loggerCfg.Path
		// log.Printf("abs %s, run %s", loggerFilePath, runpath)
	} else {
		loggerFilePath = path.Join(runpath, loggerCfg.Path)
		// log.Printf("relative %s, run %s", loggerFilePath, runpath)
	}

	if prefix != "" {
		namePrefix = fmt.Sprintf("%s_", prefix)
	} else {
		namePrefix = ""
	}

	return nil
}

func initDefaultLogger() {
	defaultLogger = logger.New()
	defaultLogger.SetFormatter(&loggerFormatter)
	defaultLogger.SetOutput(io.MultiWriter(defaultLogger.NewLogFile(&logger.WriterFile{Filename: loggerFilePath + "/" + namePrefix + "default.log", RotateMode: loggerRotateMode})))
	defaultLogger.SetLevel(logger.TraceLevel)
}

func initAPILogger() {
	apiLogger = logger.New()
	apiLogger.SetFormatter(&loggerFormatter)
	// apiLogger.SetFormatter(&logger.FormatterText{DisableTimestamp: true})
	apiLogger.SetOutput(io.MultiWriter(apiLogger.NewLogFile(&logger.WriterFile{Filename: loggerFilePath + "/" + namePrefix + "api.log", RotateMode: loggerRotateMode})))
	apiLogger.SetLevel(logger.TraceLevel)
}

func initPriceLogger() {
	priceLogger = logger.New()
	priceLogger.SetFormatter(&loggerFormatter)
	// spotCacheLogger.SetFormatter(&logger.FormatterText{DisableTimestamp: true})
	priceLogger.SetOutput(io.MultiWriter(priceLogger.NewLogFile(&logger.WriterFile{Filename: loggerFilePath + "/" + namePrefix + "price.log", RotateMode: loggerRotateMode})))
	priceLogger.SetLevel(logger.TraceLevel)
}
