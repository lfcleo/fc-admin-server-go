package util

import (
	"errors"
	"fc-admin-server-go/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtData struct {
	AdminID        uint      `json:"adminID"`        //管理员ID
	AdminUpdateAt  time.Time `json:"adminUpdateAt"`  //管理员信息更新时间，是否需要重新获取管理员信息
	PwdUpdateAt    time.Time `json:"pwdUpdateAt"`    //管理员密码更新时间，是否需重重新登录
	RoleIDs        []uint    `json:"roleIDs"`        //管理员角色ID数组
	RoleUpdateUnix int64     `json:"roleUpdateUnix"` //管理员角色的最新更新时间（取角色列表中最新的,毫秒级时间戳）
}

type JwtClaims struct {
	JwtData
	jwt.RegisteredClaims
}

// GenAdminARToken 颁发token，refresh_token
func GenAdminARToken(jwtData JwtData, rTokenExpireTime time.Duration) (aToken, rToken string, err error) {
	//设置aToken
	rc := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Data.Token.AccountExpireTime * time.Minute)), //定义过期时间
		IssuedAt:  jwt.NewNumericDate(time.Now()),                                                        // 签发时间
	}
	aJC := JwtClaims{}
	aJC.JwtData = jwtData
	aJC.RegisteredClaims = rc
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, aJC).SignedString([]byte(config.Data.Token.Secret))

	//设置rToken
	rt := rc
	rt.ExpiresAt = jwt.NewNumericDate(time.Now().Add(rTokenExpireTime * time.Hour)) //定义过期时间
	rJC := JwtClaims{}
	rJC.JwtData = jwtData
	rJC.RegisteredClaims = rt
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rJC).SignedString([]byte(config.Data.Token.Secret))
	return
}

// GenAdminAToken 颁发token,设置的是过期时间。
func GenAdminAToken(jwtData JwtData, expireTime time.Time) (aToken string, err error) {
	//设置aToken
	rc := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireTime), //定义过期时间
		IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
	}
	aJC := JwtClaims{}
	aJC.JwtData = jwtData
	aJC.RegisteredClaims = rc
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, aJC).SignedString([]byte(config.Data.Token.Secret))
	return
}

// VerifyAdminToken 验证管理员token
func VerifyAdminToken(tokenString string) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Data.Token.Secret), nil
	})
	return
}

// ParseAdminToken 解析Token管理员信息
func ParseAdminToken(token *jwt.Token) (*JwtClaims, error) {
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return claims, errors.New("无效的Token")
	}
}
