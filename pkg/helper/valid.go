package helper

import (
	"fmt"
	"reflect"
)

func IsTruePanic(shouldTrue bool, msg string) {
	if shouldTrue {
		panic(fmt.Errorf("%s", msg))
	}
}

func IsFalsePanic(shouldFalse bool, msg string) {
	if !shouldFalse {
		panic(fmt.Errorf("%s", msg))
	}
}

// IsEmptyPanic 如果输入值为空字符串，触发Panic断言，抛出错误消息。
func IsEmptyPanic(str string, msg string) {
	if str == "" {
		panic(fmt.Errorf("%s", msg))
	}
}

// IsNotEmpty 判断输入字符串是否非空
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// IsEmpty 判断输入字符串是否为空
func IsEmpty(str string) bool {
	if str == "" {
		return true
	}
	return false
}

// IsNil 判断输入值是否为Nil值，只针对引用类型判断有效，任何数值类型、结构体非指针类型等均为非Nil值。
func IsNil(i interface{}) bool {
	if nil == i {
		return true
	}
	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Interface, reflect.Slice,
		reflect.Ptr, reflect.UnsafePointer:
		return value.IsNil()
	}
	return false
}

// NotNil 判断输入值是否为非Nil值（包括：nil、类型非Nil但值为Nil），用于检查类型值是否为Nil。
// 只针对引用类型判断有效，任何数值类型、结构体非指针类型等均为非Nil值。
func NotNil(v interface{}) bool {
	return !IsNil(v)
}

// IsNotNil 对输入值断言，期望为非Nil值；断言成功时返回原值。当值为Nil时，触发Panic断言，抛出错误消息。
func IsNotNil(v interface{}, msg string) interface{} {
	if IsNil(v) {
		panic(fmt.Errorf("%s", msg))
	}
	return v
}
