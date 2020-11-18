// @Author : DAIPENGYUAN
// @File : wrap
// @Time : 2020/11/6 12:46 
// @Description : 

package errno

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type Ewrap struct {
	Message    string    // 保存自定义的错误信息
	StatusCode int       // 错误状态码
	rawErr     error     // 保存原始错误信息
	stackPC    []uintptr // 保存函数调用栈指针
}

func (e *Ewrap) Error() string {
	return e.Message
}

// RawErr the origin err
func (e Ewrap) RawErr() error {
	return e.rawErr
}

// CallStack get function call stack
func (e Ewrap) CallStack() string {
	frames := runtime.CallersFrames(e.stackPC)
	var (
		f      runtime.Frame
		more   bool
		result string
		index  int
	)
	for {
		f, more = frames.Next()
		if index = strings.Index(f.File, "src"); index != -1 {
			// trim GOPATH or GOROOT prifix
			f.File = string(f.File[index+4:])
		}
		result = fmt.Sprintf("%s%s\n\t%s:%d\n", result, f.Function, f.File, f.Line)
		if !more {
			break
		}
	}
	return result
}

func wrapErr(err error, code int, fmtAndArgs ...interface{}) *Ewrap {
	msg := fmtErrMsg(fmtAndArgs...)
	if err == nil {
		err = errors.New(msg)
	}
	if e, ok := err.(*Ewrap); ok {
		if msg != "" {
			e.Message = msg
		}
		if code != 0 {
			e.StatusCode = code
		}
		return e
	}

	pcs := make([]uintptr, 32)
	// skip the first 3 invocations
	count := runtime.Callers(3, pcs)
	e := &Ewrap{
		StatusCode: code,
		Message:    msg,
		rawErr:     err,
		stackPC:    pcs[:count],
	}
	if e.Message == "" {
		e.Message = err.Error()
	}
	return e
}

// fmtErrMsg used to format error message
func fmtErrMsg(msgs ...interface{}) string {
	if len(msgs) > 1 {
		return fmt.Sprintf(msgs[0].(string), msgs[1:]...)
	}
	if len(msgs) == 1 {
		if v, ok := msgs[0].(string); ok {
			return v
		}
		if v, ok := msgs[0].(error); ok {
			return v.Error()
		}
	}
	return ""
}

// WrapErr equal to InternalErr(err)
// notice: be careful, the returned value is *MErr, not error
func WrapErr(err error, fmtAndArgs ...interface{}) *Ewrap {
	return wrapErr(err, http.StatusInternalServerError, fmtAndArgs...)
}

// WrapErrWithCode if code is not 0, update StatusCode to code,
// if fmtAndArgs is not nil, update the Message according to fmtAndArgs
// notice: be careful, the returned value is *MErr, not error
func WrapErrWithCode(err error, code int, fmtAndArgs ...interface{}) *Ewrap {
	return wrapErr(err, code, fmtAndArgs...)
}

// NotFoundErr use http.StatusNotFound as StatusCode to express not found err
// if fmtAndArgs is not nil, update the Message according to fmtAndArgs
func NotFoundErr(err error, fmtAndArgs ...interface{}) error {
	return wrapErr(err, http.StatusNotFound, fmtAndArgs...)
}

