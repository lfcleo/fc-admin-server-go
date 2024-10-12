package v1

import (
	"errors"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/database"
	"fc-admin-server-go/models/request"
	"fc-admin-server-go/models/response"
	"fc-admin-server-go/pkg/config"
	redisUtil "fc-admin-server-go/pkg/redis"
	"fc-admin-server-go/pkg/upload"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// Login 密码登录，使用手机号/邮箱
func Login(c *gin.Context) {
	var aAuth request.Auth
	if err := c.ShouldBindBodyWithJSON(&aAuth); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//正则验证手机号/邮箱是否符合要求
	isMobile := util.ValidateMobileNumber(aAuth.Username)
	isEmail := util.ValidateEmail(aAuth.Username)
	if isMobile == false && isEmail == false {
		response.Json(response.ERROR, "手机号/邮箱格式错误", c)
		return
	}
	//根据手机号/邮箱查询管理员信息
	queryKey := "mobile"
	if isEmail {
		queryKey = "email"
	}
	//对传递过来的密码进行aes解密操作
	ivKey := strconv.FormatInt(aAuth.Timestamp*1000, 10)
	panPwdBytes, err := util.AesCBCPkcs7Decrypt(aAuth.Password, config.Data.Server.PasswordSign, ivKey)
	panPwd := string(panPwdBytes)
	if err != nil || panPwd == "" {
		response.Json(response.ERROR, "服务器验证密码错误", c)
		return
	}
	//数据库查询管理员信息
	administrator, err := database.FindAdministratorByKey(queryKey, aAuth.Username, "RolesData")
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "服务器开小差了", c)
		return
	}
	if administrator == nil {
		response.Json(response.ERROR, "手机号/邮箱错误", c)
		return
	}
	//判断账号密码是否正确
	if panPwd != administrator.Password {
		response.Json(response.ERROR, "密码错误", c)
		return
	}
	// 判断账户是否停用
	if administrator.Status != 1 {
		response.Json(response.AuthStatusError, "当前管理员状态不可用，请联系超级管理员！", c)
		return
	}

	GenAuthInfo(administrator, aAuth, c)
}

// VerificationCode 发送验证码
func VerificationCode(c *gin.Context) {
	var aAuth request.Auth
	if err := c.ShouldBindBodyWithJSON(&aAuth); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//验证手机号是否正确
	if isMobile := util.ValidateMobileNumber(aAuth.Username); isMobile == false {
		response.Json(response.ERROR, "手机号错误", c)
		return
	}
	vCode, err := util.GenCaptchaCode() //生成验证码（随机6位数）
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "发送短信失败1", c)
		return
	}
	//验证码保存在redis中
	if err = redisUtil.AddSmsCode(aAuth.Username, vCode); err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "发送短信失败2", c)
		return
	}
	response.Json(response.SUCCESS, vCode, c)
}

// MobileLogin 手机号登录
func MobileLogin(c *gin.Context) {
	var aAuth request.Auth
	if err := c.ShouldBindBodyWithJSON(&aAuth); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//正则验证手机号是否正确
	if isMobile := util.ValidateMobileNumber(aAuth.Username); isMobile == false {
		response.Json(response.ERROR, "手机号错误", c)
		return
	}
	vCode, err := redisUtil.GetSmsCode(aAuth.Username)
	if err != nil {
		//如果是没查询到手机号
		if err.Error() == "redigo: nil returned" {
			response.Json(response.ERROR, "请先发送验证码", c)
			return
		}
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取验证码失败", c)
		return
	}
	if vCode != aAuth.Password {
		response.Json(response.ERROR, "验证码错误", c)
		return
	}
	//验证通过，删除redis中保存的验证码
	if err = redisUtil.DelSmsCode(aAuth.Username); err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "服务器内部错误", c)
		return
	}
	//数据库查询管理员信息
	administrator, err := database.FindAdministratorByKey("mobile", aAuth.Username, "RolesData")
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "服务器开小差了", c)
		return
	}
	//如果管理员为空
	if administrator == nil {
		response.Json(response.ERROR, "该手机号还没有注册，请联系管理员注册。", c)
		return
	}
	// 判断账户是否停用
	if administrator.Status != 1 {
		response.Json(response.AuthStatusError, "当前管理员状态不可用，请联系超级管理员！", c)
		return
	}

	GenAuthInfo(administrator, aAuth, c)
}

