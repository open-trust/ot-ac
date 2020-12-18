package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/middleware"
	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/teambition/gear"
)

// Unit ..
type Unit struct {
	blls *bll.Blls
}

// BatchAdd 批量添加管理单元，当检测到将形成环时会返回 400 错误
func (a *Unit) BatchAdd(ctx *gear.Context) error {
	input := tpl.TargetBatchAddInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.BatchAdd(model.ContextWithPrefer(ctx), *tenant, input.Targets, input.Parent, input.Scope)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// AddFromOrg ...
func (a *Unit) AddFromOrg(ctx *gear.Context) error {
	input := tpl.UnitAddFromOrgInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.AddFromOrg(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Org, input.Parent, input.Scope)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// AddFromOU ...
func (a *Unit) AddFromOU(ctx *gear.Context) error {
	input := tpl.UnitAddFromOUInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.AddFromOU(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Org, input.OU, input.Parent, input.Scope)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// AddFromMembers ...
func (a *Unit) AddFromMembers(ctx *gear.Context) error {
	input := tpl.UnitAddFromMembersInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.AddFromMembers(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Org, input.Subjects, input.Parent, input.Scope)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// AssignParent 建立管理单元与父级管理单元的关系，当检测到将形成环时会返回 400 错误
func (a *Unit) AssignParent(ctx *gear.Context) error {
	input := tpl.UnitAssignParentInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.AssignParent(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Parent)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// AssignScope 建立管理单元与范围约束的关系
func (a *Unit) AssignScope(ctx *gear.Context) error {
	input := tpl.UnitAssignScopeInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.AssignScope(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Scope)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// AssignObject 建立管理单元与资源对象的关系
func (a *Unit) AssignObject(ctx *gear.Context) error {
	input := tpl.UnitAssignObjectInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.AssignObject(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Object)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// RemoveParent 清除管理单元与父级对象的关系
func (a *Unit) RemoveParent(ctx *gear.Context) error {
	return nil
}

// RemoveScope 清除管理单元与范围约束的关系
func (a *Unit) RemoveScope(ctx *gear.Context) error {
	return nil
}

// RemoveObject 清除管理单元与资源对象的关系
func (a *Unit) RemoveObject(ctx *gear.Context) error {
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
	input := tpl.UnitAddSubjectsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.AddSubjects(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Subjects)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// RemoveSubjects 管理单元批量移除请求主体
func (a *Unit) RemoveSubjects(ctx *gear.Context) error {
	return nil
}

// AddPermissions 给管理单元添加权限，权限必须预先存在
func (a *Unit) AddPermissions(ctx *gear.Context) error {
	input := tpl.UnitAddPermissionsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Unit.AddPermissions(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Permissions)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdatePermissions 覆盖管理单元的权限，权限必须预先存在，当 permissions 为空时会清空权限
func (a *Unit) UpdatePermissions(ctx *gear.Context) error {
	return nil
}

// RemovePermissions 移除管理单元的权限
func (a *Unit) RemovePermissions(ctx *gear.Context) error {
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
