package server

import (
	"fmt"
	"kago.fly/pkg/config"
	"net/http"
	_ "net/http/pprof"
)

type PprofLifecycle int

func NewPprofLifecycle() *PprofLifecycle {
	return new(PprofLifecycle)
}

func (p *PprofLifecycle) OnPrepare() error {
	conf := config.GetWrapper("listeners.pprof")
	enable := conf.GetBool("enable")
	if !enable {
		return nil
	}
	port := conf.GetString("port")
	if port == "" {
		port = "6885"
	}
	addr := "localhost:" + port

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("pprof listen on %s\n", addr)
	}()
	return nil
}

func (p *PprofLifecycle) Title() string {
	return "pprof"
}
