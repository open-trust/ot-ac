package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/teambition/gear"
)

// Admin ..
type Admin struct {
	blls *bll.Blls
}

// AddTenant ...
func (a *Admin) AddTenant(ctx *gear.Context) error {
	input := tpl.TenantAddInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Admin.AddTenant(model.ContextWithPrefer(ctx), input.Tenant)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateTenantStatus ...
func (a *Admin) UpdateTenantStatus(ctx *gear.Context) error {
	input := tpl.TenantAddInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Admin.UpdateTenantStatus(model.ContextWithPrefer(ctx), input.Tenant, input.Status)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// DeleteTenant ...
func (a *Admin) DeleteTenant(ctx *gear.Context) error {
	input := tpl.TenantAddInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Admin.DeleteTenant(model.ContextWithPrefer(ctx), input.Tenant)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// ListTenants ...
func (a *Admin) ListTenants(ctx *gear.Context) error {
	input := tpl.Pagination{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}
	res, err := a.blls.Admin.ListTenants(ctx, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// BatchAddSubjects ...
func (a *Admin) BatchAddSubjects(ctx *gear.Context) error {
	input := tpl.SubjectsInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Admin.BatchAddSubjects(model.ContextWithPrefer(ctx), input.Subjects)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateSubjectStatus ...
func (a *Admin) UpdateSubjectStatus(ctx *gear.Context) error {
	input := tpl.SubjectUpdateInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Admin.UpdateSubjectStatus(model.ContextWithPrefer(ctx), input.Sub, input.Status)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// ListSubjects ...
func (a *Admin) ListSubjects(ctx *gear.Context) error {
	input := tpl.Pagination{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}
	res, err := a.blls.Admin.ListSubjects(ctx, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}
