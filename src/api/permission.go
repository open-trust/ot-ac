package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/middleware"
	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/teambition/gear"
)

// Permission ..
type Permission struct {
	blls *bll.Blls
}

// BatchAdd 批量添加权限
func (a *Permission) BatchAdd(ctx *gear.Context) error {
	input := &tpl.BatchAddPermissionsInput{}
	if err := ctx.ParseBody(input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Permission.BatchAdd(model.ContextWithPrefer(ctx), tenant, input.Permissions)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Delete 删除权限
func (a *Permission) Delete(ctx *gear.Context) error {
	input := &tpl.DeletePermissionInput{}
	if err := ctx.ParseBody(input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Permission.Delete(model.ContextWithPrefer(ctx), tenant, input.Permission)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// List 列出该系统当前指定资源类型的权限，当 resource 为空时列出所有权限
func (a *Permission) List(ctx *gear.Context) error {
	input := &tpl.ListPermissionsInput{}
	if err := ctx.ParseBody(input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.Permission.List(model.ContextWithPrefer(ctx), tenant, input.Resources, &input.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}
