package tidctx

import (
	"context"
	"github.com/chnyangzhen/kago-fly/pkg/constant"
	"github.com/labstack/echo/v4"
)

// InitWebTid 从echo.Context中获取tid并添加到context中
func InitWebTid(ctx echo.Context, tid string) {
	ctx.Set(constant.Tid, tid)
}

// InitTidCtx 为context添加tid
// 适用于非web请求
func InitTidCtx(tid string) context.Context {
	return context.WithValue(context.Background(), constant.Tid, tid)
}

// WrapWebCtx 从echo.Context中获取tid并添加到context中
func WrapWebCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), constant.Tid, c.Get(constant.Tid))
}

func WrapCtx(kv map[string]string) context.Context {
	if kv == nil || len(kv) < 1 {
		return context.Background()
	}
	// 创建一个空的context
	ctx := context.Background()
	for k, v := range kv {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

// WebTid 从echo.Context中获取tid
func WebTid(ctx echo.Context) string {
	tid := ctx.Get(constant.Tid)
	if tid == nil {
		return ""
	}
	return tid.(string)
}

// Tid 从context.Context中获取tid
func Tid(ctx context.Context) string {
	tid := ctx.Value(constant.Tid)
	if tid == nil {
		return ""
	}
	return tid.(string)
}
