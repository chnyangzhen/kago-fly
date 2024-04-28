package server

import (
	"fmt"
	"github.com/chnyangzhen/kago-fly/pkg/constant"
	"github.com/chnyangzhen/kago-fly/pkg/logger"
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
}

func Run() error {
	fmt.Printf(constant.Banner)
	server := NewServer()
	server.StartGraceful()
	server.StartGraceful()
	return nil
}

// 初始化生命周期
func initLifecycle() {
	logLifecycle := logger.NewLogLifecycle()
	RegisterDestroy(logLifecycle)

	pprofLifecycle := NewPprofLifecycle()
	RegisterPrepare(pprofLifecycle)
}
