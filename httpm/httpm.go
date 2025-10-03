package httpm

type ResultOk struct {
	Code int32  `json:"code" default:"0"`
	Msg  string `json:"msg" default:"success"`
	Data any    `json:"data,omitempty"`
}

type ResultErrI interface {
	GetCode() int32
	GetMsg() string
	error
}

type ResultErr struct {
	Code int32  `json:"code" default:"1"`
	Msg  string `json:"msg" default:"error"`
}

func (r ResultErr) GetCode() int32 {
	return r.Code
}

func (r ResultErr) GetMsg() string {
	return r.Msg
}

func (r ResultErr) Error() string {
	return r.Msg
}

var _ ResultErrI = (*ResultErr)(nil)

func NewResultOk(data any, msg ...string) *ResultOk {
	msgStr := "success"
	if len(msg) > 0 {
		msgStr = msg[0]
	}
	return &ResultOk{
		Code: 0,
		Msg:  msgStr,
		Data: data,
	}
}

func NewResultErr(msg string, code ...int32) *ResultErr {
	codeInt := int32(1)
	if len(code) > 0 {
		codeInt = code[0]
	}
	return &ResultErr{
		Code: codeInt,
		Msg:  msg,
	}
}
