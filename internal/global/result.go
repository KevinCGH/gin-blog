package g

import "fmt"

const (
	SUCCESS = 0
	FAIL    = 500
)

type Result struct {
	code int
	msg  string
}

func (e Result) Code() int {
	return e.code
}

func (e Result) Msg() string {
	return e.msg
}

var (
	_codes    = map[int]struct{}{}   // 注册过的错误码集合，防止重复
	_messages = make(map[int]string) // 根据错误码获取错误信息
)

func RegisterResult(code int, msg string) Result {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	if msg == "" {
		panic("错误信息不能为空")
	}

	_codes[code] = struct{}{}
	_messages[code] = msg

	return Result{
		code: code,
		msg:  msg,
	}
}

func GetMsg(code int) string {
	return _messages[code]
}

var (
	OkResult   = RegisterResult(SUCCESS, "OK")
	FailResult = RegisterResult(FAIL, "FAIL")
)

var (
	ErrRequest  = RegisterResult(9001, "请求参数格式错误")
	ErrDbOp     = RegisterResult(9002, "数据库操作异常")
	ErrUserAuth = RegisterResult(9003, "用户认证异常")

	ErrUserNotExist = RegisterResult(1000, "用户不存在")
	ErrPassword     = RegisterResult(1001, "用户认证异常")

	ErrTokenNotExist       = RegisterResult(1200, "Token 不存在，请重新登录")
	ErrTokenRuntime        = RegisterResult(1201, "Token 已过期，请重新登录")
	ErrTokenWrong          = RegisterResult(1202, "Token 不正确，请重新登录")
	ErrTokenType           = RegisterResult(1203, "Token 格式错误，请重新登录")
	ErrTokenCreate         = RegisterResult(1204, "Token 生成失败")
	ErrUserHasNoPermission = RegisterResult(1205, "你无权进行此操作")

	ErrSendEmail      = RegisterResult(8000, "发送邮件失败")
	ErrCodeNotExist   = RegisterResult(8001, "Code 不存在")
	ErrParseEmailCode = RegisterResult(8002, "解析邮件 Code 失败")
	ErrUserExist      = RegisterResult(8003, "该邮箱已经注册")
)
