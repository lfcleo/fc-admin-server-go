package initialize

import (
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/database"
	"fc-admin-server-go/pkg/config"
	"fc-admin-server-go/pkg/util"
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

// BaseModel gorm.Model 的定义
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type StringArray []string // 字符串数组类型
type UintArray []uint     // uint数组类型

func GormSetUp() *gorm.DB {
	var (
		err error
		//databaseType = setting.DatabaseSetting.Type     //数据库类型
		user = config.Data.Database.User     //数据库的用户
		pass = config.Data.Database.Password //数据库的密码
		host = config.Data.Database.Host     //数据库地址
		name = config.Data.Database.Name     //数据库名称
	)

	//使用gorm链接数据库
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, name)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   config.Data.Database.TablePrefix, //设置表名称的前缀,表名前缀，需要的话解除此注释
			SingularTable: true, //使用单数表名，启用该选项，此时
		},
		DisableForeignKeyConstraintWhenMigrating: true, //自动迁移时，禁用外键约束,不禁用
		//Logger:                                   logger.Default.LogMode(logger.Info), //配置日志级别，打印出所有的sql
	})
	if err != nil {
		global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】链接失败：%v", err)) //数据库链接失败是致命的错误，链接失败后可以关闭程序了，所以使用logging.Fatal方法
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(config.Data.Database.MaxIdleConnects) //设置空闲时的最大连接数
	sqlDB.SetMaxOpenConns(config.Data.Database.MaxOpenConnects) //设置数据库的最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)                         //链接池中链接的最大可复用时间

	//自动检查 Tag 结构是否变化，变化则进行迁移，需要的参数为数据库模型结构体
	err = db.AutoMigrate(&database.Administrator{}, &database.Menu{}, &database.Role{}, &database.AdministratorRole{}, &database.Api{},
		&database.Menu{}, &database.RoleMenu{}, &gormadapter.CasbinRule{}, &database.SysDictType{}, &database.SysDictData{}, &database.SysOperationLog{})
	if err != nil {
		global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】迁移失败：%v", err)) //数据库链接失败是致命的错误，链接失败后可以关闭程序了，所以使用logging.Fatal方法
	}

	return db
}

