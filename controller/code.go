package controller

type ResCode int64

const (
	CodeSucess ResCode = 1000 + iota //是 Go 常量块里 "从 1000 开始自动递增" 的写法，生成的是 1000、1001、1002
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
)

var codeMsgMap = map[ResCode]string{
	CodeSucess:          "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeNeedLogin:       "需要登录",
	CodeInvalidToken:    "无效的token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
