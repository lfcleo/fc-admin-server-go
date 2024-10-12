package request

import (
	"strings"
)

type Api struct {
	BaseList
	Method      string `json:"method,omitempty"`      //接口请求方式
	Description string `json:"description,omitempty"` //接口描述
}

// BuildQueryConditions 构建查询条件字符串
func (rA *Api) BuildQueryConditions() string {
	var conditions []string
	if rA.Method != "" {
		conditions = append(conditions, "method = '"+rA.Method+"'")
	}
	if rA.Description != "" {
		conditions = append(conditions, "description LIKE '%"+rA.Description+"%'")
	}
	//如果有查询值，拼接查询字符串
	if len(conditions) > 0 {
		return strings.Join(conditions, " AND ")
	}
	return ""
}

type CreateApi struct {
	BaseID
	Path        string `json:"path"`        //请求路径
	Method      string `json:"method"`      //请求方式
	Description string `json:"description"` //接口描述
}
