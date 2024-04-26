package response

import (
	"github.com/chnyangzhen/kago-fly/pkg/helper"
)

type (
	InnerError struct {
		Code      int    // business status code
		ErrorCode string // business error code
		Message   string // error message
	}

	ParamError struct {
		InnerError
	}

	Result struct {
		Result  interface{} `json:"result"`
		Msg     string      `json:"msg"`
		Success bool        `json:"success"`
		T       int64       `json:"t"`
		Tid     string      `json:"tid"`
	}
)

func NewSuccess(result interface{}, tid string) *Result {
	return &Result{
		Result:  result,
		Success: true,
		T:       helper.GetTimeMillis(),
		Tid:     tid,
	}
}

func NewFailed(msg string, tid string) *Result {
	return &Result{
		Msg:     msg,
		Success: false,
		T:       helper.GetTimeMillis(),
		Tid:     tid,
	}
}

func NewInnerErrorFailedWith(innerError *InnerError, tid string) *Result {
	return &Result{
		Msg:     innerError.Message,
		Success: false,
		T:       helper.GetTimeMillis(),
		Tid:     tid,
	}
}

func NewParamErrorWith(paramError *ParamError, tid string) *Result {
	return &Result{
		Msg:     paramError.Message,
		Success: false,
		T:       helper.GetTimeMillis(),
		Tid:     tid,
	}
}

func (e *InnerError) Error() string {
	return e.Message
}

func NewParamError(message string) *ParamError {
	return &ParamError{
		InnerError: InnerError{
			Code:      1,
			ErrorCode: "param_error",
			Message:   message,
		},
	}
}

func NewError(code int, errorCode string, message string) *InnerError {
	return &InnerError{
		Code:      code,
		ErrorCode: errorCode,
		Message:   message,
	}
}
