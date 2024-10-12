package database

import (
	"errors"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/request"
	casbinUtil "fc-admin-server-go/pkg/casbin"
	"fc-admin-server-go/pkg/config"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// Administrator 管理员模型
type Administrator struct {
	BaseModel
	Mobile    string   `gorm:"type:char(11);index;unique;not null;comment:管理员手机号;" json:"mobile,omitempty"`  //手机号,加索引，唯一，不为空
	Password  string   `gorm:"type:char(40);comment:管理员密码;" json:"-"`                                        //密码，使用sha-1加密，生成40个字符
	Email     string   `gorm:"type:varchar(25);index;unique;not null;comment:管理员邮箱;" json:"email,omitempty"` //邮箱,加索引，唯一，不为空
	Name      string   `gorm:"type:varchar(20);comment:管理员姓名;" json:"name,omitempty"`                        //用户姓名
	Avatar    string   `gorm:"type:varchar(255);comment:管理员头像;" json:"avatar,omitempty"`                     //用户头像
	Sex       int      `gorm:"type:tinyint(1);default:1;comment:管理员性别;" json:"sex"`                          //性别，1=未知，2=男，3=女
	Status    int      `gorm:"type:tinyint(1);default:1;comment:管理员状态;" json:"status"`                       //用户状态,默认=1正常,=2停用
	RolesData []*Role  `gorm:"many2many:sys_administrator_role;" json:"-"`                                   //角色,多对多关系
	Roles     []string `gorm:"-" json:"roles,omitempty"`                                                     //RolesData中的code字段数组，不保存在数据库中
}

// TableName Administrator 表名重命名
func (Administrator) TableName() string {
	return "sys_administrator"
}

// BeforeCreate 创建之前的钩子，生成用户默认头像
func (administrator *Administrator) BeforeCreate(tx *gorm.DB) (err error) {
	administrator.Avatar = config.Data.Server.ImageSavePath + "avatar.jpeg" //生成用户默认头像
	return
}

// AfterCreate 创建之后的钩子，拼接用户头像地址
func (administrator *Administrator) AfterCreate(tx *gorm.DB) (err error) {
	administrator.Avatar = config.Data.Server.DomainName + administrator.Avatar //拼接用户头像地址
	return
}

// BeforeUpdate 更新之前的钩子，处理用户头像地址
func (administrator *Administrator) BeforeUpdate(tx *gorm.DB) (err error) {
	//判断头像是否是本服务域名头像，是的话，去除前缀
	if strings.HasPrefix(administrator.Avatar, config.Data.Server.DomainName) {
		administrator.Avatar = administrator.Avatar[len(config.Data.Server.DomainName):]
	}
	return
}

// AfterFind 查询之后的钩子，拼接头像,转化权限数组
func (administrator *Administrator) AfterFind(tx *gorm.DB) (err error) {
	if !strings.HasPrefix(administrator.Avatar, "http") {
		administrator.Avatar = config.Data.Server.DomainName + administrator.Avatar //拼接用户头像地址
	}
	for _, role := range administrator.RolesData {
		administrator.Roles = append(administrator.Roles, role.Code)
	}
	return
}

// CreateAdministrator 创建管理员信息
func (administrator *Administrator) CreateAdministrator(roleIDs []uint) (err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	////创建管理员
	//err = tx.Create(administrator).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// casbin 用户添加到组
	//err = casbinUtil.CasbinServiceApp.UpdateRolePolicies(administrator.ID, roleIDs)
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return

	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		//创建管理员
		if err = tx.Create(administrator).Error; err != nil {
			return err
		}
		// casbin 用户添加到组
		return casbinUtil.CasbinServiceApp.UpdateRolePolicies(administrator.ID, roleIDs)
	})
}

// FindAdministratorByKey 根据所需字段查找用户，如id,mobile，预加载字段
func FindAdministratorByKey(key, value interface{}, preloadFields ...string) (admin *Administrator, err error) {
	qu := fmt.Sprintf("%s = ?", key)
	var tempDB = global.FC_DB
	for _, field := range preloadFields {
		tempDB = tempDB.Preload(field)
	}
	result := tempDB.Where(qu, value).First(&admin)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, err
	}
	return
}

// UpdateAdministrator 更新管理员信息
func (administrator *Administrator) UpdateAdministrator() (err error) {
	return global.FC_DB.Model(administrator).Updates(administrator).Error
}

//// UpdateAdministratorByKey 根据字段更新管理员单个信息
//func (administrator *Administrator) UpdateAdministratorByKey(key string, value interface{}) (err error) {
//	return global.FC_DB.Model(administrator).Update(key, value).Error
//}

// UnscopedDeleteAdministrator 删除管理员信息(永久删除)
func (administrator *Administrator) UnscopedDeleteAdministrator() (err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 删除casbin中的组
	//_, err = casbinUtil.CasbinServiceApp.ClearGroupingPolicies(0, strconv.Itoa(int(administrator.ID)))
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 删除管理员信息
	//err = tx.Unscoped().Delete(administrator).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return
	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		// 删除casbin中的组
		if _, err = casbinUtil.CasbinServiceApp.ClearGroupingPolicies(0, strconv.Itoa(int(administrator.ID))); err != nil {
			return err
		}
		// 删除管理员信息
		return tx.Unscoped().Delete(administrator).Error
	})
}

// QueryAdministratorList 管理员列表
func QueryAdministratorList(queryStr string, page, pageSize int, preloadFields ...string) (list []*Administrator, count int64, err error) {
	page = request.GetPage(page, pageSize)

	var tempDB = global.FC_DB
	for _, field := range preloadFields {
		tempDB = tempDB.Preload(field)
	}
	result := tempDB.Model(&Administrator{}).Where(queryStr).Count(&count).Offset(page).Limit(pageSize).Find(&list)
	err = result.Error
	return
}
