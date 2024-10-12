package routers

import (
	"fc-admin-server-go/middleware"
	"fc-admin-server-go/pkg/config"
	v1 "fc-admin-server-go/routers/v1"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New()           //创建gin框架路由实例
	r.Use(middleware.Cors()) //跨域中间件
	r.Use(gin.Logger())      //使用gin框架中的打印中间件
	r.Use(gin.Recovery())    //使用gin框架中的恢复中间件，可以从任何恐慌中恢复，如果有，则写入500

	mode := "release"
	if config.Data.Debug {
		mode = "debug"
	}
	gin.SetMode(mode) //设置运行模式，debug或release,如果放在gin.New或者gin.Default之后，还是会打印一些信息的。放之前则不会

	//客户端API接口
	apiV1 := r.Group("/v1/")             //路由分组，apiV1代表v1版本的路由组
	apiV1.Use(middleware.VerEncrypt())   //使用校验接口签名的中间件
	apiV1.Use(middleware.OperationLog()) //记录操作日志
	{
		apiV1.POST("auth/login", v1.Login)                        //手机号/邮箱 密码登录
		apiV1.POST("auth/verification_code", v1.VerificationCode) //手机号发送验证码
		apiV1.POST("auth/login/mobile", v1.MobileLogin)           //手机号登录
		apiV1Token := apiV1.Group("token/")                       //创建使用token中间件的路由组
		apiV1Token.Use(middleware.VerToken())
		{
			apiV1Token.POST("test", v1.Test)     //测试接口
			authAPI := apiV1Token.Group("auth/") //用户管理员鉴权路由组
			authAPI.Use(middleware.VerCache())   //用户角色信息授权
			{
				authAPI.POST("refresh", v1.RefreshToken)        //刷新token
				authAPI.POST("menu", v1.MenuPermissionList)     //动态路由及权限组
				authAPI.POST("update/info", v1.UpdateAdminInfo) //更改管理员信息
				authAPI.POST("update/pwd", v1.UpdateAdminPwd)   //更改管理员密码
				authAPI.POST("info", v1.AdminInfo)              //获取管理员信息
				authAPI.POST("dict", v1.DictInfo)               //获取字典数据
				authAPI.POST("upload", v1.UploadFile)           //上传文件（图片等）
			}
			apiV1Token.POST("auth/logout", v1.Logout) //退出登录(不参与用户角色信息授权)

			adminAPI := apiV1Token.Group("admin/") //管理员管理路由组
			adminAPI.Use(middleware.VerCache())    //用户角色信息授权
			{
				adminAPI.POST("list", v1.AdminList)          //管理员列表
				adminAPI.POST("create", v1.AdminCreate)      //新建管理员
				adminAPI.POST("update", v1.AdminUpdate)      //编辑管理员
				adminAPI.POST("reset/pwd", v1.AdminResetPwd) //重置管理员密码
				adminAPI.POST("delete", v1.AdminDelete)      //删除管理员
				adminAPI.POST("set/role", v1.AdminSetRole)   //设置管理员角色
			}
			roleAPI := apiV1Token.Group("role/") //角色路由组
			roleAPI.Use(middleware.VerCache())   //用户角色信息授权
			{
				roleAPI.POST("list", v1.RoleList)      //角色列表
				roleAPI.POST("create", v1.RoleCreate)  //新建角色
				roleAPI.POST("update", v1.RoleUpdate)  //编辑角色
				roleAPI.POST("apis", v1.RoleApis)      //获取角色的API权限列表
				roleAPI.POST("set/apis", v1.SetApis)   //设置角色的API权限列表
				roleAPI.POST("menus", v1.RoleMenus)    //获取角色的菜单列表
				roleAPI.POST("set/menus", v1.SetMenus) //设置角色的菜单列表
				roleAPI.POST("delete", v1.RoleDelete)  //删除角色
			}
			menuAPI := apiV1Token.Group("menu/") //菜单路由组
			menuAPI.Use(middleware.VerCache())   //用户角色信息授权
			{
				menuAPI.POST("list", v1.MenuList) //菜单列表
			}
			apiV1Token.POST("menu/create", v1.MenuCreate) //新建菜单(不参与用户角色信息授权)
			apiV1Token.POST("menu/update", v1.MenuUpdate) //编辑菜单(不参与用户角色信息授权)
			apiV1Token.POST("menu/delete", v1.MenuDelete) //删除菜单(不参与用户角色信息授权)

			apisAPI := apiV1Token.Group("api/") //接口路由组
			apisAPI.Use(middleware.VerCache())  //用户角色信息授权
			{
				apisAPI.POST("list", v1.ApiList)        //接口列表
				apisAPI.POST("create", v1.ApiCreate)    //新建接口
				apisAPI.POST("update", v1.ApiUpdate)    //编辑接口
				apisAPI.POST("delete", v1.ApiDelete)    //删除接口
				apisAPI.POST("list/all", v1.ApiAllList) //所有接口列表
			}
			dictAPI := apiV1Token.Group("dict/") //字典路由组
			dictAPI.Use(middleware.VerCache())   //用户角色信息授权
			{
				dictAPI.POST("type/list", v1.DictTypeList)     //字典分类列表
				dictAPI.POST("type/create", v1.DictTypeCreate) //新建字典分类
				dictAPI.POST("type/update", v1.DictTypeUpdate) //更新字典分类
				dictAPI.POST("type/delete", v1.DictTypeDelete) //删除字典分类
				dictAPI.POST("data/list", v1.DictDataList)     //字典数据列表
				dictAPI.POST("data/create", v1.DictDataCreate) //新建字典数据
				dictAPI.POST("data/update", v1.DictDataUpdate) //更新字典数据
				dictAPI.POST("data/delete", v1.DictDataDelete) //删除字典数据
			}
			logsAPI := apiV1Token.Group("log/") //操作日志路由组
			logsAPI.Use(middleware.VerCache())  //用户角色信息授权
			{
				logsAPI.POST("list", v1.LogList)     //操作日志列表
				logsAPI.POST("delete", v1.LogDelete) //删除日志（多选）
			}
		}
	}

	//当访问 $HOST/upload/apks 时，将会读取到 项目/runtime/upload/apks 下的文件 这样就能让外部访问到图片资源了
	r.Static("/documents", "./runtime/documents")
	//静态资源与静态html文件
	r.Static("/web", "./web")
	return r
}
