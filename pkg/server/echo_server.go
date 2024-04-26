package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"kago.fly/pkg/config"
	"kago.fly/pkg/constant"
	"kago.fly/pkg/context"
	"kago.fly/pkg/helper"
	"kago.fly/pkg/logger"
	"kago.fly/pkg/response"
	"kago.fly/pkg/validator"
	"net/http"
	"sync"
	"time"
)

var (
	s *Server
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
		if err := prepare.OnPrepare(); err != nil {
			return err
		}
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

func (s *Server) Start() error {
	// 生命周期准备阶段
	s.prepare()

	config := config.GetWrapper("listeners.web")

	s.waiting.Add(1)
	// 服务器启动
	// 服务器初始化
	// 自定义的前置过滤器
	go func(e *echo.Echo, routes sync.Map, waiting *sync.WaitGroup) {
		defer waiting.Done()
		routes.Range(func(key, value interface{}) bool {
			r := value.(*apiRoute)
			e.Add(r.method, r.path, r.handler, r.middleware...)
			return true
		})
		e.Start(config.GetString("address") + ":" + config.GetString("port"))
	}(s.Echo, s.routes, s.waiting)
	s.waiting.Wait()
	return s.StartedAfter()
}

func (s *Server) StartedAfter() error {
	for _, startedAfter := range AfterLifecycle() {
		if err := startedAfter.OnAfter(); err != nil {
			return err
		}
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
	logger.Infow("server shutdown")

	for _, shutdown := range ShutdownLifecycle() {
		if err := shutdown.OnFinalize(ctx); err != nil {
			logger.Errorw("lifecycle shutdown error", "shutdown", shutdown, "error", err)
		}
	}

	for _, destroy := range DestroyLifecycle() {
		if err := destroy.OnDestroy(ctx); err != nil {
			logger.Errorw("lifecycle destroy error", "destroy", destroy, "error", err)
		}
	}
	time.Sleep(5 * time.Second)
	return nil
}

func (s *Server) Stop(timeout time.Duration) {
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
		panic(api + " already exists")
	}
	s.routes.Store(api, route)
}

type apiRoute struct {
	method     string
	path       string
	handler    echo.HandlerFunc
	middleware []echo.MiddlewareFunc
}