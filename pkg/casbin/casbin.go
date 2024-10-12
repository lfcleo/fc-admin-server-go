package casbinUtil

import (
	"errors"
	"fc-admin-server-go/global"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type CasbinService struct{}

var CasbinServiceApp = new(CasbinService)

// // RolePolicy (RoleID, Path, Method) 对应于 `CasbinRule` 表中的 (v0, v1, v2)
//
//	type RolePolicy struct {
//		RoleID string `gorm:"column:v0"`
//		Path   string `gorm:"column:v1"`
//		Method string `gorm:"column:v2"`
//	}
//
// // GetRoles 获取所有角色组
//
//	func (c *CasbinService) GetRoles() ([]string, error) {
//		e := c.Casbin()
//		return e.GetAllRoles()
//	}
//
// // GetRolePolicy 获取所有角色组权限
//
//	func (c *CasbinService) GetRolePolicy() (roles []RolePolicy, err error) {
//		err = global.FC_DB.Model(&gormadapter.CasbinRule{}).Where("ptype = 'p'").Find(&roles).Error
//		if err != nil {
//			return nil, err
//		}
//		return
//	}
//
// // CreateRolePolicy 创建角色组权限（单个）, 已有的会忽略
//
//	func (c *CasbinService) CreateRolePolicy(r RolePolicy) error {
//		e := c.Casbin()
//		_, err := e.AddPolicy(r.RoleID, r.Method, r.Method)
//		if err != nil {
//			return err
//		}
//		return e.SavePolicy()
//	}
//
//// CreateRolePolicies 创建角色组权限(多个)
//func (c *CasbinService) CreateRolePolicies(r [][]string) error {
//	e := c.Casbin()
//	_, err := e.AddPolicies(r)
//	if err != nil {
//		return err
//	}
//	return e.SavePolicy()
//}
//
//// UpdateRolePolicy 修改角色组权限
//func (c *CasbinService) UpdateRolePolicy(old, new RolePolicy) error {
//	e := c.Casbin()
//	_, err := e.UpdatePolicy([]string{old.RoleID, old.Method, old.Method},
//		[]string{new.RoleID, new.Method, new.Method})
//	if err != nil {
//		return err
//	}
//	return e.SavePolicy()
//}
//
//// DeleteRolePolicy 删除角色组权限
//func (c *CasbinService) DeleteRolePolicy(r RolePolicy) error {
//	e := c.Casbin()
//	_, err := e.RemovePolicy(r.RoleID, r.Method, r.Method)
//	if err != nil {
//		return err
//	}
//	return e.SavePolicy()
//}
//
//// DeleteRolePolicies 删除角色组权限（多个）
//func (c *CasbinService) DeleteRolePolicies(r [][]string) error {
//	e := c.Casbin()
//	_, err := e.RemovePolicies(r)
//	if err != nil {
//		return err
//	}
//	return e.SavePolicy()
//}
//
//type User struct {
//	UserName  string
//	RoleNames []string
//}
//
//// GetUsers 获取所有用户以及关联的角色
//func (c *CasbinService) GetUsers() (users []User, err error) {
//	e := c.Casbin()
//	p, err := e.GetGroupingPolicy()
//	if err != nil {
//		return
//	}
//	usernameUser := make(map[string]*User, 0)
//	for _, _p := range p {
//		username, usergroup := _p[0], _p[1]
//		if v, ok := usernameUser[username]; ok {
//			usernameUser[username].RoleNames = append(v.RoleNames, usergroup)
//		} else {
//			usernameUser[username] = &User{UserName: username, RoleNames: []string{usergroup}}
//		}
//	}
//	for _, v := range usernameUser {
//		users = append(users, *v)
//	}
//	return
//}
//
//// UpdateUserRole 角色组中用户, 没有组默认创建
//func (c *CasbinService) UpdateUserRole(username, rolename string) error {
//	e := c.Casbin()
//	_, err := e.AddGroupingPolicy(username, rolename)
//	if err != nil {
//		return err
//	}
//	return e.SavePolicy()
//}
//
//// UpdateUserRoles 角色组中添加用户（多个）, 没有组默认创建
//func (c *CasbinService) UpdateUserRoles(r [][]string) error {
//	e := c.Casbin()
//	_, err := e.AddGroupingPolicies(r)
//	if err != nil {
//		return err
//	}
//	return e.SavePolicy()
//}
//
//// DeleteUserRole 角色组中删除用户
//func (c *CasbinService) DeleteUserRole(username, rolename string) error {
//	e := c.Casbin()
//	_, err := e.RemoveGroupingPolicy(username, rolename)
//	if err != nil {
//		return err
//	}
//	return e.SavePolicy()
//}
//

// CanAccess 验证用户权限
func (c *CasbinService) CanAccess(sub, url, method string) (ok bool, err error) {
	e := c.Casbin()
	return e.Enforce(sub, url, method)
}

type CasbinInfo struct {
	Path   string `json:"path"`   // 路径
	Method string `json:"method"` // 方法
}

// UpdateCasbin 更新角色的api权限列表
func (cs *CasbinService) UpdateCasbin(rID uint, casbinInfos []CasbinInfo) error {
	roleID := strconv.Itoa(int(rID))
	_, err := cs.ClearCasbin(0, roleID)
	if err != nil {
		return err
	}
	var rules [][]string
	//做权限去重处理
	deduplicateMap := make(map[string]bool)
	for _, v := range casbinInfos {
		key := roleID + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{roleID, v.Path, v.Method})
		}
	}
	e := cs.Casbin()
	success, err := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return err
}