// GenAuthInfo 生成aToken，rToken,管理员信息
func GenAuthInfo(administrator *database.Administrator, aAuth request.Auth, c *gin.Context) {
	//取出管理员角色ID数组
	var roleIDs []uint
	for _, roleData := range administrator.RolesData {
		roleIDs = append(roleIDs, roleData.ID)
	}
	// 判断用户是否有角色权限，
	if len(roleIDs) == 0 {
		response.Json(response.ERROR, "当前账号无角色权限，请联系超级管理员处理！", c)
		return
	}

	//是否免登录，免登录refreshToken过期时间长，非免登录refreshToken过期时间短。
	rTokenExpireTime := config.Data.Token.RefreshExpireTime
	if aAuth.Auto {
		rTokenExpireTime = config.Data.Token.RefreshAutoExpireTime
	}

	//取出管理员角色中最新的更新时间，加入到token中
	newTime := administrator.RolesData[0].UpdatedAt
	for _, t := range administrator.RolesData {
		if t.UpdatedAt.After(newTime) {
			newTime = t.UpdatedAt
		}
	}

	//生成token
	jwtData := util.JwtData{
		AdminID:        administrator.ID,
		AdminUpdateAt:  administrator.UpdatedAt,
		PwdUpdateAt:    administrator.UpdatedAt,
		RoleIDs:        roleIDs,
		RoleUpdateUnix: newTime.Unix(),
	}
	aToken, rToken, err := util.GenAdminARToken(jwtData, rTokenExpireTime)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "创建token失败", c)
		return
	}

	//管理员信息存储redis中
	rAI := redisUtil.RedisAdminInfo{
		AdminStatus:   administrator.Status,
		AdminUpdateAt: administrator.UpdatedAt,
		PwdUpdateAt:   administrator.UpdatedAt,
		RoleIDs:       roleIDs,
	}
	//是否做同端唯一登录
	if config.Data.Token.Unique {
		rAI.Token = aToken
	}
	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	err = redisUtil.AddAdminInfo(administrator.ID, appTypeString, rAI)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储管理员信息错误", c)
		return
	}

	administratorModel := response.Auth{
		ID:     administrator.ID,
		Mobile: administrator.Mobile,
		Email:  administrator.Email,
		Name:   administrator.Name,
		Avatar: administrator.Avatar,
		Sex:    administrator.Sex,
		Roles:  administrator.Roles,
	}

	response.Json(response.SUCCESS, gin.H{
		"userInfo":     administratorModel,
		"token":        aToken,
		"refreshToken": rToken,
	}, c)
}

// RefreshToken 刷新token
func RefreshToken(c *gin.Context) {
	var rT request.RefreshToken
	if err := c.ShouldBindBodyWithJSON(&rT); err != nil {
		response.Json(response.AuthError, "参数错误", c)
		return
	}

	//判断refreshToken是否可以使用
	rToken, err := util.VerifyAdminToken(rT.Token)
	if err != nil {
		//判断refreshToken是否过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			response.Json(response.AuthError, "token expire", c)
			return
		}
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.AuthError, "verify token error", c)
		return
	}
	//如果token没有通过校验，返回401
	if !rToken.Valid {
		response.Json(response.AuthError, "token Valid false", c)
		return
	}
	//取出aToken（aToken中存储的用户信息等是新的。包括判断等，所以要解析过期的token）
	tokenString := c.GetHeader("Authorization") //从请求的header中获取token字符串
	aToken, err := util.VerifyAdminToken(tokenString)
	//此token已经是过期了的,所以这样判断
	if !errors.Is(err, jwt.ErrTokenExpired) || aToken == nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.AuthError, "verify token error", c)
		return
	}
	//取aToken中保存的用户信息
	claims, err := util.ParseAdminToken(aToken)
	if err != nil && err.Error() != "无效的Token" {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.AuthError, "token Valid false", c)
		return
	}
	//从redis中获取存储的用户信息
	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	redisAdminInfo, err := redisUtil.GetAdminInfo(claims.AdminID, appTypeString)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			global.FC_LOGGER.Error(fmt.Sprint(err))
		}
		response.Json(response.AuthError, "get admin info error", c)
		return
	}
	//是否做同端唯一登录,比较token
	if config.Data.Token.Unique {
		if redisAdminInfo.Token != tokenString {
			response.Json(response.AuthError, "token参数错误2", c)
			c.Abort()
			return
		}
	}
	//判断管理员状态是否可用
	if redisAdminInfo.AdminStatus != 1 {
		response.Json(response.AuthError, "当前账号被冻结,联系超级管理员处理", c)
		return
	}
	//判断管理员密码的更新时间，是否需要登录
	if claims.PwdUpdateAt.Equal(redisAdminInfo.PwdUpdateAt) == false {
		response.Json(response.AuthPwdUpdate, "当前账号密码已修改，请重新登录！", c)
		return
	}
	//生成新的token
	newToken, err := util.GenAdminAToken(claims.JwtData, time.Now().Add(config.Data.Token.AccountExpireTime*time.Minute))
	if err != nil {
		response.Json(response.AuthError, "token error", c)
		return
	}
	//是否做同端唯一登录,是的话更改redis中存储的管理员token
	if config.Data.Token.Unique {
		redisAdminInfo.Token = newToken
		err = redisUtil.AddAdminInfo(claims.AdminID, appTypeString, redisAdminInfo)
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			response.Json(response.AuthError, "redis存储管理员信息错误", c)
			return
		}
	}
	response.Json(response.SUCCESS, newToken, c)
}