// InitSQLData 创建默认数据库数据
func InitSQLData() {
	//自动创建超级管理员,包括超级管理员权限，关联表
	var count int64
	err := global.FC_DB.Model(&database.Administrator{}).Count(&count).Error
	if err != nil {
		global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】查询默认超级管理员数据库失败：%v", err)) //链接失败后可以关闭程序了，所以使用logging.Fatal方法
	}
	if count == 0 {
		// 如果没有数据，则插入默认数据
		err = global.FC_DB.Create(&database.Administrator{
			Mobile:   "18888888888",
			Password: util.EncodeSha1("Qwe123"),
			Email:    "admin@admin.com",
			Name:     "超级管理员",
			RolesData: []*database.Role{
				{
					Name:            "超级管理员",
					Code:            "ALL",
					Notes:           "最高权限，不可修改。",
					AdministratorID: 1,
				},
			},
		}).Error
		if err != nil {
			global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】创建默认超级管理员失败：%v", err)) //链接失败后可以关闭程序了，所以使用logging.Fatal方法
		}
	}

	//自动创建api接口
	err = global.FC_DB.Model(&database.Api{}).Count(&count).Error
	if err != nil {
		global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】查询默认API接口数据库失败：%v", err)) //链接失败后可以关闭程序了，所以使用logging.Fatal方法
	}
	if count == 0 {
		entities := []database.Api{
			//{Path: "auth/refresh", Method: "POST", Description: "刷新token(必选)", SystemApi: true, AdministratorID: 1},
			//
			//{Path: "auth/menu", Method: "POST", Description: "动态路由及权限组(必选)", SystemApi: true, AdministratorID: 1},
			//{Path: "auth/logout", Method: "POST", Description: "退出登录(必选)", SystemApi: true, AdministratorID: 1},
			//{Path: "auth/info", Method: "POST", Description: "获取管理员信息(必选)", SystemApi: true, AdministratorID: 1},
			//{Path: "auth/update/info", Method: "POST", Description: "更改管理员信息(建议选)", AdministratorID: 1},
			//{Path: "auth/update/pwd", Method: "POST", Description: "更改管理员密码(建议选)", AdministratorID: 1},
			//{Path: "auth/upload", Method: "POST", Description: "上传文件/图片(建议选)", AdministratorID: 1},

			{Path: "admin/list", Method: "POST", Description: "管理员列表", AdministratorID: 1},
			{Path: "admin/create", Method: "POST", Description: "新建管理员", AdministratorID: 1},
			{Path: "admin/update", Method: "POST", Description: "编辑管理员", AdministratorID: 1},
			{Path: "admin/reset/pwd", Method: "POST", Description: "重置管理员密码", AdministratorID: 1},
			{Path: "admin/delete", Method: "POST", Description: "删除管理员", AdministratorID: 1},
			{Path: "admin/set/role", Method: "POST", Description: "设置管理员角色", AdministratorID: 1},

			{Path: "role/list", Method: "POST", Description: "角色列表", AdministratorID: 1},
			{Path: "role/create", Method: "POST", Description: "新建角色", AdministratorID: 1},
			{Path: "role/update", Method: "POST", Description: "编辑角色", AdministratorID: 1},
			{Path: "role/apis", Method: "POST", Description: "获取角色的API权限列表", AdministratorID: 1},
			{Path: "role/set/apis", Method: "POST", Description: "设置角色的API权限列表", AdministratorID: 1},
			{Path: "role/menus", Method: "POST", Description: "获取角色的菜单列表", AdministratorID: 1},
			{Path: "role/set/menus", Method: "POST", Description: "设置角色的菜单列表", AdministratorID: 1},
			{Path: "role/delete", Method: "POST", Description: "删除角色", AdministratorID: 1},

			{Path: "menu/list", Method: "POST", Description: "菜单列表", AdministratorID: 1},
			{Path: "menu/create", Method: "POST", Description: "新建菜单", AdministratorID: 1},
			{Path: "menu/update", Method: "POST", Description: "编辑菜单", AdministratorID: 1},
			{Path: "menu/delete", Method: "POST", Description: "删除菜单", AdministratorID: 1},

			{Path: "api/list", Method: "POST", Description: "接口列表", AdministratorID: 1},
			{Path: "api/create", Method: "POST", Description: "新建接口", AdministratorID: 1},
			{Path: "api/update", Method: "POST", Description: "编辑接口", AdministratorID: 1},
			{Path: "api/delete", Method: "POST", Description: "删除接口", AdministratorID: 1},
			{Path: "api/list/all", Method: "POST", Description: "全部接口列表", AdministratorID: 1},

			{Path: "dict/type/list", Method: "POST", Description: "字典分类列表", AdministratorID: 1},
			{Path: "dict/type/create", Method: "POST", Description: "新建字典分类", AdministratorID: 1},
			{Path: "dict/type/update", Method: "POST", Description: "编辑字典分类", AdministratorID: 1},
			{Path: "dict/type/delete", Method: "POST", Description: "删除字典分类", AdministratorID: 1},
			{Path: "dict/data/list", Method: "POST", Description: "字典数据列表", AdministratorID: 1},
			{Path: "dict/data/create", Method: "POST", Description: "新建字典数据", AdministratorID: 1},
			{Path: "dict/data/update", Method: "POST", Description: "编辑字典数据", AdministratorID: 1},
			{Path: "dict/data/delete", Method: "POST", Description: "删除字典数据", AdministratorID: 1},

			{Path: "log/list", Method: "POST", Description: "操作日志列表", AdministratorID: 1},
			{Path: "log/delete", Method: "POST", Description: "删除操作日志", AdministratorID: 1},
		}
		if err := global.FC_DB.Create(&entities).Error; err != nil {
			global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】创建默认API接口失败：%v", err)) //链接失败后可以关闭程序了，所以使用logging.Fatal方法
		}
	}

	//自动创建菜单接口
	err = global.FC_DB.Model(&database.Menu{}).Count(&count).Error
	if err != nil {
		global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】查询默认菜单数据库失败：%v", err)) //链接失败后可以关闭程序了，所以使用logging.Fatal方法
	}
	if count == 0 {
		entities := []database.Menu{
			{SystemMenu: true, Sort: 1, Path: "/dashboard", Name: "dashboard", Component: "home/index", AdministratorID: 1, Meta: database.MenuMeta{
				Icon:    "Eleme",
				Title:   "控制台",
				Type:    "MENU",
				IsAffix: true,
			}},
			{SystemMenu: true, Sort: 2, Path: "/setting", Name: "setting", Redirect: "/setting/menu", Component: "setting", AdministratorID: 1, Meta: database.MenuMeta{
				Icon:  "Setting",
				Title: "配置",
				Type:  "MENU",
			}},
			{ParentID: 2, SystemMenu: true, Sort: 1, Name: "settingMenu", Path: "/setting/menu", Component: "setting/menu/index", AdministratorID: 1, Meta: database.MenuMeta{
				Icon:        "Menu",
				Title:       "菜单管理",
				Type:        "MENU",
				IsKeepAlive: true,
			}},
			{ParentID: 3, SystemMenu: false, Sort: 1, Name: "menu:create:root", Path: "menu:create:root", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "新增根菜单",
				Type:  "BUTTON",
			}},
			{ParentID: 3, SystemMenu: false, Sort: 2, Name: "menu:create:sub", Path: "menu:create:sub", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "添加子菜单",
				Type:  "BUTTON",
			}},
			{ParentID: 3, SystemMenu: false, Sort: 3, Name: "menu:update", Path: "menu:update", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "编辑菜单",
				Type:  "BUTTON",
			}},
			{ParentID: 3, SystemMenu: false, Sort: 4, Name: "menu:delete", Path: "menu:delete", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "删除菜单",
				Type:  "BUTTON",
			}},
			{ParentID: 3, SystemMenu: false, Sort: 5, Name: "menu:read", Path: "menu:read", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "查询菜单",
				Type:  "BUTTON",
			}},
			{ParentID: 2, SystemMenu: true, Sort: 2, Name: "settingApi", Path: "/setting/api", Component: "setting/api/index", AdministratorID: 1, Meta: database.MenuMeta{
				Icon:        "Platform",
				Title:       "接口管理",
				Type:        "MENU",
				IsKeepAlive: true,
			}},
			{ParentID: 9, SystemMenu: false, Sort: 1, Name: "api:create", Path: "api:create", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "新增API接口",
				Type:  "BUTTON",
			}},
			{ParentID: 9, SystemMenu: false, Sort: 2, Name: "api:update", Path: "api:update", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "编辑API接口",
				Type:  "BUTTON",
			}},
			{ParentID: 9, SystemMenu: false, Sort: 3, Name: "api:delete", Path: "api:delete", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "删除API接口",
				Type:  "BUTTON",
			}},
			{ParentID: 9, SystemMenu: false, Sort: 4, Name: "api:read", Path: "api:read", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "查询API接口",
				Type:  "BUTTON",
			}},
			{ParentID: 2, SystemMenu: true, Sort: 3, Name: "settingRole", Path: "/setting/role", Component: "setting/role/index", AdministratorID: 1, Meta: database.MenuMeta{
				Icon:        "Avatar",
				Title:       "角色管理",
				Type:        "MENU",
				IsKeepAlive: true,
			}},
			{ParentID: 14, SystemMenu: false, Sort: 1, Name: "role:create", Path: "role:create", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "新增角色",
				Type:  "BUTTON",
			}},
			{ParentID: 14, SystemMenu: false, Sort: 2, Name: "role:update", Path: "role:update", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "编辑角色",
				Type:  "BUTTON",
			}},
			{ParentID: 14, SystemMenu: false, Sort: 3, Name: "role:delete", Path: "role:delete", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "删除角色",
				Type:  "BUTTON",
			}},
			{ParentID: 14, SystemMenu: false, Sort: 4, Name: "role:read", Path: "role:read", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "查询角色",
				Type:  "BUTTON",
			}},
			{ParentID: 14, SystemMenu: false, Sort: 5, Name: "role:set", Path: "role:set", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "设置权限",
				Type:  "BUTTON",
			}},
			{ParentID: 2, SystemMenu: true, Sort: 4, Name: "settingAdmin", Path: "/setting/admin", Component: "setting/admin/index", AdministratorID: 1, Meta: database.MenuMeta{
				Icon:        "UserFilled",
				Title:       "用户管理",
				Type:        "MENU",
				IsKeepAlive: true,
			}},
			{ParentID: 20, SystemMenu: false, Sort: 1, Name: "admin:create", Path: "admin:create", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "新增管理员",
				Type:  "BUTTON",
			}},
			{ParentID: 20, SystemMenu: false, Sort: 2, Name: "admin:update", Path: "admin:update", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "编辑管理员",
				Type:  "BUTTON",
			}},
			{ParentID: 20, SystemMenu: false, Sort: 3, Name: "admin:delete", Path: "admin:delete", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "删除管理员",
				Type:  "BUTTON",
			}},
			{ParentID: 20, SystemMenu: false, Sort: 4, Name: "admin:read", Path: "admin:read", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "查询管理员",
				Type:  "BUTTON",
			}},
			{ParentID: 20, SystemMenu: false, Sort: 5, Name: "admin:password", Path: "admin:password", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "重置密码",
				Type:  "BUTTON",
			}},
			{ParentID: 2, SystemMenu: true, Sort: 5, Name: "settingDict", Path: "/setting/dict", Component: "setting/dict/index", AdministratorID: 1, Meta: database.MenuMeta{
				Icon:        "Memo",
				Title:       "字典管理",
				Type:        "MENU",
				IsKeepAlive: true,
			}},
			{ParentID: 26, SystemMenu: false, Sort: 1, Name: "dict:type:create", Path: "dict:type:create", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "新增字典分类",
				Type:  "BUTTON",
			}},
			{ParentID: 26, SystemMenu: false, Sort: 2, Name: "dict:type:update", Path: "dict:type:update", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "编辑字典分类",
				Type:  "BUTTON",
			}},
			{ParentID: 26, SystemMenu: false, Sort: 3, Name: "dict:type:delete", Path: "dict:type:delete", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "删除字典分类",
				Type:  "BUTTON",
			}},
			{ParentID: 26, SystemMenu: false, Sort: 4, Name: "dict:data:create", Path: "dict:data:create", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "新增字典数据",
				Type:  "BUTTON",
			}},
			{ParentID: 26, SystemMenu: false, Sort: 5, Name: "dict:data:update", Path: "dict:data:update", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "编辑字典数据",
				Type:  "BUTTON",
			}},
			{ParentID: 26, SystemMenu: false, Sort: 6, Name: "dict:data:delete", Path: "dict:data:delete", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "删除字典数据",
				Type:  "BUTTON",
			}},
			{ParentID: 2, SystemMenu: true, Sort: 6, Name: "settingLog", Path: "/setting/log", Component: "setting/log/index", AdministratorID: 1, Meta: database.MenuMeta{
				Icon:        "Clock",
				Title:       "操作日志",
				Type:        "MENU",
				IsKeepAlive: true,
			}},
			{ParentID: 33, SystemMenu: false, Sort: 1, Name: "log:delete", Path: "log:delete", Component: "", AdministratorID: 1, Meta: database.MenuMeta{
				Title: "删除操作日志",
				Type:  "BUTTON",
			}},
		}
		if err := global.FC_DB.Create(&entities).Error; err != nil {
			global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】创建默认菜单失败：%v", err)) //链接失败后可以关闭程序了，所以使用logging.Fatal方法
		}
	}

	//自动创建字典分类及字典数据
	err = global.FC_DB.Model(&database.SysDictType{}).Count(&count).Error
	if err != nil {
		global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】查询默认菜单数据库失败：%v", err)) //链接失败后可以关闭程序了，所以使用logging.Fatal方法
	}
	if count == 0 {
		dictType := database.SysDictType{
			Name:            "性别",
			Type:            "sex",
			Status:          1,
			Notes:           "管理员/用户性别",
			AdministratorID: 1,
			SysDictData: []*database.SysDictData{
				{Label: "未知", Value: "0", Status: 1, Sort: 1, Notes: "", DictType: "sex"},
				{Label: "男", Value: "1", Status: 1, Sort: 2, Notes: "", DictType: "sex"},
				{Label: "女", Value: "2", Status: 1, Sort: 3, Notes: "", DictType: "sex"},
			},
		}
		if err := global.FC_DB.Create(&dictType).Error; err != nil {
			global.FC_LOGGER.Fatal(fmt.Sprintf("【数据库】创建默认字典类型失败：%v", err)) //链接失败后可以关闭程序了，所以使用logging.Fatal方法
		}
	}
}
