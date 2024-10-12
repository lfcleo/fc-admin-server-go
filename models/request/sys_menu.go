package request

import (
	"strings"
)

type Menu struct {
	Title string `json:"title,omitempty"` //菜单名称
}

// BuildQueryConditions 构建查询条件字符串
func (rM *Menu) BuildQueryConditions() string {
	var conditions []string
	if rM.Title != "" {
		conditions = append(conditions, "meta->>'$.title' LIKE '%"+rM.Title+"%'")
	}
	//如果有查询值，拼接查询字符串
	if len(conditions) > 0 {
		return strings.Join(conditions, " AND ")
	}
	return ""
}
