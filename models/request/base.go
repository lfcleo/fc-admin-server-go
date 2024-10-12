package request

import "fc-admin-server-go/pkg/config"

// Base 请求参数公共结构体
type Base struct {
	Timestamp int64  `json:"timestamp"` //时间戳
	Version   string `json:"version"`   //版本号
}

// BaseList 请求列表基类参数
type BaseList struct {
	Page     int `json:"page,omitempty"`                             //查看列表分页使用
	PageSize int `json:"pageSize,omitempty" binding:"gte=0,lte=100"` //每页显示多少条数据,最多50条
}

// BaseID 请求id基类参数
type BaseID struct {
	ID uint `json:"id"`
}

// BaseIDs 请求ids基类参数
type BaseIDs struct {
	IDs []uint `json:"ids"`
}

// DefaultAPIModel 默认api路由模型
type DefaultAPIModel struct {
	Path   string
	Method string
}

// GetPage 获取页数page
func GetPage(page, pageSize int) int {
	if pageSize == 0 {
		pageSize = config.Data.Server.PageSize
	}
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}
	return result
}
