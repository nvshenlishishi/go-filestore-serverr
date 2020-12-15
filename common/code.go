package common

// 错误码
type ErrorCode int32

const (
	_                    int32 = iota + 9999 // 9999
	StatusOK                                 // 10000 正常
	StatusParamInvalid                       // 10001 请求参数无效
	StatusServerError                        // 10002 服务出错
	StatusRegisterFailed                     // 10003 注册失败
	StatusLoginFailed                        // 10004 登陆失败
	StatusTokenInvalid                       // 10005 token无效
	StatusUserNotExists                      // 10006 用户不存在
)
