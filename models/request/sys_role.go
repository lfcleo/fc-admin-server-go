package request

import (
	casbinUtil "fc-admin-server-go/pkg/casbin"
	"strings"
)

// DefaultApiIDs 创建角色Role的默认api接口列表
func DefaultApiIDs() []uint {
	return []uint{1, 2, 3, 4}
}

type Role struct {
	BaseList
	Name string `json:"name,omitempty"` //角色名称
	Code string `json:"code,omitempty"` //角色编码
}

// BuildQueryConditions 构建查询条件字符串
func (rR *Role) BuildQueryConditions() string {
	var conditions []string
	if rR.Name != "" {
		conditions = append(conditions, "name LIKE '%"+rR.Name+"%'")
	}
	if rR.Code != "" {
		conditions = append(conditions, "code LIKE '%"+rR.Code+"%'")
	}
	//如果有查询值，拼接查询字符串
	if len(conditions) > 0 {
		return strings.Join(conditions, " AND ")
	}
	return ""
}

type CreateRole struct {
	BaseID
	Name  string `json:"name"`  //角色昵称
	Code  string `json:"code"`  //角色编码
	Notes string `json:"notes"` //角色备注
}

// SetMenuIDs 设置角色的菜单权限
type SetMenuIDs struct {
	BaseID
	IDs []uint `json:"ids"`
}

// SetApis 设置角色的菜单权限
type SetApis struct {
	BaseID
	Apis []casbinUtil.CasbinInfo `json:"apis"`
}
