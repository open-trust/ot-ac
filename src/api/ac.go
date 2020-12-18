package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/middleware"
	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/teambition/gear"
)

// AC ..
type AC struct {
	blls *bll.Blls
}

// CheckUnit 检查请求主体到指定管理单元有没有指定权限
func (a *AC) CheckUnit(ctx *gear.Context) error {
	input := tpl.ACCheckPermissionsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.AC.CheckUnit(model.ContextWithPrefer(ctx), *tenant, input.Subject, input.Target, input.Permissions, input.WithOrganization)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// CheckScope ...
func (a *AC) CheckScope(ctx *gear.Context) error {
	input := tpl.ACCheckPermissionsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.AC.CheckScope(model.ContextWithPrefer(ctx), *tenant, input.Subject, input.Target, input.Permissions, input.WithOrganization)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// CheckObject ...
func (a *AC) CheckObject(ctx *gear.Context) error {
	input := tpl.ACCheckPermissionsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	tenant, err := middleware.TenantFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := a.blls.AC.CheckObject(model.ContextWithPrefer(ctx), *tenant, input.Subject, input.Target, input.Permissions, input.WithOrganization, input.IgnoreScope)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

func (a *AC) ListPermissionsByUnit(ctx *gear.Context) error {
	return nil
}

func (a *AC) ListPermissionsByScope(ctx *gear.Context) error {
	return nil
}

func (a *AC) ListPermissionsByObject(ctx *gear.Context) error {
	return nil
}

func (a *AC) ListObject(ctx *gear.Context) error {
	return nil
}

func (a *AC) SearchObject(ctx *gear.Context) error {
	return nil
}
