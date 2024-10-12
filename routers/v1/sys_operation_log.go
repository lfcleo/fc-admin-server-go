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

// LogList 操作日志列表
func LogList(c *gin.Context) {
	var rO request.OperationLog
	if err := c.ShouldBindBodyWithJSON(&rO); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	queryStr := rO.BuildQueryConditions()
	list, total, err := database.QueryOperationLogList(queryStr, rO.Page, rO.PageSize, "Administrator")
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取列表信息失败", c)
		return
	}
	response.Json(response.SUCCESS, gin.H{
		"total":    total,
		"page":     rO.Page,
		"pageSize": rO.PageSize,
		"list":     list,
	}, c)
}

// LogDelete 删除操作日志
func LogDelete(c *gin.Context) {
	var reqData request.BaseIDs
	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	err := database.UnscopedDeleteSysOperationLog(reqData.IDs)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "删除错误", c)
		return
	}

	response.Json(response.SUCCESS, "success", c)
}
