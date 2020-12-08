package api

import (
	"strconv"

	"github.com/open-trust/ot-ac/src/bll"
	"github.com/teambition/gear"
)

// Object ..
type Object struct {
	blls *bll.Blls
}

// Serve ..
func (a *Object) Serve(ctx *gear.Context) error {
	action := ctx.Param("action")
	switch action {
	case "BatchAdd":
		return a.BatchAdd(ctx)
	case "AddWithUnit":
		return a.AddWithUnit(ctx)
	case "AssignParent":
		return a.AssignParent(ctx)
	case "AssignScope":
		return a.AssignScope(ctx)
	case "ClearParent":
		return a.ClearParent(ctx)
	case "ClearScope":
		return a.ClearScope(ctx)
	case "Delete":
		return a.Delete(ctx)
	case "UpdateTerms":
		return a.UpdateTerms(ctx)
	case "UpdatePermissions":
		return a.UpdatePermissions(ctx)
	case "OverridePermissions":
		return a.OverridePermissions(ctx)
	case "ClearPermissions":
		return a.ClearPermissions(ctx)
	case "ListChildren":
		return a.ListChildren(ctx)
	case "ListDescendant":
		return a.ListDescendant(ctx)
	case "ListPermissions":
		return a.ListPermissions(ctx)
	case "GetDAG":
		return a.GetDAG(ctx)
	case "Search":
		return a.Search(ctx)
	}
	return gear.ErrBadRequest.WithMsgf("unknown action %s", strconv.Quote(action))
}

// BatchAdd 批量添加资源对象，当检测到将形成环时会返回 400 错误
func (a *Object) BatchAdd(ctx *gear.Context) error {
	return nil
}

// AddWithUnit 添加资源对象，并同时添加对应的管理单元和建立连接关系，当检测到将形成环时会返回 400 错误
func (a *Object) AddWithUnit(ctx *gear.Context) error {
	return nil
}

// AssignParent 建立资源对象与父级对象的关系，当检测到将会形成环时会返回 400 错误
func (a *Object) AssignParent(ctx *gear.Context) error {
	return nil
}

// AssignScope 建立资源对象与范围约束的关系
func (a *Object) AssignScope(ctx *gear.Context) error {
	return nil
}

// ClearParent 清除资源对象与父级对象的关系
func (a *Object) ClearParent(ctx *gear.Context) error {
	return nil
}

// ClearScope 清除资源对象与范围约束的关系
func (a *Object) ClearScope(ctx *gear.Context) error {
	return nil
}

// Delete 删除资源对象及其所有子孙资源对象和链接关系
func (a *Object) Delete(ctx *gear.Context) error {
	return nil
}

// UpdateTerms 更新资源对象的搜索关键词
func (a *Object) UpdateTerms(ctx *gear.Context) error {
	return nil
}

// UpdatePermissions 给资源对象添加可透传的权限，权限必须预先存在
func (a *Object) UpdatePermissions(ctx *gear.Context) error {
	return nil
}

// OverridePermissions 覆盖资源对象可透传的权限，权限必须预先存在，当 permissions 为空时会清空权限
func (a *Object) OverridePermissions(ctx *gear.Context) error {
	return nil
}

// ClearPermissions 移除资源对象可透传的权限
func (a *Object) ClearPermissions(ctx *gear.Context) error {
	return nil
}

// ListChildren 列出资源对象的指定目标类型的子级资源对象
func (a *Object) ListChildren(ctx *gear.Context) error {
	return nil
}

// ListDescendant 列出资源对象的所有指定目标类型的子孙资源对象
// depth 定义对 *tpl.TargetType 类型资源对象的递归查询深度，而不是指定 object 到 *tpl.TargetType 类型资源对象的深度，默认对 *tpl.TargetType 类型资源对象查到底
func (a *Object) ListDescendant(ctx *gear.Context) error {
	return nil
}

// ListPermissions 列出资源对象可透传的权限
func (a *Object) ListPermissions(ctx *gear.Context) error {
	return nil
}

// GetDAG 根据 start 和 ends 找出一个 DAG，其中 start 为 Object，ends 为 0 到多个 Object
func (a *Object) GetDAG(ctx *gear.Context) error {
	return nil
}

// Search 根据关键词在资源对象的所有指定类型的子孙资源对象中进行搜索，term 为空不匹配任何资源对象
func (a *Object) Search(ctx *gear.Context) error {
	return nil
}
