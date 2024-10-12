package database

import (
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/request"
	"gorm.io/gorm"
)

// SysOperationLog 操作日志
type SysOperationLog struct {
	BaseModel
	IP              string         `gorm:"type:varchar(20);comment:IP地址;" json:"ip"`           //IP主机地址
	Method          string         `gorm:"type:varchar(20);comment:请求方式;" json:"method"`     //请求方式
	Path            string         `gorm:"type:varchar(255);comment:请求路径地址;" json:"path"`  //请求路径
	Agent           string         `gorm:"type:text;comment:代理;" json:"agent"`                 //代理
	Request         string         `gorm:"type:text;comment:请求参数;" json:"request"`           //请求参数
	Response        string         `gorm:"type:text;comment:响应结果;" json:"response"`          //响应结果
	Code            int            `gorm:"type:tinyint unsigned;comment:响应状态码" json:"code"` //响应状态码
	AdministratorID uint           `gorm:"index;comment:管理员ID;" json:"administratorID"`       //管理员ID
	Administrator   *Administrator `json:"administrator"`                                        //管理员模型
}

// CreateOperationLog 创建操作日志
func (log *SysOperationLog) CreateOperationLog() error {
	//return global.FC_DB.Create(log).Error
	return global.FC_DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(log).Error
	})
}

// QueryOperationLogList 操作日志列表
func QueryOperationLogList(queryStr string, page, pageSize int, preloadFields ...string) (list []*SysOperationLog, count int64, err error) {
	page = request.GetPage(page, pageSize)

	var tempDB = global.FC_DB
	for _, field := range preloadFields {
		tempDB = tempDB.Preload(field)
	}
	result := tempDB.Model(&SysOperationLog{}).Where(queryStr).Count(&count).Offset(page).Limit(pageSize).Order("id desc").Find(&list)
	err = result.Error
	return
}

// UnscopedDeleteSysOperationLog 根据ID数组批量删除
func UnscopedDeleteSysOperationLog(ids []uint) (err error) {
	return global.FC_DB.Unscoped().Delete(&SysOperationLog{}, "id IN ?", ids).Error
}
