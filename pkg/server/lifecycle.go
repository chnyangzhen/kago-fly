package server

import (
	"context"
	"sort"
)

var (
	prepares = make([]Preparer, 0, 8)
	destroys = make([]Destroyer, 0, 8)
	afters   = make([]After, 0, 8)
)

type (
	Lifecycle interface {
		Title() string
	}
	// Preparer 应用启动时准备数据，如Apollo只需要初始化一次，无论应用启动了多少个web server，只需要初始化一个Apollo
	Preparer interface {
		OnPrepare() error
		Title() string
	}

	After interface {
		OnAfter() error
		Title() string
	}

	// Destroyer 应用销毁前处理函数，如Apollo关闭连接、Logger关闭文件等等
	Destroyer interface {
		OnDestroy(ctx context.Context) error
		Title() string
	}
)

// RegisterPrepare 注册Prepare
func RegisterPrepare(prepare Preparer) {
	prepares = append(prepares, prepare)
}

// RegisterDestroy 注册Destroy
func RegisterDestroy(destroy Destroyer) {
	destroys = append(destroys, destroy)
}

func RegisterStartedAfter(after After) {
	afters = append(afters, after)
}

// AfterLifecycle 返回After列表的副本
func AfterLifecycle() []After {
	dst := make([]After, len(afters))
	copy(dst, afters)
	return dst
}

// PrepareLifecycle 返回Prepare列表的副本
func PrepareLifecycle() []Preparer {
	dst := make([]Preparer, len(prepares))
	copy(dst, prepares)
	return dst
}

// DestroyLifecycle 返回Destroy列表的副本
func DestroyLifecycle() []Destroyer {
	dst := make([]Destroyer, len(destroys))
	copy(dst, destroys)
	sort.Slice(dst, func(i, j int) bool {
		return i > j
	})
	return dst
}
