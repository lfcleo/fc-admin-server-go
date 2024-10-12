package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fc-admin-server-go/global"
	"gorm.io/gorm"
	"time"
)

// Menu 菜单模型
type Menu struct {
	BaseModel
	ParentID        uint           `gorm:"comment:父菜单路由ID;" json:"parentID"`                              //上级菜单路由的ID，=0 代表一级菜单路由
	SystemMenu      bool           `gorm:"type:bool;default:false;comment:是否是系统路由;" json:"systemMenu"`    //是否是系统菜单/路由，系统路由禁止删除
	Sort            uint8          `gorm:"type:tinyint unsigned;default:0;comment:排序标记" json:"sort"`      // 排序标记
	Path            string         `gorm:"type:varchar(255);unique;not null;comment:菜单路由路径;" json:"path"` //菜单路由路径
	Name            string         `gorm:"type:varchar(20);unique;not null;comment:菜单路由名称;" json:"name"`  //菜单路由名称
	Redirect        string         `gorm:"type:varchar(255);comment:菜单路由重定向地址;" json:"redirect"`          //菜单路由重定向地址
	Component       string         `gorm:"type:varchar(255);comment:菜单路由文件所在地址;" json:"component"`        //菜单路由文件所在地址
	Meta            MenuMeta       `gorm:"type:json;comment:菜单路由元数据;" json:"meta"`                        //菜单路由元数据
	Children        []*Menu        `gorm:"-" json:"children,omitempty"`                                   //子菜单数据
	AdministratorID uint           `gorm:"index" json:"administratorID"`                                  //创建菜单的管理员ID
	Administrator   *Administrator `json:"administrator,omitempty"`                                       //创建菜单的管理员模型
	Roles           []*Role        `gorm:"many2many:sys_role_menu;" json:"roles,omitempty"`               //接口包含的角色权限
}

// MenuMeta 菜单路由元数据
type MenuMeta struct {
	Icon        string `json:"icon,omitempty"`        //菜单和面包屑对应的图标
	Title       string `json:"title"`                 //菜单标题
	Type        string `json:"type"`                  //菜单类型，MENU=菜单，LINK=外链，BUTTON=按钮
	IsHide      bool   `json:"isHide,omitempty"`      //是否在菜单中隐藏
	IsFull      bool   `json:"isFull,omitempty"`      //单是否全屏
	IsAffix     bool   `json:"isAffix,omitempty"`     //菜单是否固定在标签页中
	IsKeepAlive bool   `json:"isKeepAlive,omitempty"` //是否缓存路由
	Tag         string `json:"tag,omitempty"`         //标签，会在菜单栏中显示红色角标
}

// TableName Menu 表名重命名
func (Menu) TableName() string {
	return "sys_menu"
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (mm *MenuMeta) Scan(value interface{}) error {
	arr, ok := value.([]byte)
	if !ok {
		return errors.New("不匹配的数据类型")
	}
	return json.Unmarshal(arr, mm)
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (mm MenuMeta) Value() (driver.Value, error) {
	return json.Marshal(mm)
}

// CreateMenu 创建菜单
func (menu *Menu) CreateMenu() error {
	//return global.FC_DB.Create(menu).Error
	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(menu).Error
	})
}

//// AfterUpdate 更新之后的钩子
//func (menu *Menu) AfterUpdate(tx *gorm.DB) (err error) {
//	err = tx.Model(Role{}).Where("id = 1").Update("updated_at", menu.UpdatedAt).Error
//	return
//}

// QueryMenuList 菜单接口列表
func QueryMenuList(queryStr string) (list []*Menu, err error) {
	err = global.FC_DB.Model(&Menu{}).Where(queryStr).Order("sort").Find(&list).Error
	return
}

// QueryMenuCountByParentID 查询菜单的子菜单个数
func QueryMenuCountByParentID(mID uint) (count int64, err error) {
	err = global.FC_DB.Model(&Menu{}).Where("parent_id = ?", mID).Count(&count).Error
	return
}

// BuildMenuTree 从给定的菜单列表中构建树形结构
func BuildMenuTree(menus []*Menu) []*Menu {
	// 创建一个 map 来存储每个菜单项
	menuMap := make(map[uint]*Menu)
	// 创建一个 slice 来存储顶级节点
	var roots []*Menu
	// 首先将所有节点添加到 map 中
	for _, m := range menus {
		menuMap[m.ID] = m
	}
	// 遍历所有菜单项，并将它们添加到它们各自的父节点下
	for _, m := range menus {
		if parent, ok := menuMap[m.ParentID]; ok {
			// 如果找到了父节点，就将当前节点添加到父节点的 Children 列表中
			parent.Children = append(parent.Children, m)
		} else {
			// 如果没有找到父节点，则当前节点是顶级节点
			roots = append(roots, m)
		}
	}
	return roots
}

