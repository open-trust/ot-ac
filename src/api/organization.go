package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/teambition/gear"
)

// Organization ..
type Organization struct {
	blls *bll.Blls
}

// AddOrg ...
func (a *Organization) AddOrg(ctx *gear.Context) error {
	input := tpl.OrganizationInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Organization.AddOrg(model.ContextWithPrefer(ctx), input.Org)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateOrgStatus ...
func (a *Organization) UpdateOrgStatus(ctx *gear.Context) error {
	input := tpl.OrganizationStatusInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Organization.UpdateOrgStatus(model.ContextWithPrefer(ctx), input.Org, input.Status)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// DeleteOrg ...
func (a *Organization) DeleteOrg(ctx *gear.Context) error {
	return nil
}

// ListOrgs ...
func (a *Organization) ListOrgs(ctx *gear.Context) error {
	input := tpl.Pagination{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}
	res, err := a.blls.Organization.ListOrgs(ctx, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// ListSubjectOrgs ...
func (a *Organization) ListSubjectOrgs(ctx *gear.Context) error {
	return nil
}

// AddOU ...
func (a *Organization) AddOU(ctx *gear.Context) error {
	input := tpl.OrganizationAddOUInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Organization.AddOU(model.ContextWithPrefer(ctx), input.Org, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateOUParent ...
func (a *Organization) UpdateOUParent(ctx *gear.Context) error {
	input := tpl.OrganizationUpdateOUParentInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Organization.UpdateOUParent(model.ContextWithPrefer(ctx), input.Org, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateOUStatus ...
func (a *Organization) UpdateOUStatus(ctx *gear.Context) error {
	return nil
}

// UpdateOUTerms ...
func (a *Organization) UpdateOUTerms(ctx *gear.Context) error {
	return nil
}

// DeleteOU ...
func (a *Organization) DeleteOU(ctx *gear.Context) error {
	return nil
}

// ListOUs ...
func (a *Organization) ListOUs(ctx *gear.Context) error {
	return nil
}

// ListSubjectOUs ...
func (a *Organization) ListSubjectOUs(ctx *gear.Context) error {
	return nil
}

// SearchOUs ...
func (a *Organization) SearchOUs(ctx *gear.Context) error {
	return nil
}

// BatchAddMember ...
func (a *Organization) BatchAddMember(ctx *gear.Context) error {
	input := tpl.OrganizationBatchAddMemberInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Organization.BatchAddMember(model.ContextWithPrefer(ctx), input.Org, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateMemberStatus ...
func (a *Organization) UpdateMemberStatus(ctx *gear.Context) error {
	return nil
}

// UpdateMemberTerms ...
func (a *Organization) UpdateMemberTerms(ctx *gear.Context) error {
	return nil
}

// RemoveMember ...
func (a *Organization) RemoveMember(ctx *gear.Context) error {
	return nil
}

// ListMembers ...
func (a *Organization) ListMembers(ctx *gear.Context) error {
	return nil
}

// SearchMember ...
func (a *Organization) SearchMember(ctx *gear.Context) error {
	return nil
}

// BatchAddOUMember ...
func (a *Organization) BatchAddOUMember(ctx *gear.Context) error {
	input := tpl.OrganizationBatchAddOUMemberInput{}
	if err := ctx.ParseBody(&input); err != nil {
		return err
	}

	res, err := a.blls.Organization.BatchAddOUMember(model.ContextWithPrefer(ctx), input.Org, input)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// RemoveOUMember ...
func (a *Organization) RemoveOUMember(ctx *gear.Context) error {
	return nil
}

// ListOUMembers ...
func (a *Organization) ListOUMembers(ctx *gear.Context) error {
	return nil
}

// ListOUDescendantMembers ...
func (a *Organization) ListOUDescendantMembers(ctx *gear.Context) error {
	return nil
}
