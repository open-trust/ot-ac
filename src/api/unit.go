package api

import (
	"strconv"

	"github.com/open-trust/ot-ac/src/bll"
	"github.com/teambition/gear"
)

// Unit ..
type Unit struct {
	blls *bll.Blls
}

// Serve ..
func (a *Unit) Serve(ctx *gear.Context) error {
	action := ctx.Param("action")
	switch action {
	case "BatchAdd":
		return a.BatchAdd(ctx)
	case "AssignParent":
		return a.AssignParent(ctx)
	case "AssignScope":
		return a.AssignScope(ctx)
	case "AssignObject":
		return a.AssignObject(ctx)
	case "ClearParent":
		return a.ClearParent(ctx)
	case "ClearScope":
		return a.ClearScope(ctx)
	case "ClearObject":
		return a.ClearObject(ctx)
	case "Delete":
		return a.Delete(ctx)
	case "UpdateStatus":
		return a.UpdateStatus(ctx)
	case "AddSubjects":
		return a.AddSubjects(ctx)
	case "ClearSubjects":
		return a.ClearSubjects(ctx)
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
	case "ListSubjects":
		return a.ListSubjects(ctx)
	case "ListDescendantSubjects":
		return a.ListDescendantSubjects(ctx)
	case "GetDAG":
		return a.GetDAG(ctx)
	}
	return gear.ErrBadRequest.WithMsgf("unknown action %s", strconv.Quote(action))
}

// BatchAdd 批量添加管理单元，当检测到将形成环时会返回 400 错误
func (a *Unit) BatchAdd(ctx *gear.Context) error {
	return nil
}

// AssignParent 建立管理单元与父级管理单元的关系，当检测到将形成环时会返回 400 错误
func (a *Unit) AssignParent(ctx *gear.Context) error {
	return nil
}

// AssignScope 建立管理单元与范围约束的关系
func (a *Unit) AssignScope(ctx *gear.Context) error {
	return nil
}

// AssignObject 建立管理单元与资源对象的关系
func (a *Unit) AssignObject(ctx *gear.Context) error {
	return nil
}

// ClearParent 清除管理单元与父级对象的关系
func (a *Unit) ClearParent(ctx *gear.Context) error {
	return nil
}

// ClearScope 清除管理单元与范围约束的关系
func (a *Unit) ClearScope(ctx *gear.Context) error {
	return nil
}

// ClearObject 清除管理单元与资源对象的关系
func (a *Unit) ClearObject(ctx *gear.Context) error {
	return nil
}

// Delete 删除管理单元及其所有子孙管理单元和链接关系
func (a *Unit) Delete(ctx *gear.Context) error {
	return nil
}

// UpdateStatus 更新管理单元的状态，-1 表示停用
func (a *Unit) UpdateStatus(ctx *gear.Context) error {
	return nil
}

// AddSubjects 管理单元批量添加请求主体，当请求主体不存在时会自动创建
func (a *Unit) AddSubjects(ctx *gear.Context) error {
	return nil
}

// ClearSubjects 管理单元批量移除请求主体
func (a *Unit) ClearSubjects(ctx *gear.Context) error {
	return nil
}

// UpdatePermissions 给管理单元添加权限，权限必须预先存在
func (a *Unit) UpdatePermissions(ctx *gear.Context) error {
	return nil
}

// OverridePermissions 覆盖管理单元的权限，权限必须预先存在，当 permissions 为空时会清空权限
func (a *Unit) OverridePermissions(ctx *gear.Context) error {
	return nil
}

// ClearPermissions 移除管理单元的权限
func (a *Unit) ClearPermissions(ctx *gear.Context) error {
	return nil
}

// ListChildren 列出管理单元的指定目标类型的子级管理单元，不包含 status 为 -1 的节点
func (a *Unit) ListChildren(ctx *gear.Context) error {
	return nil
}

// ListDescendant 列出管理单元的指定目标类型的所有子孙管理单元，不包含 status 为 -1 的管理单元
// depth 定义对 targetType 类型管理单元的递归查询深度，而不是指定 unit 到 targetType 类型管理单元的深度，默认对 targetType 类型管理单元查到底
func (a *Unit) ListDescendant(ctx *gear.Context) error {
	return nil
}

// ListPermissions 列出管理单元的直属权限
func (a *Unit) ListPermissions(ctx *gear.Context) error {
	return nil
}

// ListSubjects 列出管理单元的直属请求主体，不包含 status 为 -1 的请求主体
func (a *Unit) ListSubjects(ctx *gear.Context) error {
	return nil
}

// ListDescendantSubjects 列出管理单元及子孙管理单元下所有的请求主体，不包含 status 为 -1 的请求主体
func (a *Unit) ListDescendantSubjects(ctx *gear.Context) error {
	return nil
}

// GetDAG 根据 start 和 ends 找出一个 DAG，其中 start 可以为 Subject 或 Unit，ends 为 0 到多个 Unit，不包含 status 为 -1 的节点
func (a *Unit) GetDAG(ctx *gear.Context) error {
	return nil
}
