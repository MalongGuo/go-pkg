// Package errm 提供错误包装与位置信息记录, 便于在调用处快速定位问题。
//
// 该包通过构造实现接口 Cause 的错误值, 将调用点的文件、行号、函数名
// 与自定义消息一起写入错误文本, 从而在日志与返回值中携带足够的排查信息。
package errm

import (
	"errors"
	"fmt"
	"runtime"
)

type Cause interface {
	Cause() error
	error
}

// message 是 Cause 的具体实现, 保存底层错误与格式化后的错误消息。
type message struct {
	cause error
	msg   string
}

// 确保 message 实现 Cause 接口
var _ Cause = (*message)(nil)

// ---- struct 可以调用的方法 ----

// Cause 返回底层错误。
func (m message) Cause() error { return m.cause }

// Error 返回带位置信息的错误消息。
func (m message) Error() string { return m.msg }

// newWithMessage 基于底层错误与消息构造包含位置信息的错误。
//
// 参数:
//   - cause: 底层错误, 可为 nil, 会被位置信息和 msg 包装
//   - msg: 附加的错误描述, 会与位置信息一并写入
//
// 返回: 满足 Cause 接口的错误值, 其 Error 文本包含文件、行号与函数名。
func newWithMessage(cause error, msg string) Cause {
	skip := 1 + 1
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	errMsg := fmt.Sprintf("%s:%d\n%s:\n%v=>%s", file, line, funcName, cause, msg)
	cause = errors.New(errMsg)

	return message{cause: cause, msg: errMsg}
}

// New 基于给定的底层错误与消息创建一个带位置信息的错误。
//
// 参数:
//   - cause: 底层错误, 可为 nil
//   - msg: 附加的错误描述
//
// 返回: 满足 Cause 接口的错误值。
func New(cause error, msg string) Cause {
	return newWithMessage(cause, msg)
}

// WithMessagef 使用格式化字符串与可变参数构造带位置信息的错误。
//
// 参数:
//   - cause: 底层错误, 可为 nil
//   - format: 格式化模板, 语义同 fmt.Printf
//   - a: 模板参数
//
// 返回: 满足 Cause 接口的错误值。
func WithMessagef(cause error, format string, a ...any) Cause {
	return newWithMessage(cause, fmt.Sprintf(format, a...))
}

// WithMessage 基于给定的消息构造带位置信息的错误。
//
// 参数:
//   - cause: 底层错误, 可为 nil
//   - msg: 附加的错误描述
//
// 返回: 满足 Cause 接口的错误值。
func WithMessage(cause error, msg string) Cause {
	return newWithMessage(cause, msg)
}
