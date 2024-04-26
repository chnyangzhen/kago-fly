package helper

import (
	"time"
)

func GetTimeMillis() int64 {
	return time.Now().UnixMilli()
}
