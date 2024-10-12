package request

import (
	"fc-admin-server-go/models/response"
	"fmt"
	"strings"
)

// Administrator 网络请求管理员模型
type Administrator struct {
	BaseList
	Name   string `json:"name,omitempty"`   //管理员名称
	Status int    `json:"status,omitempty"` //管理员状态，=-1全部,=0关闭，=1正常
	Mobile string `json:"mobile,omitempty"` //管理员手机号
	Email  string `json:"email,omitempty"`  //管理员邮箱
}

// BuildQueryConditions 构建查询条件字符串
func (rA *Administrator) BuildQueryConditions() string {
	var conditions []string
	if rA.Name != "" {
		conditions = append(conditions, "name LIKE '%"+rA.Name+"%'")
	}
	if rA.Status >= 0 {
		conditions = append(conditions, fmt.Sprintf("status = %d", rA.Status))
	}
	if rA.Mobile != "" {
		conditions = append(conditions, "mobile LIKE '%"+rA.Mobile+"%'")
	}
	if rA.Email != "" {
		conditions = append(conditions, "email LIKE '%"+rA.Email+"%'")
	}
	//如果有查询值，拼接查询字符串
	if len(conditions) > 0 {
		return strings.Join(conditions, " AND ")
	}
	return ""
}

// CrateAdministrator 创建管理员
type CrateAdministrator struct {
	Base
	response.Auth
	Password string `json:"password" binding:"required"` //加密的密码
	RoleIDs  []uint `json:"roleIDs" binding:"required"`  //角色ID数组
}

// PasswordAdministrator 修改管理员密码
type PasswordAdministrator struct {
	Base
	BaseID
	Password string `json:"password"`
}

// SetAdministratorRole 设置管理员角色
type SetAdministratorRole struct {
	BaseID
	RoleIDs []uint `json:"roleIDs"`
}
