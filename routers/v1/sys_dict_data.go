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

// DictDataList 字典数据列表
func DictDataList(c *gin.Context) {
	var dD request.DictData
	if err := c.ShouldBindBodyWithJSON(&dD); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	queryStr := dD.BuildQueryConditions()
	list, total, err := database.QueryDictDataList(queryStr, dD.Page, dD.PageSize)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取列表信息失败", c)
		return
	}
	response.Json(response.SUCCESS, gin.H{
		"total":    total,
		"page":     dD.Page,
		"pageSize": dD.PageSize,
		"list":     list,
	}, c)
}

// DictDataCreate 创建字典数据
func DictDataCreate(c *gin.Context) {
	var sDD database.SysDictData
	if err := c.ShouldBindBodyWithJSON(&sDD); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	err := sDD.CreateDictData()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "创建字典分类错误", c)
		return
	}

	response.Json(response.SUCCESS, sDD, c)
}

// DictDataUpdate 更新字典数据
func DictDataUpdate(c *gin.Context) {
	var sDD database.SysDictData
	if err := c.ShouldBindBodyWithJSON(&sDD); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	err := sDD.UpdateDictData()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "更新字典数据错误", c)
		return
	}

	response.Json(response.SUCCESS, sDD, c)
}

// DictDataDelete 删除字典数据
func DictDataDelete(c *gin.Context) {
	var reqData request.BaseIDs
	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	err := database.UnscopedDeleteSysDictData(reqData.IDs)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "删除错误", c)
		return
	}

	response.Json(response.SUCCESS, "success", c)
}
