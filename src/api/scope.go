package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/middleware"
	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/teambition/gear"
)

// Scope ..
type Scope struct {
	blls *bll.Blls
}

// Add 创建范围约束
func (a *Scope) Add(ctx *gear.Context) error {
	input := tpl.Target{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Scope.Add(model.ContextWithPrefer(ctx), *tenant, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Delete 删除范围约束
func (a *Scope) Delete(ctx *gear.Context) error {
	input := tpl.Target{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Scope.Delete(model.ContextWithPrefer(ctx), *tenant, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// DeleteAll 删除范围约束及范围内的所有 Unit 和 Object
func (a *Scope) DeleteAll(ctx *gear.Context) error {
	input := tpl.Target{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Scope.DeleteAll(model.ContextWithPrefer(ctx), *tenant, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateStatus 更新范围约束的状态，-1 表示停用
func (a *Scope) UpdateStatus(ctx *gear.Context) error {
	input := tpl.ScopeUpdateInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Scope.UpdateStatus(model.ContextWithPrefer(ctx), *tenant, input.Target, input.Status)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// List 列出该系统当前所有指定目标类型的范围约束
func (a *Scope) List(ctx *gear.Context) error {
	input := tpl.ScopeListInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Scope.List(model.ContextWithPrefer(ctx), *tenant, input.TargetType, input.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// ListUnits 列出范围约束下指定目标类型的直属的管理单元
func (a *Scope) ListUnits(ctx *gear.Context) error {
	input := tpl.ScopeListUnitObjectsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Scope.ListUnits(model.ContextWithPrefer(ctx), *tenant, input.Scope, input.TargetType, input.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// ListObjects 列出范围约束下指定目标类型的直属的资源对象
func (a *Scope) ListObjects(ctx *gear.Context) error {
	input := tpl.ScopeListUnitObjectsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Scope.ListObjects(model.ContextWithPrefer(ctx), *tenant, input.Scope, input.TargetType, input.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}