// ClearCasbin 清除匹配的权限
func (cs *CasbinService) ClearCasbin(v int, p ...string) (bool, error) {
	e := cs.Casbin()
	isOk, err := e.RemoveFilteredPolicy(v, p...)
	if err != nil {
		return false, err
	}
	return isOk, nil
}

// GetPolicyPathByRoleID 获取角色的API接口权限列表
func (cs *CasbinService) GetPolicyPathByRoleID(rID uint) (cis []CasbinInfo, err error) {
	e := cs.Casbin()
	roleID := strconv.Itoa(int(rID))
	list, err := e.GetFilteredPolicy(0, roleID)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		cis = append(cis, CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return cis, err
}

// GetRolePolicyByApiInfo 获取API接口的权限列表
func (cs *CasbinService) GetRolePolicyByApiInfo(tx *gorm.DB, v1, v2 string) (aIDs []uint, err error) {
	var gRules []gormadapter.CasbinRule
	err = tx.Model(&gormadapter.CasbinRule{}).Where("ptype = 'p' AND v1 = ? AND v2 = ?", v1, v2).Find(&gRules).Error
	if err != nil {
		return nil, err
	}
	for _, rule := range gRules {
		id, err := strconv.ParseUint(rule.V0, 10, 64)
		if err != nil {
			return nil, err
		}
		aIDs = append(aIDs, uint(id))
	}
	return
}

/**--------**/

// UpdateRolePolicies 更新角色组
func (cs *CasbinService) UpdateRolePolicies(aID uint, rIDs []uint) error {
	roleID := strconv.Itoa(int(aID))
	_, err := cs.ClearGroupingPolicies(0, roleID)
	if err != nil {
		return err
	}
	var rules [][]string
	for _, v := range rIDs {
		rules = append(rules, []string{roleID, strconv.Itoa(int(v))})
	}
	e := cs.Casbin()
	_, err = e.AddGroupingPolicies(rules)
	if err != nil {
		return err
	}
	return e.SavePolicy()
}

// ClearGroupingPolicies 清除用户组的权限
func (cs *CasbinService) ClearGroupingPolicies(v int, g ...string) (bool, error) {
	e := cs.Casbin()
	isOk, err := e.RemoveFilteredGroupingPolicy(v, g...)
	if err != nil {
		return false, err
	}
	return isOk, nil
}

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

func (casbinService *CasbinService) Casbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		a, err := gormadapter.NewAdapterByDB(global.FC_DB)
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			return
		}
		text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub, p.sub) && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		m, err := model.NewModelFromString(text)
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			return
		}
		syncedCachedEnforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
		syncedCachedEnforcer.SetExpireTime(60 * 60)
		_ = syncedCachedEnforcer.LoadPolicy()
	})
	return syncedCachedEnforcer
}
