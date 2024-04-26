package user

import (
	"context"
	"github.com/chnyangzhen/kago-fly/pkg/helper"
	"github.com/chnyangzhen/kago-fly/pkg/logger"
	"github.com/chnyangzhen/kago-fly/pkg/response"
	"github.com/chnyangzhen/kago-fly/pkg/server"
	"github.com/chnyangzhen/kago-fly/pkg/tidctx"
	"github.com/labstack/echo/v4"
)

func init() {
	server.RegisterRoute("GET", "/user", Query)
	server.RegisterRoute("POST", "/user", Post)
}

type User struct {
	Name string `validate:"required" json:"name" `
	Age  int    `validate:"gte=1,lte=130" json:"age"`
}

func Post(ctx echo.Context) error {
	user := new(User)
	ctx.Bind(user)
	ctx.Validate(user)

	// 获取echo中设置的tid
	webTid := tidctx.WebTid(ctx)
	logger.TraceId(webTid).Info("yyyyyyyyyyyyyyyyyy")

	// web调用下游服务转换context
	web := tidctx.WrapWebCtx(ctx)
	t(web)

	// 非web调用下游服务转换context
	msgId := helper.Uuid()
	nonWeb := tidctx.InitTidCtx(msgId)
	t(nonWeb)
	return server.WriteSuccess(ctx, user)
}

func Query(ctx echo.Context) error {
	r := ctx.Request()
	logger.Info("=====")
	//time.Sleep(5 * time.Second)
	select {
	case <-r.Context().Done():
		logger.Infow("请求已取消", "err", r.Context().Err())
		// 请求已取消
		ctx.JSON(200, "请求已取消")
	default:
		logger.Info("请求已完成", "default")
	}
	//ctx.JSON(200, "hello world")
	return response.NewError(1, "hello world", "666")
	//return nil
}

func t(ctx context.Context) {
	logger.TraceId(tidctx.Tid(ctx)).Info("tid info xxx")
	logger.Trace(ctx).Info("ctx info xxx")
}
