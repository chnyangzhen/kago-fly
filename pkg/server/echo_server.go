package server

import (
	"context"
	"fmt"
	"github.com/chnyangzhen/kago-fly/pkg/config"
	"github.com/chnyangzhen/kago-fly/pkg/constant"
	"github.com/chnyangzhen/kago-fly/pkg/helper"
	"github.com/chnyangzhen/kago-fly/pkg/logger"
	"github.com/chnyangzhen/kago-fly/pkg/response"
	"github.com/chnyangzhen/kago-fly/pkg/tidctx"
	"github.com/chnyangzhen/kago-fly/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	s               *Server
	shutdownSignals = []os.Signal{os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM}
)

func init() {
	s = &Server{echo.New(), sync.Map{}, &sync.WaitGroup{}}
}

// Server 服务器实例
type Server struct {
	*echo.Echo // Web服务实例的包装，如echo框架的包装
	routes     sync.Map
	waiting    *sync.WaitGroup
}

func (s *Server) prepare() error {
	for _, prepare := range PrepareLifecycle() {
		logger.Infof("Prepare lifecycle title: %s is ready.", prepare.Title())
		if err := prepare.OnPrepare(); err != nil {
			logger.Errorf("Prepare lifecycle title: %s error with %s", prepare.Title(), err.Error())
			return err
		}
		logger.Infof("Prepare lifecycle title: %s completed.", prepare.Title())
	}
	return nil
}

func NewServer() *Server {

	webConfig := config.GetWrapper("listeners.web")
	e := s.Echo

	targetHeader := config.GetStringWithDefault("echo.request-id", constant.XRequestID)
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return helper.Uuid()
		},
		RequestIDHandler: func(ctx echo.Context, tid string) {
			tidctx.InitWebTid(ctx, tid)
		},

		TargetHeader: targetHeader,
	}))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if helper.IsNil(err) {
			return
		}
		var r *response.Result
		switch x := err.(type) {
		case *response.InnerError:
			r = response.NewInnerErrorFailedWith(x, tidctx.WebTid(c))
		case *response.ParamError:
			r = response.NewParamErrorWith(x, tidctx.WebTid(c))
		default:
			r = response.NewFailed("unknown error", tidctx.WebTid(c))
		}
		WriteJson(c, r)
	}

	e.RouteNotFound("/*", func(c echo.Context) error {
		notFound := response.NewFailed("api not found", tidctx.WebTid(c))
		return WriteJsonWithCode(c, http.StatusNotFound, notFound)
	})

	// 设置BodyLimit
	if bodyLimit := webConfig.GetString("features.body_limit"); bodyLimit != "" {
		logger.Infof("开启BodyLimit限制, body-limit: size= %s", bodyLimit)
		e.Pre(middleware.BodyLimit(bodyLimit))
	}

	// CORS（是否开启支持跨域访问特性）
	if enabled := webConfig.GetBool("features.cors_enable"); enabled {
		logger.Infof("开启跨域访问")
		e.Pre(middleware.CORS())
	}

	// CSRF（是否开启检查跨站请求伪造特性）
	if enabled := webConfig.GetBool("features.csrf_enable"); enabled {
		logger.Infof("开启CSRF")
		e.Pre(middleware.CSRF())
	}

	e.HideBanner = true

	// 捕获error，InnerError错误可以直接panic，不打印堆栈，其他错误需要打印堆栈
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll:   true,
		DisablePrintStack: true,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			_, ok := err.(*response.ParamError)
			if !ok {
				tid := tidctx.WebTid(c)
				logger.TraceId(tid).Errorw("recover error", "err", err, "stack", string(stack))
			}
			return err
		},
	}))
	e.Validator = &validator.DataValidator{}
	return s
}

func (s *Server) StartGraceful() {
	quit := make(chan os.Signal, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorw("系统运行异常,即将停止,请检查!", "error", r)
				quit <- os.Interrupt
			}
		}()
		if err := s.start(); err != nil {
			quit <- os.Interrupt
		}
	}()

	signal.Notify(quit, shutdownSignals...)
	sig := <-quit
	logger.Info("receive signal: ", sig)

	s.stop(10 * time.Second)
}

func (s *Server) start() error {
	// 生命周期准备阶段
	s.prepare()

	config := config.GetWrapper("listeners.web")

	s.waiting.Add(1)
	// 服务器启动
	// 服务器初始化
	// 自定义的前置过滤器
	go func(e *echo.Echo, routes sync.Map, waiting *sync.WaitGroup) {
		routes.Range(func(key, value interface{}) bool {
			r := value.(*apiRoute)
			e.Add(r.method, r.path, r.handler, r.middleware...)
			return true
		})
		waiting.Done()
		e.Start(config.GetString("address") + ":" + config.GetString("port"))
	}(s.Echo, s.routes, s.waiting)
	s.waiting.Wait()
	return s.StartedAfter()
}

func (s *Server) StartedAfter() error {
	for _, startedAfter := range AfterLifecycle() {
		logger.Infof("After lifecycle title: %s is ready.", startedAfter.Title())
		if err := startedAfter.OnAfter(); err != nil {
			logger.Errorf("After lifecycle title: %s error with %s", startedAfter.Title(), err.Error())
			return err
		}
		logger.Infof("After lifecycle title: %s completed.", startedAfter.Title())
	}
	return nil
}

func (s *Server) registerRouter(method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	s.Echo.Add(method, path, handler, middleware...)
}

func (s *Server) destroy(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		logger.Errorw("server shutdown error", "error", err)
	}

	for _, destroy := range DestroyLifecycle() {
		fmt.Printf("destroy title: %s is ready.\n", destroy.Title())
		if err := destroy.OnDestroy(ctx); err != nil {
			fmt.Errorf("lifecycle destroy title: %s error: %s\n", destroy.Title(), err.Error())
		} else {
			fmt.Printf("destroy title: %s completed.\n", destroy.Title())
		}
	}
	return nil
}

func (s *Server) stop(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.destroy(ctx)
}

func WriteFailed(ctx echo.Context, err error) error {
	return ctx.JSON(200, response.NewFailed(err.Error(), tidctx.WebTid(ctx)))
}

func WriteSuccess(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, response.NewSuccess(data, tidctx.WebTid(ctx)))
}

func WriteJson(ctx echo.Context, data *response.Result) error {
	data.Tid = tidctx.WebTid(ctx)
	return ctx.JSON(200, data)
}

func WriteJsonWithCode(ctx echo.Context, code int, data *response.Result) error {
	data.Tid = tidctx.WebTid(ctx)
	return ctx.JSON(code, data)
}

func RegisterRoute(method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	route := &apiRoute{
		method:     method,
		path:       path,
		handler:    handler,
		middleware: middleware,
	}
	api := method + ":" + path
	if _, ok := s.routes.Load(api); ok {
		panic(api + " already exists.")
	}
	s.routes.Store(api, route)
}

type apiRoute struct {
	method     string
	path       string
	handler    echo.HandlerFunc
	middleware []echo.MiddlewareFunc
}
