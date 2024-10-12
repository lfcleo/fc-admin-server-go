package v1

import (
	"errors"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/database"
	"fc-admin-server-go/models/request"
	"fc-admin-server-go/models/response"
	casbinUtil "fc-admin-server-go/pkg/casbin"
	redisUtil "fc-admin-server-go/pkg/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

// RoleList 角色列表
func RoleList(c *gin.Context) {
	var rR request.Role
	if err := c.ShouldBindBodyWithJSON(&rR); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	queryStr := rR.BuildQueryConditions()
	list, total, err := database.QueryRoleList(queryStr, rR.Page, rR.PageSize)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取列表信息失败", c)
		return
	}
	response.Json(response.SUCCESS, gin.H{
		"total":    total,
		"page":     rR.Page,
		"pageSize": rR.PageSize,
		"list":     list,
	}, c)
}

// RoleCreate 新建角色
func RoleCreate(c *gin.Context) {
	var cR request.CreateRole
	if err := c.ShouldBindBodyWithJSON(&cR); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	adminID := c.MustGet("AID").(uint)
	role := database.Role{
		Name:            cR.Name,
		Code:            cR.Code,
		Notes:           cR.Notes,
		AdministratorID: adminID,
	}
	err := role.CreateRole()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "创建管理员错误", c)
		return
	}

	//角色更新时间信息存储redis中
	err = redisUtil.AddRoleInfo(role.ID, role.UpdatedAt)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储角色信息错误", c)
		return
	}

	response.Json(response.SUCCESS, role, c)
}

// RoleUpdate 更新角色
func RoleUpdate(c *gin.Context) {
	var cR request.CreateRole
	if err := c.ShouldBindBodyWithJSON(&cR); err != nil || cR.ID == 0 {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	role := database.Role{
		Name:  cR.Name,
		Code:  cR.Code,
		Notes: cR.Notes,
	}
	role.ID = cR.ID

	// 数据库更新
	err := role.UpdateRole()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "修改信息失败", c)
		return
	}

	//角色更新时间信息存储redis中
	err = redisUtil.AddRoleInfo(role.ID, role.UpdatedAt)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储角色信息错误", c)
		return
	}

	response.Json(response.SUCCESS, role, c)
}

// RoleApis 获取角色的API权限列表
func RoleApis(c *gin.Context) {
	var rID request.BaseID
	if err := c.ShouldBindBodyWithJSON(&rID); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	apis, err := casbinUtil.CasbinServiceApp.GetPolicyPathByRoleID(rID.ID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取管理员API列表权限错误", c)
		return
	}

	response.Json(response.SUCCESS, apis, c)
}

// SetApis 设置角色的API接口权限列表
func SetApis(c *gin.Context) {
	var rSA request.SetApis
	if err := c.ShouldBindBodyWithJSON(&rSA); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	if rSA.ID == 0 || len(rSA.Apis) == 0 {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//更新角色API接口中间表
	tTime, err := database.UpdateRoleApis(rSA.ID, rSA.Apis)
	if err != nil {
		response.Json(response.ERROR, "更新角色API接口列表失败", c)
		return
	}
	//设置redis中角色信息
	err = redisUtil.AddRoleInfo(rSA.ID, tTime)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储管理员信息错误", c)
		return
	}
	response.Json(response.SUCCESS, "success", c)
}

// RoleMenus 获取角色的菜单列表
func RoleMenus(c *gin.Context) {
	var rID request.BaseID
	if err := c.ShouldBindBodyWithJSON(&rID); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	var role database.Role
	role.ID = rID.ID

	menus, err := database.QueryRoleMenus(&role)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取管理员菜单列表权限错误", c)
		return
	}
	var menuIDs []uint
	for _, menu := range menus {
		menuIDs = append(menuIDs, menu.ID)
	}

	response.Json(response.SUCCESS, menuIDs, c)
}

// SetMenus 设置角色的菜单列表
func SetMenus(c *gin.Context) {
	var rSA request.SetMenuIDs
	if err := c.ShouldBindBodyWithJSON(&rSA); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	if rSA.ID == 0 || len(rSA.IDs) == 0 {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	//更新角色API接口中间表
	tTime, err := database.UpdateRoleMenus(rSA.ID, rSA.IDs)
	if err != nil {
		response.Json(response.ERROR, "更新角色菜单列表失败", c)
		return
	}
	//设置redis中角色信息
	err = redisUtil.AddRoleInfo(rSA.ID, tTime)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储管理员信息错误", c)
		return
	}
	response.Json(response.SUCCESS, "success", c)
}

// RoleDelete 删除角色信息
func RoleDelete(c *gin.Context) {
	var reqData request.BaseID
	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	var role database.Role
	role.ID = reqData.ID
	err := role.UnscopedDeleteRole()
	if err != nil {
		response.Json(response.ERROR, "删除角色失败", c)
		return
	}
	// 设置redis中角色更新时间为 0 值，因为此数据已经删除了
	err = redisUtil.AddRoleInfo(reqData.ID, time.Time{})
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储管理员信息错误", c)
		return
	}

	response.Json(response.SUCCESS, "success", c)
}
