package logger

import (
	"context"
	"github.com/chnyangzhen/kago-fly/pkg/constant"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// TraceLogger 带有TraceId的Logger
type TraceLogger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Panic(args ...interface{})
	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
	Panicf(fmt string, args ...interface{})
	Infow(msg string, keyAndValues ...interface{})
	Warnw(msg string, keyAndValues ...interface{})
	Errorw(msg string, keyAndValues ...interface{})
	Debugw(msg string, keyAndValues ...interface{})
	Panicw(msg string, keyAndValues ...interface{})
}

// TraceId 带有TraceId的Logger
func TraceId(tid string) TraceLogger {
	return baseLogger.Sugar().With(constant.Tid, tid)
}

// Trace 从Context中获取TraceId，如果不存在，则返回原始Logger
func Trace(ctx context.Context) TraceLogger {
	if ctx == nil {
		return baseLogger.Sugar()
	}
	traceId := ctx.Value(constant.Tid)
	if traceId != nil {
		return baseLogger.Sugar().With(zap.String(constant.Tid, cast.ToString(traceId)))
	}
	return baseLogger.Sugar()
}