// MenuPermissionList 获取动态菜单和权限
func MenuPermissionList(c *gin.Context) {
	adminID := c.MustGet("AID").(uint)
	//如果是超级管理员，获取所有
	if adminID == 1 {
		list, err := database.QueryMenuList("")
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			response.Json(response.ERROR, "获取列表信息失败", c)
			return
		}
		//treeList, _ := database.MenuTree(list, 0, true)
		treeList := database.BuildMenuTree(list)
		//查询用户的权限
		response.Json(response.SUCCESS, gin.H{
			"menu":        treeList,
			"permissions": []string{"ALL"},
		}, c)
	} else {
		list, permissions, err := database.QueryAdminMenu(adminID)
		if err != nil {
			response.Json(response.ERROR, "test", c)
			return
		}

		//查询用户的权限
		response.Json(response.SUCCESS, gin.H{
			"menu":        list,
			"permissions": permissions,
		}, c)
	}

}

// Logout 退出登录
func Logout(c *gin.Context) {
	response.Json(response.SUCCESS, "success", c)
}

// UploadFile 上传文件接口（图片等）
func UploadFile(c *gin.Context) {
	adminID := c.MustGet("AID").(uint)
	adminIDStr := strconv.FormatUint(uint64(adminID), 10)
	//获取要上传的文件夹路径
	director := c.PostForm("director")
	//获取上传的图片
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "上传图片错误", c)
		return
	}
	//解析文件的类型
	fileHeaderBytes := make([]byte, 512)
	if _, err := file.Read(fileHeaderBytes); err != nil {
		response.Json(response.ERROR, "图片格式错误", c)
		return
	}

	//判断文件类型,
	if upload.CheckHttpTypeImage(http.DetectContentType(fileHeaderBytes)) {
		//判断图片格式是否正确
		if !upload.CheckImageExt(fileHeader.Filename) {
			response.Json(response.ERROR, "图片格式错误", c)
			return
		}
		//判断图片大小
		if !upload.CheckImageSizeByNum(fileHeader.Size) {
			response.Json(response.ERROR, "图片太大了，请压缩后上传", c)
			return
		}
		imageName := upload.SetImageName(adminIDStr, fileHeader.Filename) //为文件设置新名称
		savePath := upload.GetImagePath()                                 //保存的目录 document/images/
		fullPath := upload.GetImageFullPath(adminIDStr + "/" + director)  //图片在项目中的目录 runtime/document/images/用户uuid/"director上传的名称"/
		src := fullPath + imageName                                       //图片在项目中的位置

		//检查文件路径，这里面做了包括创建文件夹，检查权限等操作
		err = upload.CheckImage(fullPath)
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			response.Json(response.ERROR, "上传图片失败1", c)
			return
		}
		//使用c.SaveUploadedFile把上传的文件移动到指定位置
		if err = c.SaveUploadedFile(fileHeader, src); err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			response.Json(response.ERROR, "上传图片失败2", c)
			return
		}
		path := savePath + adminIDStr + "/" + director + "/" + imageName       //拼接完整的图片地址
		response.Json(response.SUCCESS, config.Data.Server.DomainName+path, c) //返还给前端图片路径
	}
}

