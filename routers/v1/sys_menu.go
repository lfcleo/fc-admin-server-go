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

// MenuList 菜单列表
func MenuList(c *gin.Context) {
	var rM request.Menu
	if err := c.ShouldBindBodyWithJSON(&rM); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	queryStr := rM.BuildQueryConditions()
	list, err := database.QueryMenuList(queryStr)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "获取列表信息失败", c)
		return
	}
	//treeList, _ := database.MenuTree(list, 0, false)
	treeList := database.BuildMenuTree(list)
	response.Json(response.SUCCESS, treeList, c)
}

// MenuCreate 创建菜单
func MenuCreate(c *gin.Context) {
	var cM database.Menu
	if err := c.ShouldBindBodyWithJSON(&cM); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	adminID := c.MustGet("AID").(uint)
	cM.AdministratorID = adminID
	err := cM.CreateMenu()
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "创建菜单错误", c)
		return
	}

	err = redisUtil.AddRoleInfos([]uint{1}, cM.UpdatedAt)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储角色信息错误", c)
		return
	}

	response.Json(response.SUCCESS, cM, c)
}

// MenuUpdate 更新菜单信息
func MenuUpdate(c *gin.Context) {
	var cM database.Menu
	if err := c.ShouldBindBodyWithJSON(&cM); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}
	if cM.ID == 0 {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	// 数据库更新
	roleIDs, err := database.UpdateMenu(&cM)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "修改信息失败", c)
		return
	}
	//if len(roleIDs) > 0 {
	//	// 设置redis中角色更新时间
	//	for _, id := range roleIDs {
	//		err = redisUtil.AddRoleInfo(id, cM.UpdatedAt)
	//		if err != nil {
	//			global.FC_LOGGER.Error(fmt.Sprint(err))
	//			response.Json(response.ERROR, "redis存储管理员信息错误", c)
	//			return
	//		}
	//	}
	//}
	err = redisUtil.AddRoleInfos(roleIDs, cM.UpdatedAt)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储角色信息错误", c)
		return
	}

	response.Json(response.SUCCESS, cM, c)
}

// MenuDelete 删除菜单
func MenuDelete(c *gin.Context) {
	var reqData request.BaseID
	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		response.Json(response.ERROR, "参数错误", c)
		return
	}

	//判断菜单是否有子菜单，有的话不允许删除
	count, err := database.QueryMenuCountByParentID(reqData.ID)
	if err != nil {
		response.Json(response.ERROR, "删除菜单失败", c)
		return
	}
	if count > 0 {
		response.Json(response.ERROR, "删除失败，当前菜单目录存在子菜单，不允许删除！", c)
		return
	}
	var menu database.Menu
	menu.ID = reqData.ID
	roleIDs, tTime, err := database.UnscopedDeleteMenu(&menu)
	if err != nil {
		response.Json(response.ERROR, "删除菜单失败", c)
		return
	}
	//if len(roleIDs) > 0 {
	//	// 设置redis中角色更新时间
	//	for _, id := range roleIDs {
	//		err = redisUtil.AddRoleInfo(id, tTime)
	//		if err != nil {
	//			global.FC_LOGGER.Error(fmt.Sprint(err))
	//			response.Json(response.ERROR, "redis存储管理员信息错误", c)
	//			return
	//		}
	//	}
	//}
	err = redisUtil.AddRoleInfos(roleIDs, tTime)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		response.Json(response.ERROR, "redis存储角色信息错误", c)
		return
	}

	response.Json(response.SUCCESS, "success", c)
}
