package database

import (
	"fc-admin-server-go/global"
	"gorm.io/gorm"
	"time"
)

// RoleMenu 角色表和菜单表多对多关系的中间表
type RoleMenu struct {
	RoleID uint
	MenuID uint
}

// TableName RoleApi 表名重命名
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

// QueryRoleMenus 查询角色的菜单列表
func QueryRoleMenus(role *Role) (menus []*Menu, err error) {
	err = global.FC_DB.Model(role).Association("Menus").Find(&menus)
	return
}

// QueryMenuRoles 查询菜单的角色列表
func QueryMenuRoles(menu *Menu) (roles []*Role, err error) {
	err = global.FC_DB.Model(menu).Association("Roles").Find(&roles)
	return
}

// UpdateRoleMenus 更新角色的菜单列表
func UpdateRoleMenus(rID uint, menuIDs []uint) (tTime time.Time, err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 删除关联表中的 role_id = ？的所有数据
	//err = tx.Delete(&RoleMenu{}, "role_id = ?", rID).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 拼接中间表属性
	//var roleMenus []RoleMenu
	//for _, v := range menuIDs {
	//	roleMenus = append(roleMenus, RoleMenu{
	//		RoleID: rID,
	//		MenuID: v,
	//	})
	//}
	//// 创建中间表数据
	//err = tx.Create(&roleMenus).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	////更新角色的updatedAt
	//tTime = time.Now()
	//err = tx.Model(Role{}).Where("id = ?", rID).Update("updated_at", tTime).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return

	err = global.FC_DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联表中的 role_id = ？的所有数据
		if err = tx.Delete(&RoleMenu{}, "role_id = ?", rID).Error; err != nil {
			return err
		}
		// 拼接中间表属性
		var roleMenus []RoleMenu
		for _, v := range menuIDs {
			roleMenus = append(roleMenus, RoleMenu{
				RoleID: rID,
				MenuID: v,
			})
		}
		// 创建中间表数据
		if err = tx.Create(&roleMenus).Error; err != nil {
			return err
		}
		//更新角色的updatedAt
		tTime = time.Now()
		return tx.Model(Role{}).Where("id = ?", rID).Update("updated_at", tTime).Error
	})
	return
}