// UpdateAdminInfo 更新管理员信息（管理员自行更改）
func UpdateAdminInfo(c *gin.Context) {
	var rAI response.Auth
	if err := c.ShouldBindBodyWithJSON(&rAI); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	adminID := c.MustGet("AID").(uint)
	var administrator database.Administrator
	administrator.ID = adminID
	administrator.Name = rAI.Name
	administrator.Sex = rAI.Sex
	administrator.Avatar = rAI.Avatar
	// 其余角色只可以修改姓名，性别，头像，超级管理员可以修改邮箱，姓名，性别，头像。
	if adminID == 1 {
		administrator.Mobile = rAI.Mobile
		administrator.Email = rAI.Email
	}
	// 数据库更新
	err := administrator.UpdateAdministrator()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "修改信息失败", c)
		return
	}

	//更新redis中管理员信息
	newToken, err := redisUtil.JwtSetAdminInfoUpdateAt(administrator.UpdatedAt, c)
	if err != nil {
		response.Json(response.AuthError, "更新缓存失败", c)
		return
	}

	rAI = response.Auth{
		ID:     administrator.ID,
		Mobile: administrator.Mobile,
		Email:  administrator.Email,
		Name:   administrator.Name,
		Avatar: config.Data.Server.DomainName + administrator.Avatar,
		Sex:    administrator.Sex,
		Roles:  administrator.Roles,
	}

	c.Header("Token", newToken)
	response.Json(response.SUCCESS, rAI, c)
}

// UpdateAdminPwd 更新管理员密码（管理员自行更改）
func UpdateAdminPwd(c *gin.Context) {
	var rP request.SetPassword
	if err := c.ShouldBindBodyWithJSON(&rP); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	adminID := c.MustGet("AID").(uint)
	//查询管理员信息
	administrator, err := database.FindAdministratorByKey("id", adminID)
	if err != nil && administrator == nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取管理员信息失败", c)
		return
	}

	//对传递过来的密码进行aes解密操作
	ivKey := strconv.FormatInt(rP.Timestamp*1000, 10)
	pwdBytes, err := util.AesCBCPkcs7Decrypt(rP.UsePassword, config.Data.Server.PasswordSign, ivKey)
	pwd := string(pwdBytes)
	if err != nil || pwd == "" {
		response.Json(response.ERROR, "服务器验证密码错误", c)
		return
	}
	//判断账号密码是否正确
	if pwd != administrator.Password {
		response.Json(response.ERROR, "密码错误", c)
		return
	}

	//对传递过来的新密码进行aes解密操作
	newPwdBytes, err := util.AesCBCPkcs7Decrypt(rP.NewPassword, config.Data.Server.PasswordSign, ivKey)
	newPwd := string(newPwdBytes)
	if err != nil || newPwd == "" {
		response.Json(response.ERROR, "服务器验证密码错误", c)
		return
	}

	//修改新密码
	var pwdAdmin database.Administrator
	pwdAdmin.ID = adminID
	pwdAdmin.Password = newPwd
	// 数据库更新
	err = pwdAdmin.UpdateAdministrator()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "修改信息失败", c)
		return
	}

	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	//更新redis中管理员信息
	err = redisUtil.JwtSetAdminPwdUpdateAt(adminID, appTypeString, pwdAdmin.UpdatedAt)
	if err != nil {
		response.Json(response.AuthError, "更新缓存失败", c)
		return
	}

	response.Json(response.SUCCESS, "success", c)
}

// AdminInfo 获取管理员信息
func AdminInfo(c *gin.Context) {
	adminID := c.MustGet("AID").(uint)
	//查询管理员信息
	administrator, err := database.FindAdministratorByKey("id", adminID, "RolesData")
	if err != nil && administrator == nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取管理员信息失败", c)
		return
	}
	rAI := response.Auth{
		ID:     administrator.ID,
		Mobile: administrator.Mobile,
		Email:  administrator.Email,
		Name:   administrator.Name,
		Avatar: administrator.Avatar,
		Sex:    administrator.Sex,
		Roles:  administrator.Roles,
	}
	response.Json(response.SUCCESS, rAI, c)
}

// DictInfo 获取字典信息
func DictInfo(c *gin.Context) {
	var aD request.AuthDict
	if err := c.ShouldBindBodyWithJSON(&aD); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	dT, err := database.QueryDictType(aD.Key)
	if err != nil && dT == nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取字典信息失败", c)
		return
	}
	response.Json(response.SUCCESS, dT.SysDictData, c)
}
