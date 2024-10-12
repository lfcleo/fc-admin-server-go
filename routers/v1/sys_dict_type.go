package v1

import (
	"errors"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/database"
	"fc-admin-server-go/models/request"
	"fc-admin-server-go/models/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DictTypeList 字典分类列表
func DictTypeList(c *gin.Context) {
	list, err := database.QueryDictTypeList()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取列表信息失败", c)
		return
	}
	response.Json(response.SUCCESS, list, c)
}

// DictTypeCreate 创建字典分类
func DictTypeCreate(c *gin.Context) {
	var sDT database.SysDictType
	if err := c.ShouldBindBodyWithJSON(&sDT); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	adminID := c.MustGet("AID").(uint)
	sDT.AdministratorID = adminID
	err := sDT.CreateDictType()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "创建字典分类错误", c)
		return
	}

	response.Json(response.SUCCESS, sDT, c)
}

// DictTypeUpdate 更新字典分类
func DictTypeUpdate(c *gin.Context) {
	var sDT database.SysDictType
	if err := c.ShouldBindBodyWithJSON(&sDT); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	err := database.UpdateDictType(&sDT)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "创建字典分类错误", c)
		return
	}

	response.Json(response.SUCCESS, sDT, c)
}

// DictTypeDelete 删除字典分类
func DictTypeDelete(c *gin.Context) {
	var reqData request.BaseID
	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	var dt database.SysDictType
	dt.ID = reqData.ID
	err := database.UnscopedDeleteSysDictType(&dt)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "删除字典分类错误", c)
		return
	}

	response.Json(response.SUCCESS, "success", c)
}