// UpdateMenu 更新菜单信息
func UpdateMenu(menu *Menu) (roleIDs []uint, err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 更新菜单信息,使用map更新，确保指定字段被更新
	//mapValue := map[string]interface{}{
	//	"path":      menu.Path,
	//	"name":      menu.Name,
	//	"redirect":  menu.Redirect,
	//	"component": menu.Component,
	//	"meta":      menu.Meta,
	//	"sort":      menu.Sort,
	//}
	//err = tx.Model(menu).Updates(mapValue).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 获取菜单的所有角色权限
	//roles, err := QueryMenuRoles(menu)
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//for _, role := range roles {
	//	roleIDs = append(roleIDs, role.ID)
	//}
	//// 设置角色的updateAt
	//err = tx.Table("sys_role").Where("id IN ?", roleIDs).Update("updated_at", menu.UpdatedAt).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return
	err = global.FC_DB.Transaction(func(tx *gorm.DB) error {
		// 更新菜单信息,使用map更新，确保指定字段被更新
		mapValue := map[string]interface{}{
			"path":      menu.Path,
			"name":      menu.Name,
			"redirect":  menu.Redirect,
			"component": menu.Component,
			"meta":      menu.Meta,
			"sort":      menu.Sort,
		}
		if err = tx.Model(menu).Updates(mapValue).Error; err != nil {
			return err
		}
		var roles []*Role
		if err = tx.Model(menu).Association("Roles").Find(&roles); err != nil {
			return err
		}
		for _, role := range roles {
			roleIDs = append(roleIDs, role.ID)
		}
		roleIDs = append(roleIDs, 1) //把超级管理员角色添加进去
		// 设置角色的updateAt
		err = tx.Table("sys_role").Where("id IN ?", roleIDs).Update("updated_at", menu.UpdatedAt).Error
		return err
	})
	return
}

// UnscopedDeleteMenu 删除接口信息(永久删除),同时删除role的中间表
func UnscopedDeleteMenu(menu *Menu) (roleIDs []uint, tTime time.Time, err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 获取菜单的所有角色权限
	//roles, err := QueryMenuRoles(menu)
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//if len(roles) > 0 {
	//	for _, role := range roles {
	//		roleIDs = append(roleIDs, role.ID)
	//	}
	//	tTime = time.Now()
	//	// 设置角色的updateAt
	//	err = tx.Table("sys_role").Where("id IN ?", roleIDs).Update("updated_at", tTime).Error
	//	if err != nil {
	//		tx.Rollback()
	//		return
	//	}
	//}
	//// 删除菜单，同时删除role的中间表
	//err = tx.Select("Roles").Unscoped().Where("system_menu = 0").Delete(menu).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return
	err = global.FC_DB.Transaction(func(tx *gorm.DB) error {
		// 获取菜单的所有角色权限
		roles, err := QueryMenuRoles(menu)
		if err != nil {
			return err
		}
		if len(roles) > 0 {
			for _, role := range roles {
				roleIDs = append(roleIDs, role.ID)
			}
			tTime = time.Now()
			// 设置角色的updateAt
			if err = tx.Table("sys_role").Where("id IN ?", roleIDs).Update("updated_at", tTime).Error; err != nil {
				return err
			}
		}
		roleIDs = append(roleIDs, 1) //把超级管理员角色添加进去
		// 删除菜单，同时删除role的中间表
		return tx.Select("Roles").Unscoped().Where("system_menu = 0").Delete(menu).Error
	})
	return
}

// QueryAdminMenu 获取管理员动态菜单和权限
func QueryAdminMenu(aID uint) (list []*Menu, permissions []string, err error) {
	var menuList []*Menu
	err = global.FC_DB.Raw(`SELECT DISTINCT
				m.*
			FROM
				sys_menu m
				LEFT JOIN sys_role_menu rm ON m.id = rm.menu_id
				LEFT JOIN sys_administrator_role ur ON rm.role_id = ur.role_id
				LEFT JOIN sys_role r ON r.id = ur.role_id 
			WHERE
				ur.administrator_id = ?
			ORDER BY
			    m.sort ASC;`, aID).
		Find(&menuList).Error
	if err != nil {
		return nil, nil, err
	}

	// 递归过滤菜单栏,获取
	list, permissions = MenuTree(menuList, 0, true)

	return
}

// MenuTree 递归菜单列表转树形结构，getPermission=true 将会在treeList排除meta.type == BUTTON的菜单，并添加到permissions字符串数组权限中
func MenuTree(menuList []*Menu, pid uint, getPermission bool) (treeList []*Menu, permissions []string) {
	for _, v := range menuList {
		if v.ParentID == pid {
			v1, excludedNames1 := MenuTree(menuList, v.ID, getPermission)
			v.Children = v1
			if len(excludedNames1) > 0 {
				// 将子菜单中排除的按钮名添加到总的排除名单中
				permissions = append(permissions, excludedNames1...)
			}
			if getPermission && v.Meta.Type == "BUTTON" {
				permissions = append(permissions, v.Name)
				continue
			}
			treeList = append(treeList, v)
		}
	}
	return
}
