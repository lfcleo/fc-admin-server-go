package v1

import (
	"errors"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/database"
	"fc-admin-server-go/models/request"
	"fc-admin-server-go/models/response"
	redisUtil "fc-admin-server-go/pkg/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApiList 接口列表
func ApiList(c *gin.Context) {
	var rA request.Api
	if err := c.ShouldBindBodyWithJSON(&rA); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	queryStr := rA.BuildQueryConditions()
	list, total, err := database.QueryApiList(queryStr, rA.Page, rA.PageSize)
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

// ApiAllList 所有api接口列表
func ApiAllList(c *gin.Context) {
	list, err := database.QueryAllApiList()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取列表信息失败", c)
		return
	}
	response.Json(response.SUCCESS, list, c)
}

// ApiCreate 创建接口
func ApiCreate(c *gin.Context) {
	var cA request.CreateApi
	if err := c.ShouldBindBodyWithJSON(&cA); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	adminID := c.MustGet("AID").(uint)
	api := database.Api{
		Path:            cA.Path,
		Method:          cA.Method,
		Description:     cA.Description,
		AdministratorID: adminID,
	}
	err := api.CreateApi()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "创建管理员错误", c)
		return
	}

	response.Json(response.SUCCESS, api, c)
}

// ApiUpdate 更新接口
func ApiUpdate(c *gin.Context) {
	var cA request.CreateApi
	if err := c.ShouldBindBodyWithJSON(&cA); err != nil || cA.ID == 0 {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	api := database.Api{
		Path:        cA.Path,
		Method:      cA.Method,
		Description: cA.Description,
	}
	api.ID = cA.ID

	// 数据库更新
	roleIDs, err := database.UpdateApi(&api)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "修改信息失败", c)
		return
	}

	//if len(roleIDs) > 0 {
	//	// 设置redis中角色更新时间
	//	for _, id := range roleIDs {
	//		err = redisUtil.AddRoleInfo(id, api.UpdatedAt)
	//		if err != nil {
	//			global.FC_LOGGER.Error(fmt.Sprint(err))
	//			response.Json(response.ERROR, "redis存储管理员信息错误", c)
	//			return
	//		}
	//	}
	//}
	err = redisUtil.AddRoleInfos(roleIDs, api.UpdatedAt)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储管理员信息错误", c)
		return
	}

	response.Json(response.SUCCESS, api, c)
}

// ApiDelete 删除api接口权限
func ApiDelete(c *gin.Context) {
	var reqData request.BaseID
	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	var api database.Api
	api.ID = reqData.ID
	roleIDs, tTime, err := database.UnscopedDeleteApi(&api)
	if err != nil {
		response.Json(response.ERROR, "删除接口信息失败", c)
		return
	}
	//if len(roleIDs) > 0 {
	// 设置redis中角色更新时间
	//for _, id := range roleIDs {
	//	err = redisUtil.AddRoleInfo(id, tTime)
	//	if err != nil {
	//		global.FC_LOGGER.Error(fmt.Sprint(err))
	//		response.Json(response.ERROR, "redis存储管理员信息错误", c)
	//		return
	//	}
	//}
	//}
	err = redisUtil.AddRoleInfos(roleIDs, tTime)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储管理员信息错误", c)
		return
	}
	response.Json(response.SUCCESS, "success", c)
}
