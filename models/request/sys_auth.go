package request

// Auth 请求数据结构体
type Auth struct {
	Base
	Username string `json:"username" binding:"required"` //手机号/邮箱
	Password string `json:"password"`                    //密码/验证码
	Auto     bool   `json:"auto"`                        //是否免登录
}

// RefreshToken 请求数据结构体
type RefreshToken struct {
	Token string `json:"token" binding:"required"` //refresh token
}

// SetPassword 设置密码
type SetPassword struct {
	Base
	UsePassword string `json:"usePassword"` //旧密码
	NewPassword string `json:"newPassword"` //新密码
}

// AuthDict 获取字典详情
type AuthDict struct {
	Key string `json:"key"` //字典key值
}
