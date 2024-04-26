package logger

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"kago.fly/pkg/constant"
	"os"
)

type LogLifecycle int

func (l *LogLifecycle) OnDestroy(ctx context.Context) error {
	baseLogger.Sync()
	baseLogger.Sugar().Sync()
	errorFileWriter.Close()
	warnFileWriter.Close()
	infoFileWriter.Close()
	return nil
}

func (l *LogLifecycle) Title() string {
	return "logger"
}

func NewLogLifecycle() *LogLifecycle {
	return new(LogLifecycle)
}

var (
	baseLogger                                      *zap.Logger
	errorFileWriter, warnFileWriter, infoFileWriter *lumberjack.Logger
	consoleWriter                                   = zapcore.Lock(os.Stdout)
)

type ConfigWrapper struct {
	Default zap.Config        `json:"default" yaml:"default"`
	Rolling RollingFileConfig `json:"rolling" yaml:"rolling"`
}

type RollingFileConfig struct {
	LogFilePath   string `json:"logFilePath" yaml:"logFilePath"`     // 日志路径
	ErrorFileName string `json:"errorFileName" yaml:"errorFileName"` // 默认名称：error.log
	WarnFileName  string `json:"warnFileName" yaml:"warnFileName"`   // 默认名称：warn.log
	InfoFileName  string `json:"infoFileName" yaml:"infoFileName"`   // 默认名称：info.log
	MaxSize       int    `json:"maxSize" yaml:"maxSize"`             // 一个文件多少Ｍ（大于该数字开始切分文件）
	MaxBackups    int    `json:"maxBackups" yaml:"maxBackups"`       // MaxBackups是要保留的最大旧日志文件数
	MaxAge        int    `json:"maxAge" yaml:"maxAge"`               // MaxAge是根据日期保留旧日志文件的最大天数
	Compress      bool   `json:"compress" yaml:"compress"`           // 是否压缩
}

func InitLogger(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read config(file:%s): %w", filename, err)
	}

	config := ConfigWrapper{
		Default: zap.Config{},
		Rolling: RollingFileConfig{},
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return buildLogger(config)
}

func buildLogger(wrapper ConfigWrapper) error {

	config := wrapper.Default
	rollingConfig := wrapper.Rolling

	logEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)

	infoFileWriter = initLumberjackLogger(rollingConfig.InfoFileName, rollingConfig)
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel && level-zapcore.InfoLevel-config.Level.Level() > -1
	})

	warnFileWriter = initLumberjackLogger(rollingConfig.WarnFileName, rollingConfig)
	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.WarnLevel && zapcore.WarnLevel-config.Level.Level() > -1
	})

	errorFileWriter = initLumberjackLogger(rollingConfig.ErrorFileName, rollingConfig)
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level > zapcore.WarnLevel && zapcore.WarnLevel-config.Level.Level() > -1
	})

	consoleLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level-config.Level.Level() > -1
	})

	zapCores := []zapcore.Core{
		zapcore.NewCore(logEncoder, zapcore.AddSync(infoFileWriter), infoLevel),
		zapcore.NewCore(logEncoder, zapcore.AddSync(warnFileWriter), warnLevel),
		zapcore.NewCore(logEncoder, zapcore.AddSync(errorFileWriter), errorLevel),
		zapcore.NewCore(logEncoder, consoleWriter, consoleLevel),
	}

	l, err := config.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(zapCores...)
	}))

	if err != nil {
		return err
	}

	baseLogger = l

	return nil
}

func initLumberjackLogger(filename string, fileConfig RollingFileConfig) *lumberjack.Logger {
	// 创建info级别的lumberjack logger实例
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileConfig.LogFilePath + constant.FilepathSeparator + filename,
		MaxSize:    fileConfig.MaxSize,
		MaxBackups: fileConfig.MaxBackups,
		MaxAge:     fileConfig.MaxAge,
		Compress:   fileConfig.Compress,
	}
	return lumberjackLogger
}

func GetLogger() *zap.SugaredLogger {
	return baseLogger.Sugar()
}

func NewWith(args ...interface{}) *zap.SugaredLogger {
	return baseLogger.Sugar().With(args...)
}

func Info(args ...interface{}) {
	baseLogger.Sugar().Info(args...)
}

func Warn(args ...interface{}) {
	baseLogger.Sugar().Warn(args...)
}

func Error(args ...interface{}) {
	baseLogger.Sugar().Error(args...)
}

func Debug(args ...interface{}) {
	baseLogger.Sugar().Debug(args...)
}

func Panic(args ...interface{}) {
	baseLogger.Sugar().Panic(args...)
}

func Infof(fmt string, args ...interface{}) {
	baseLogger.Sugar().Infof(fmt, args...)
}

func Warnf(fmt string, args ...interface{}) {
	baseLogger.Sugar().Warnf(fmt, args...)
}

func Errorf(fmt string, args ...interface{}) {
	baseLogger.Sugar().Errorf(fmt, args...)
}

func Debugf(fmt string, args ...interface{}) {
	baseLogger.Sugar().Debugf(fmt, args...)
}

func Panicf(fmt string, args ...interface{}) {
	baseLogger.Sugar().Panicf(fmt, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	baseLogger.Sugar().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	baseLogger.Sugar().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	baseLogger.Sugar().Errorw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	baseLogger.Sugar().Debugw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	baseLogger.Sugar().Panicw(msg, keysAndValues...)
}
