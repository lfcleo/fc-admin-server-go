package request

import "strings"

type DictData struct {
	BaseList
	DictType string `json:"dictType" binding:"required"`
}

// BuildQueryConditions 构建查询条件字符串
func (dD *DictData) BuildQueryConditions() string {
	var conditions []string
	if dD.DictType != "" {
		conditions = append(conditions, "dict_type = '"+dD.DictType+"'")
	}
	//如果有查询值，拼接查询字符串
	if len(conditions) > 0 {
		return strings.Join(conditions, " AND ")
	}
	return ""
}
