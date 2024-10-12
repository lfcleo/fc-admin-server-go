package database

import (
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/request"
	"gorm.io/gorm"
)

type SysDictData struct {
	BaseModel
	Label    string `gorm:"type:varchar(30);comment:字典标签;" json:"label,omitempty"`             //字典标签
	Value    string `gorm:"type:varchar(30);comment:字典值;" json:"value,omitempty"`              //字典值
	Status   int    `gorm:"type:tinyint(1);default:1;comment:字典状态,默认=1正常,=2停用;" json:"status"` //字典状态,默认=1正常,=2停用
	Sort     uint8  `gorm:"type:tinyint unsigned;default:0;comment:排序标记" json:"sort" `         // 排序标记
	Notes    string `gorm:"type:varchar(255);comment:备注;" json:"notes"`                        //备注
	DictType string `gorm:"index" json:"dictType"`                                             //关联标记
}

// TableName SysDictData 表名重命名
func (SysDictData) TableName() string {
	return "sys_dict_data"
}

// CreateDictData 创建字典详情
func (sDD *SysDictData) CreateDictData() error {
	//return global.FC_DB.Create(sDD).Error
	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(sDD).Error
	})
}

// QueryDictDataList 字典数据列表
func QueryDictDataList(queryStr string, page, pageSize int) (list []*SysDictData, count int64, err error) {
	page = request.GetPage(page, pageSize)
	result := global.FC_DB.Model(&SysDictData{}).Where(queryStr).Order("sort").Count(&count).Offset(page).Limit(pageSize).Find(&list)
	err = result.Error
	return
}

// UpdateDictData 更新字典数据
func (sDT *SysDictData) UpdateDictData() (err error) {
	return global.FC_DB.Model(sDT).Updates(sDT).Error
}

// UnscopedDeleteSysDictData 根据ID数组批量删除字典数据
func UnscopedDeleteSysDictData(ids []uint) (err error) {
	return global.FC_DB.Unscoped().Delete(&SysDictData{}, "id IN ?", ids).Error
}
