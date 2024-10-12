package v1

import (
	"errors"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/database"
	"fc-admin-server-go/models/request"
	"fc-admin-server-go/models/response"
	"fc-admin-server-go/pkg/config"
	redisutil "fc-admin-server-go/pkg/redis"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// AdminList 管理员列表
func AdminList(c *gin.Context) {
	var rA request.Administrator
	if err := c.ShouldBindBodyWithJSON(&rA); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	queryStr := rA.BuildQueryConditions()
	list, total, err := database.QueryAdministratorList(queryStr, rA.Page, rA.PageSize, "RolesData")
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取列表信息失败", c)
		return
	}
	response.Json(response.SUCCESS, gin.H{
		"total":    total,
		"page":     rA.Page,
		"pageSize": rA.PageSize,
		"list":     list,
	}, c)
}

// AdminCreate 新建管理员
func AdminCreate(c *gin.Context) {
	var cA request.CrateAdministrator
	if err := c.ShouldBindBodyWithJSON(&cA); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//对传递过来的密码进行aes解密操作
	ivKey := strconv.FormatInt(cA.Timestamp*1000, 10)
	panPwdBytes, err := util.AesCBCPkcs7Decrypt(cA.Password, config.Data.Server.PasswordSign, ivKey)
	panPwd := string(panPwdBytes)
	if err != nil || panPwd == "" {
		response.Json(response.ERROR, "服务器验证密码错误", c)
		return
	}
	// 设置管理员的角色信息
	var adminRoles []*database.Role
	for _, v := range cA.RoleIDs {
		var role database.Role
		role.ID = v
		adminRoles = append(adminRoles, &role)
	}
	administrator := database.Administrator{
		Mobile:    cA.Mobile,
		Password:  panPwd,
		Email:     cA.Email,
		Name:      cA.Name,
		Avatar:    cA.Avatar,
		Sex:       cA.Sex,
		Status:    1,
		RolesData: adminRoles,
	}
	err = administrator.CreateAdministrator(cA.RoleIDs)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "创建管理员错误", c)
		return
	}

	response.Json(response.SUCCESS, administrator, c)
}

// AdminUpdate 编辑管理员,多个信息
func AdminUpdate(c *gin.Context) {
	var cA response.Auth
	if err := c.ShouldBindBodyWithJSON(&cA); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	if cA.ID == 0 {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	administrator := database.Administrator{
		Mobile: cA.Mobile,
		Email:  cA.Email,
		Name:   cA.Name,
		Avatar: cA.Avatar,
		Sex:    cA.Sex,
		Status: cA.Status,
	}
	administrator.ID = cA.ID

	// 数据库更新
	err := administrator.UpdateAdministrator()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "修改信息失败", c)
		return
	}

	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	//更新redis中管理员信息
	err = redisutil.SetAdminInfoUpdateAtStatus(administrator.ID, appTypeString, administrator.UpdatedAt, administrator.Status)
	if err != nil {
		response.Json(response.AuthError, "更新缓存失败", c)
		return
	}

	response.Json(response.SUCCESS, administrator, c)
}

// AdminResetPwd 重置管理员密码
func AdminResetPwd(c *gin.Context) {
	var pA request.PasswordAdministrator
	if err := c.ShouldBindBodyWithJSON(&pA); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//对传递过来的密码进行aes解密操作
	ivKey := strconv.FormatInt(pA.Timestamp*1000, 10)
	pwdBytes, err := util.AesCBCPkcs7Decrypt(pA.Password, config.Data.Server.PasswordSign, ivKey)
	pwd := string(pwdBytes)
	if err != nil || pwd == "" {
		response.Json(response.ERROR, "服务器验证密码错误", c)
		return
	}

	//修改新密码
	var pwdAdmin database.Administrator
	pwdAdmin.ID = pA.ID
	pwdAdmin.Password = pwd
	// 数据库更新
	err = pwdAdmin.UpdateAdministrator()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "重置密码失败", c)
		return
	}

	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	//更新redis中管理员信息
	err = redisutil.JwtSetAdminPwdUpdateAt(pwdAdmin.ID, appTypeString, pwdAdmin.UpdatedAt)
	if err != nil {
		response.Json(response.ERROR, "重置密码缓存失败", c)
		return
	}

	//判断如果是管理员修改自己的数据，让管理员重新登录
	adminID := c.MustGet("AID").(uint)
	if pA.ID == adminID {
		response.Json(response.AuthPwdUpdate, "当前账号密码已修改，请重新登录！", c)
		return
	}
	response.Json(response.SUCCESS, "success", c)
}

// AdminDelete 删除管理员
func AdminDelete(c *gin.Context) {
	var reqData request.BaseID
	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	appTypeString := c.GetHeader("Type")
	// 删除redis中存储的用户信息
	err := redisutil.DeleteAdminInfo(reqData.ID, appTypeString)
	if err != nil {
		response.Json(response.ERROR, "删除缓存失败", c)
		return
	}
	//永久删除管理员信息
	var admin database.Administrator
	admin.ID = reqData.ID
	err = admin.UnscopedDeleteAdministrator()
	if err != nil {
		response.Json(response.ERROR, "删除管理员失败", c)
		return
	}
	response.Json(response.SUCCESS, "success", c)
}

// AdminSetRole 设置管理员角色
func AdminSetRole(c *gin.Context) {
	var rSR request.SetAdministratorRole
	if err := c.ShouldBindBodyWithJSON(&rSR); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	if rSR.ID == 0 || len(rSR.RoleIDs) == 0 {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//更新用户角色中间表
	err := database.UpdateAdministratorRoles(rSR.ID, rSR.RoleIDs)
	if err != nil {
		response.Json(response.ERROR, "更新用户角色失败", c)
		return
	}
	//从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	appTypeString := c.GetHeader("Type")
	//更新redis中管理员信息
	err = redisutil.SetAdminInfoRoleIDs(rSR.ID, appTypeString, rSR.RoleIDs)
	if err != nil {
		response.Json(response.AuthError, "更新缓存失败", c)
		return
	}
	//判断如果是管理员修改自己的数据，让管理员重新登录
	adminID := c.MustGet("AID").(uint)
	if rSR.ID == adminID {
		response.Json(response.AuthRoleUpdate, "账户角色信息有更新，请重新登录！", c)
		return
	}
	response.Json(response.SUCCESS, "success", c)
}
