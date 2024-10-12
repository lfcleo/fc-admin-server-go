package database

import (
	"fc-admin-server-go/global"
	"gorm.io/gorm"
)

type SysDictType struct {
	BaseModel
	Name            string         `gorm:"type:varchar(30);unique;comment:字典名称;" json:"name,omitempty"`       //字典名称
	Type            string         `gorm:"type:varchar(30);unique;comment:字典类型;" json:"type,omitempty"`       //字典类型
	Status          int            `gorm:"type:tinyint(1);default:1;comment:字典状态,默认=1正常,=2停用;" json:"status"` //字典状态,默认=1正常,=2停用
	Notes           string         `gorm:"type:varchar(255);comment:备注;" json:"notes"`                        //备注
	AdministratorID uint           `gorm:"index" json:"administratorID"`                                      //创建角色的管理员ID
	Administrator   *Administrator `json:"administrator,omitempty"`                                           //创建角色的管理员模型
	SysDictData     []*SysDictData `gorm:"foreignKey:DictType;references:Type;" json:"SysDictData,omitempty"` //与SysDictData一对多关系，重写外键是DictType，重写引用是Type
}

// TableName SysDictType 表名重命名
func (SysDictType) TableName() string {
	return "sys_dict_type"
}

// CreateDictType 创建字典分类
func (sDT *SysDictType) CreateDictType() error {
	//return global.FC_DB.Create(sDT).Error
	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(sDT).Error
	})
}

// UpdateDictType 更新字典分类
func UpdateDictType(sDT *SysDictType) (err error) {
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 先查询原数据，主要获取type的值
	//var dT SysDictType
	//err = tx.Where("id = ?", sDT.ID).First(&dT).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 更新dictData表中type的值
	//if dT.Type != sDT.Type {
	//	err = tx.Model(SysDictData{}).Where("dict_type = ?", dT.Type).Updates(SysDictData{DictType: sDT.Type}).Error
	//	if err != nil {
	//		tx.Rollback()
	//		return
	//	}
	//}
	//// 更新dictType表中的值
	//err = tx.Model(sDT).Updates(sDT).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return
	return global.FC_DB.Transaction(func(tx *gorm.DB) (err error) {
		// 先查询原数据，主要获取type的值
		var dT SysDictType
		if err = tx.Where("id = ?", sDT.ID).First(&dT).Error; err != nil {
			return err
		}
		// 更新dictData表中type的值
		if dT.Type != sDT.Type {
			if err = tx.Model(SysDictData{}).Where("dict_type = ?", dT.Type).Updates(SysDictData{DictType: sDT.Type}).Error; err != nil {
				return err
			}
		}
		// 更新dictType表中的值
		return tx.Model(sDT).Updates(sDT).Error
	})
}

// QueryDictTypeList 字典分类列表
func QueryDictTypeList() (list []*SysDictType, err error) {
	result := global.FC_DB.Model(&SysDictType{}).Find(&list)
	err = result.Error
	return
}

// UnscopedDeleteSysDictType 永久删除字典分类
func UnscopedDeleteSysDictType(sDT *SysDictType) (err error) {
	////return global.FC_DB.Model(SysDictType{}).Association("SysDictData").Delete(sDT)
	//tx := global.FC_DB.Begin()
	//defer func() {
	//	if r := recover(); r != nil { //如果出现问题的话，回滚事务
	//		tx.Rollback()
	//	}
	//}()
	//// 先查询原数据，主要获取type的值
	//var dT SysDictType
	//err = tx.Where("id = ?", sDT.ID).First(&dT).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 删除dictData表中type的值
	//err = tx.Model(SysDictData{}).Where("dict_type = ?", dT.Type).Unscoped().Delete(SysDictData{}).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//// 删除dictType表中的值
	//err = tx.Unscoped().Delete(sDT).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//err = tx.Commit().Error
	//return
	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		// 先查询原数据，主要获取type的值
		var dT SysDictType
		if err = tx.Where("id = ?", sDT.ID).First(&dT).Error; err != nil {
			return err
		}
		// 删除dictData表中type的值
		if err = tx.Model(SysDictData{}).Where("dict_type = ?", dT.Type).Unscoped().Delete(SysDictData{}).Error; err != nil {
			return err
		}
		// 删除dictType表中的值
		return tx.Unscoped().Delete(sDT).Error
	})
}

// QueryDictType 根据key值获取字典数据
func QueryDictType(key string) (dT *SysDictType, err error) {
	result := global.FC_DB.Model(&SysDictType{}).Where("status = 1 AND type = ?", key).Preload("SysDictData", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = 1").Order("sort")
	}).First(&dT)
	err = result.Error
	return
}
