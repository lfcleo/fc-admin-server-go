package v1

import (
	"fc-admin-server-go/models/response"
	redisutil "fc-admin-server-go/pkg/redis"
	"github.com/gin-gonic/gin"
	"log"
)

// MenuApiList 动态路由及权限组
func MenuApiList(c *gin.Context) {
	response.Json(response.SUCCESS, gin.H{
		"menu": []gin.H{
			{
				"name":      "dashboard",
				"path":      "/dashboard",
				"component": "home/index",
				"meta": gin.H{
					"icon":    "Eleme",
					"title":   "控制台",
					"type":    "menu",
					"isAffix": true,
				},
			},
			{
				"name":     "setting",
				"path":     "/setting",
				"redirect": "/setting/menu",
				"meta": gin.H{
					"icon":  "Setting",
					"title": "配置",
					"type":  "menu",
				},
				"children": []gin.H{
					{
						"name":      "settingMenu",
						"path":      "/setting/menu",
						"component": "setting/menu/index",
						"meta": gin.H{
							"title": "菜单管理",
							"icon":  "Menu",
							"type":  "menu",
						},
						"children": []gin.H{
							{
								"name":      "test",
								"path":      "",
								"component": "",
								"meta": gin.H{
									"title": "测试按钮",
									"type":  "BUTTON",
								},
							},
						},
					},
					{
						"name":      "settingApi",
						"path":      "/setting/api",
						"component": "setting/api/index",
						"meta": gin.H{
							"title": "接口管理",
							"icon":  "Platform",
							"type":  "menu",
						},
					},
					{
						"name":      "settingRole",
						"path":      "/setting/role",
						"component": "setting/role/index",
						"meta": gin.H{
							"title": "角色管理",
							"icon":  "Avatar",
							"type":  "menu",
						},
					},
					{
						"name":      "settingAdmin",
						"path":      "/setting/admin",
						"component": "setting/admin/index",
						"meta": gin.H{
							"title":       "用户管理",
							"icon":        "UserFilled",
							"type":        "menu",
							"isKeepAlive": true,
						},
					},
					{
						"name":      "settingDict",
						"path":      "/setting/dict",
						"component": "setting/dict/index",
						"meta": gin.H{
							"title": "字典管理",
							"icon":  "Memo",
							"type":  "menu",
						},
					},
					{
						"name":      "link",
						"path":      "https://www.baidu.com",
						"component": "",
						"meta": gin.H{
							"icon":  "Link",
							"title": "外部链接",
							"type":  "LINK",
						},
					},
				},
			},
			{
				"name":      "link",
				"path":      "https://www.baidu.com",
				"component": "",
				"meta": gin.H{
					"icon":  "Link",
					"title": "外部链接",
					"type":  "LINK",
				},
			},
		},
		"permissions": []string{"ALL"},
	}, c)
}

func Test(c *gin.Context) {
	roles, err := redisutil.GetAllRoleInfos()
	if err != nil {
		log.Println(err)
	}
	log.Println(roles)
	response.Json(response.SUCCESS, "success", c)
}
