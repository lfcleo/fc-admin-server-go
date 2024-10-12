package database

import (
	"fc-admin-server-go/global"
	casbinUtil "fc-admin-server-go/pkg/casbin"
	"gorm.io/gorm"
)

type AdministratorRole struct {
	AdministratorID uint
	RoleID          uint
}

// TableName AdministratorRole 表名重命名
func (AdministratorRole) TableName() string {
	return "sys_administrator_role"
}

// UpdateAdministratorRoles 更新管理员的角色
func UpdateAdministratorRoles(aID uint, roleIDs []uint) (err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 删除关联表中的 administrator_id = ？的所有数据
	//err = tx.Delete(&AdministratorRole{}, "administrator_id = ?", aID).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 拼接中间表属性
	//var adminRoles []AdministratorRole
	//for _, v := range roleIDs {
	//	adminRoles = append(adminRoles, AdministratorRole{
	//		AdministratorID: aID,
	//		RoleID:          v,
	//	})
	//}
	//// 创建中间表数据
	//err = tx.Create(&adminRoles).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// casbin 用户添加到组
	//err = casbinUtil.CasbinServiceApp.UpdateRolePolicies(aID, roleIDs)
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return
	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联表中的 administrator_id = ？的所有数据
		if err = tx.Delete(&AdministratorRole{}, "administrator_id = ?", aID).Error; err != nil {
			return err
		}
		// 拼接中间表属性
		var adminRoles []AdministratorRole
		for _, v := range roleIDs {
			adminRoles = append(adminRoles, AdministratorRole{
				AdministratorID: aID,
				RoleID:          v,
			})
		}
		// 创建中间表数据
		if err = tx.Create(&adminRoles).Error; err != nil {
			return err
		}
		// casbin 用户添加到组
		return casbinUtil.CasbinServiceApp.UpdateRolePolicies(aID, roleIDs)
	})
}
