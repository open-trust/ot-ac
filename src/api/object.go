package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/middleware"
	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/teambition/gear"
)

// Object ..
type Object struct {
	blls *bll.Blls
}

// BatchAdd 批量添加资源对象，当检测到将形成环时会返回 400 错误
func (a *Object) BatchAdd(ctx *gear.Context) error {
	input := tpl.TargetBatchAddInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Object.BatchAdd(model.ContextWithPrefer(ctx), *tenant, input.Targets, input.Parent, input.Scope)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// AssignParent 建立资源对象与父级对象的关系，当检测到将会形成环时会返回 400 错误
func (a *Object) AssignParent(ctx *gear.Context) error {
	return nil
}

// AssignScope 建立资源对象与范围约束的关系
func (a *Object) AssignScope(ctx *gear.Context) error {
	return nil
}

// RemoveParent 清除资源对象与父级对象的关系
func (a *Object) RemoveParent(ctx *gear.Context) error {
	return nil
}

// RemoveScope 清除资源对象与范围约束的关系
func (a *Object) RemoveScope(ctx *gear.Context) error {
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

// AddPermissions 给资源对象添加可透传的权限，权限必须预先存在
func (a *Object) AddPermissions(ctx *gear.Context) error {
	input := tpl.ObjectAddPermissionsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Object.AddPermissions(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Permissions)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdatePermissions 覆盖资源对象可透传的权限，权限必须预先存在，当 permissions 为空时会清空权限
func (a *Object) UpdatePermissions(ctx *gear.Context) error {
	return nil
}

// RemovePermissions 移除资源对象可透传的权限
func (a *Object) RemovePermissions(ctx *gear.Context) error {
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
