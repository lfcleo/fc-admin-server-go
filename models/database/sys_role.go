package database

import (
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/request"
	casbinUtil "fc-admin-server-go/pkg/casbin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Role 角色模型
type Role struct {
	BaseModel
	Name            string           `gorm:"type:varchar(20);unique;comment:角色名称;" json:"name,omitempty"`       //角色名称
	Code            string           `gorm:"type:varchar(20);unique;comment:角色编码;" json:"code,omitempty"`       //角色编码
	Notes           string           `gorm:"type:varchar(255);comment:备注;" json:"notes"`                        //备注
	AdministratorID uint             `gorm:"index" json:"administratorID"`                                      //创建角色的管理员ID
	Administrator   *Administrator   `json:"administrator,omitempty"`                                           //创建角色的管理员模型
	Administrators  []*Administrator `gorm:"many2many:sys_administrator_role;" json:"administrators,omitempty"` //角色包含的管理员
	Menus           []*Menu          `gorm:"many2many:sys_role_menu;" json:"menus,omitempty"`                   //角色包含的菜单权限
}

// TableName Administrator 表名重命名
func (Role) TableName() string {
	return "sys_role"
}

// CreateRole 创建角色
func (role *Role) CreateRole() error {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	////创建角色
	//err = tx.Create(role).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	////创建角色默认菜单，拼接中间表属性
	//roleMenus := []RoleMenu{
	//	{RoleID: role.ID, MenuID: 1},
	//}
	//// 创建中间表数据
	//err = tx.Create(&roleMenus).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return
	return global.FC_DB.Transaction(func(tx *gorm.DB) (err error) {
		//创建角色
		if err = tx.Create(role).Error; err != nil {
			return
		}
		//创建角色默认菜单，拼接中间表属性
		roleMenus := []RoleMenu{
			{RoleID: role.ID, MenuID: 1},
		}
		// 创建中间表数据
		return tx.Create(&roleMenus).Error
	})
}

// UpdateRole 更新角色信息
func (role *Role) UpdateRole() (err error) {
	return global.FC_DB.Model(role).Updates(role).Error
}

// UnscopedDeleteRole 删除角色信息(永久删除),同时删除管理员，菜单的中间表
func (role *Role) UnscopedDeleteRole() (err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 删除casbin中权限
	//_, err = casbinUtil.CasbinServiceApp.ClearCasbin(0, strconv.Itoa(int(role.ID)))
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 删除casbin中的组
	//_, err = casbinUtil.CasbinServiceApp.ClearGroupingPolicies(0, strconv.Itoa(int(role.ID)))
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	////删除角色信息
	//err = tx.Select("Administrators", "Menus").Unscoped().Delete(role).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return

	return global.FC_DB.Transaction(func(tx *gorm.DB) (err error) {
		// 删除casbin中权限
		if _, err = casbinUtil.CasbinServiceApp.ClearCasbin(0, strconv.Itoa(int(role.ID))); err != nil {
			return
		}
		// 删除casbin中的组
		if _, err = casbinUtil.CasbinServiceApp.ClearGroupingPolicies(0, strconv.Itoa(int(role.ID))); err != nil {
			return
		}
		//删除角色信息
		return tx.Select("Administrators", "Menus").Unscoped().Delete(role).Error
	})
}

// QueryRoleList 角色列表
func QueryRoleList(queryStr string, page, pageSize int) (list []*Role, count int64, err error) {
	page = request.GetPage(page, pageSize)
	result := global.FC_DB.Model(&Role{}).Where(queryStr).Count(&count).Offset(page).Limit(pageSize).Find(&list)
	err = result.Error
	return
}

// QueryAllRole 获取所有角色列表
func QueryAllRole() (list []*Role, err error) {
	// 获取全部记录
	err = global.FC_DB.Find(&list).Error
	return
}

// UpdateRoleApis 更新角色的API接口列表
func UpdateRoleApis(rID uint, cis []casbinUtil.CasbinInfo) (tTime time.Time, err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	////更新casbin权限
	//err = casbinUtil.CasbinServiceApp.UpdateCasbin(rID, cis)
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 更新角色的updatedAt
	//tTime = time.Now()
	//err = tx.Model(Role{}).Where("id = ?", rID).Update("updated_at", tTime).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return

	err = global.FC_DB.Transaction(func(tx *gorm.DB) (err error) {
		//更新casbin权限
		if err = casbinUtil.CasbinServiceApp.UpdateCasbin(rID, cis); err != nil {
			return
		}
		// 更新角色的updatedAt
		tTime = time.Now()
		return tx.Model(Role{}).Where("id = ?", rID).Update("updated_at", tTime).Error
	})
	return
}
