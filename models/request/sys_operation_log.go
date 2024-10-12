package request

import (
	"fmt"
	"strings"
)

// OperationLog 网络请求操作日志模型
type OperationLog struct {
	BaseList
	Method          string   `json:"method,omitempty"` //请求方式
	Code            int      `json:"code,omitempty"`   //响应状态码
	Path            string   `json:"path,omitempty"`   //请求路径
	AdministratorID uint     `json:"administratorID"`  //管理员ID
	Date            []string `json:"date,omitempty"`   //查询日期
}

// BuildQueryConditions 构建查询条件字符串
func (rA *OperationLog) BuildQueryConditions() string {
	var conditions []string
	if rA.Method != "" {
		conditions = append(conditions, "method = '"+rA.Method+"'")
	}
	if rA.Code > 0 {
		conditions = append(conditions, fmt.Sprintf("code = %d", rA.Code))
	}
	if rA.Path != "" {
		conditions = append(conditions, "path LIKE '%"+rA.Path+"%'")
	}
	if rA.AdministratorID > 0 {
		conditions = append(conditions, fmt.Sprintf("administrator_id = %d", rA.AdministratorID))
	}
	if len(rA.Date) == 2 {
		conditions = append(conditions, fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", rA.Date[0], rA.Date[1]))
	}
	//如果有查询值，拼接查询字符串
	if len(conditions) > 0 {
		return strings.Join(conditions, " AND ")
	}
	return ""
}
