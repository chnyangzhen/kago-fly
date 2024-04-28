package cmd

import (
	"github.com/chnyangzhen/kago-fly/pkg/config"
	"github.com/chnyangzhen/kago-fly/pkg/constant"
	"github.com/chnyangzhen/kago-fly/pkg/logger"
	"github.com/chnyangzhen/kago-fly/pkg/server"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"strconv"
	"time"
)

// App 应用信息
type App struct {
	Name        string
	Version     string
	Copyright   string
	Description string
	Banner      string
}

var _app = App{
	Name:        constant.AppName,
	Version:     "1.0",
	Copyright:   "(c) " + strconv.Itoa(time.Now().Year()),
	Description: "",
	Banner:      constant.Banner,
}

func DefaultGo() {
	Go(_app, nil)
}

func Go(app App, l ...server.Lifecycle) {
	newApp := NewApp(app,
		NewActions(
			InitViperComponent(),
			InitZapLoggerComponent(),
			InitLifecycle(l...),
			RunApplication(app.Banner),
		),
	)
	err := newApp.Run(os.Args)
	if err != nil {
		panic(any(err))
	}
}

func NewApp(app App, action cli.ActionFunc) *cli.App {
	inst := &cli.App{

		// 应用名、版本、版本、描述等等基本信息填充
		Name:        app.Name,
		Version:     app.Version,
		Copyright:   app.Copyright,
		Description: app.Description,

		// 不带后缀名的文件名
		// 环境变量配置时，多个配置项使用逗号分隔，如：CONFIG_NAMES=application,hystrix,go2sky,consumer
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    constant.ConfigNames,                     // 程序中使用Name来获取配置项
				Value:   cli.NewStringSlice(constant.Application), // 参数默认值
				EnvVars: []string{constant.EnvConfigNames},        // 使用EnvVars来指定接收环境变量的名称，与Name进行了映射
				Usage:   "application config names",               // 功能描述
			},

			&cli.StringFlag{
				Name:    constant.LogConfigName,                  // 参数名称
				Value:   constant.DefaultLogConfigName,           // 参数默认值
				EnvVars: []string{constant.EnvLogConfigName},     // 接收环境变量名称
				Usage:   "Load logger configuration from `FILE`", // 功能描述
			},
		},
		// 该程序执行的代码（未指定子命令时执行的操作），参考02:https://www.cnblogs.com/wangjq19920210/p/15352101.html
		Action: action,
	}

	// 排序"启动参数flag标志"、命令行列表
	sort.Sort(cli.FlagsByName(inst.Flags))
	sort.Sort(cli.CommandsByName(inst.Commands))
	return inst
}

// NewActions 用于组合多个ActionFunc
func NewActions(actions ...cli.ActionFunc) cli.ActionFunc {
	return func(context *cli.Context) error {
		for _, action := range actions {
			if err := action(context); err != nil {
				return err
			}
		}
		return nil
	}
}

// InitViperComponent 用于初始化Viper组件
func InitViperComponent() cli.ActionFunc {
	return func(context *cli.Context) error {
		if err := config.InitConfig(context.StringSlice(constant.ConfigNames)); err != nil {
			return err
		}
		return nil
	}
}

// InitZapLoggerComponent 用于初始化Zap日志组件
func InitZapLoggerComponent() cli.ActionFunc {
	return func(context *cli.Context) error {
		if err := logger.InitLogger(context.String(constant.LogConfigName)); err != nil {
			return err
		}
		return nil
	}
}

// RunApplication 用于启动应用
func RunApplication(banner string) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := server.Run(banner); err != nil {
			return err
		}
		return nil
	}
}

func InitLifecycle(l ...server.Lifecycle) cli.ActionFunc {
	return func(context *cli.Context) error {
		for _, lifecycle := range l {
			server.AddLifecycle(lifecycle)
		}
		return nil
	}
}
