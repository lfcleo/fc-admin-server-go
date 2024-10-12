package database

import (
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/request"
	casbinUtil "fc-admin-server-go/pkg/casbin"
	"gorm.io/gorm"
	"time"
)

type Api struct {
	BaseModel
	Path            string         `gorm:"type:varchar(255);not null;comment:接口路径;" json:"path,omitempty"` //接口路径
	Method          string         `gorm:"type:varchar(10);not null;comment:请求方式;" json:"method"`          //接口请求方式
	Description     string         `gorm:"type:varchar(20);not null;comment:接口描述;" json:"description"`     //接口描述
	AdministratorID uint           `gorm:"index" json:"administratorID"`                                   //创建角色的管理员ID
	Administrator   *Administrator `json:"administrator,omitempty"`                                        //创建角色的管理员模型
}

// TableName Api 表名重命名
func (Api) TableName() string {
	return "sys_api"
}

// CreateApi 创建接口
func (api *Api) CreateApi() (err error) {
	//return global.FC_DB.Create(api).Error
	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(api).Error
	})
}

// QueryApiList Api接口列表
func QueryApiList(queryStr string, page, pageSize int) (list []*Api, count int64, err error) {
	page = request.GetPage(page, pageSize)
	result := global.FC_DB.Model(&Api{}).Where(queryStr).Count(&count).Order("id desc").Offset(page).Limit(pageSize).Find(&list)
	err = result.Error
	return
}

// QueryAllApiList 获取所有Api接口列表
func QueryAllApiList() (list []*Api, err error) {
	result := global.FC_DB.Model(&Api{}).Order("id desc").Find(&list)
	err = result.Error
	return
}

// UpdateApi 更新接口信息
func UpdateApi(api *Api) (roleIDs []uint, err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 更新api接口信息
	//err = tx.Model(api).Updates(api).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//
	//// 获取api接口的所有角色权限
	//roleIDs, err = casbinUtil.CasbinServiceApp.GetRolePolicyByApiInfo(tx, api.Path, api.Method)
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//
	//// 设置角色的updateAt
	//if len(roleIDs) > 0 {
	//	err = tx.Table("sys_role").Where("id IN ?", roleIDs).Update("updated_at", api.UpdatedAt).Error
	//	if err != nil {
	//		tx.Rollback()
	//		return
	//	}
	//}
	//err = tx.Commit().Error
	//return
	err = global.FC_DB.Transaction(func(tx *gorm.DB) error {
		// 更新api接口信息
		if err = tx.Model(api).Updates(api).Error; err != nil {
			return err
		}
		// 获取api接口的所有角色权限
		roleIDs, err = casbinUtil.CasbinServiceApp.GetRolePolicyByApiInfo(tx, api.Path, api.Method)
		if err != nil {
			return err
		}
		// 设置角色的updateAt
		if len(roleIDs) > 0 {
			err = tx.Table("sys_role").Where("id IN ?", roleIDs).Update("updated_at", api.UpdatedAt).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

// UnscopedDeleteApi 删除接口信息(永久删除),
func UnscopedDeleteApi(api *Api) (roleIDs []uint, tTime time.Time, err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//
	//err = tx.First(&api).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 获取api接口的所有角色权限
	//roleIDs, err = casbinUtil.CasbinServiceApp.GetRolePolicyByApiInfo(tx, api.Path, api.Method)
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//
	//if len(roleIDs) > 0 {
	//	tTime = time.Now()
	//	// 设置角色的updateAt
	//	err = tx.Table("sys_role").Where("id IN ?", roleIDs).Update("updated_at", tTime).Error
	//	if err != nil {
	//		tx.Rollback()
	//		return
	//	}
	//}
	//
	//// 删除 casbin 中的数据
	//_, err = casbinUtil.CasbinServiceApp.ClearCasbin(1, api.Path, api.Method)
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//
	//// 删除api接口
	//err = tx.Unscoped().Delete(api).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return

	err = global.FC_DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.First(&api).Error; err != nil {
			return err
		}
		// 获取api接口的所有角色权限
		roleIDs, err = casbinUtil.CasbinServiceApp.GetRolePolicyByApiInfo(tx, api.Path, api.Method)
		if err != nil {
			return err
		}
		if len(roleIDs) > 0 {
			tTime = time.Now()
			// 设置角色的updateAt
			if err = tx.Table("sys_role").Where("id IN ?", roleIDs).Update("updated_at", tTime).Error; err != nil {
				return err
			}
		}
		// 删除 casbin 中的数据
		if _, err = casbinUtil.CasbinServiceApp.ClearCasbin(1, api.Path, api.Method); err != nil {
			return err
		}
		// 删除api接口
		return tx.Unscoped().Delete(api).Error
	})
	return
}
