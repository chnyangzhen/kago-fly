package server

import (
	"fmt"
	"kago.fly/pkg/constant"
	"kago.fly/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	initLifecycle()
}

func AddLifecycle(l Lifecycle) {
	if v, ok := l.(Preparer); ok {
		RegisterPrepare(v)
	}

	if v, ok := l.(After); ok {
		RegisterStartedAfter(v)
	}

	if v, ok := l.(Destroyer); ok {
		RegisterDestroy(v)
	}

	if v, ok := l.(Finalizer); ok {
		RegisterShutdown(v)
	}
}

func Run() error {
	fmt.Printf(constant.Banner)

	server := NewServer()

	quit := make(chan os.Signal, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorw("系统运行异常,即将停止,请检查!", "error", r)
				quit <- os.Interrupt
			}
		}()
		if err := server.Start(); err != nil {
			quit <- os.Interrupt
		}
	}()

	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM)
	sig := <-quit
	logger.Info("receive signal: ", sig)

	server.Stop(10 * time.Second)

	return nil
}

// 初始化生命周期
func initLifecycle() {
	logLifecycle := logger.NewLogLifecycle()
	RegisterDestroy(logLifecycle)

	pprofLifecycle := NewPprofLifecycle()
	RegisterPrepare(pprofLifecycle)
}
