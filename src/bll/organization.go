package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
)

// Organization ...
type Organization struct {
	ms *model.Models
}

// AddOrg ...
func (b *Organization) AddOrg(ctx context.Context, org string) (*tpl.SuccessResponseType, error) {
	ok, err := b.ms.Organization.AddOrg(ctx, org)
	if err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: ok}, nil
}

// UpdateOrgStatus ...
func (b *Organization) UpdateOrgStatus(ctx context.Context, org string, status int) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Organization.UpdateOrgStatus(ctx, org, status); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// DeleteOrg ...
func (b *Organization) DeleteOrg(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListOrgs ...
func (b *Organization) ListOrgs(ctx context.Context, pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	data, err := b.ms.Organization.ListOrgs(ctx, pg.PageSize, pg.Skip, pg.PageToken)
	if err != nil {
		return nil, err
	}
	res := &tpl.SuccessResponseType{Result: data, NextToken: ""}
	if len(data) >= pg.PageSize {
		res.NextToken = data[len(data)-1].UID
	}
	return res, nil
}

// ListSubjectOrgs ...
func (b *Organization) ListSubjectOrgs(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AddOU ...
func (b *Organization) AddOU(ctx context.Context, org string, input tpl.OrganizationAddOUInput) (*tpl.SuccessResponseType, error) {
	ok, err := b.ms.Organization.AddOU(ctx, org, input)
	if err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: ok}, nil
}

// UpdateOUParent ...
func (b *Organization) UpdateOUParent(ctx context.Context, org string, input tpl.OrganizationUpdateOUParentInput) (*tpl.SuccessResponseType, error) {
	err := b.ms.Organization.UpdateOUParent(ctx, org, input)
	if err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// UpdateOUStatus ...
func (b *Organization) UpdateOUStatus(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// UpdateOUTerms ...
func (b *Organization) UpdateOUTerms(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// DeleteOU ...
func (b *Organization) DeleteOU(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListOUs ...
func (b *Organization) ListOUs(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListSubjectOUs ...
func (b *Organization) ListSubjectOUs(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// SearchOUs ...
func (b *Organization) SearchOUs(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// BatchAddMember ...
func (b *Organization) BatchAddMember(ctx context.Context, org string, input tpl.OrganizationBatchAddMemberInput) (*tpl.SuccessResponseType, error) {
	subjects := make([]string, 0, len(input.Subjects))
	for _, sub := range input.Subjects {
		subjects = append(subjects, sub.Sub)
	}

	subs, err := b.ms.Subject.AcquireUIDsOrAdd(ctx, subjects)
	if err != nil {
		return nil, err
	}

	for i := range input.Subjects {
		input.Subjects[i].UID = tpl.GetSubjectUID(subs, input.Subjects[i].Sub)
	}
	if err := b.ms.Organization.BatchAddMember(ctx, org, input); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// UpdateMemberStatus ...
func (b *Organization) UpdateMemberStatus(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// UpdateMemberTerms ...
func (b *Organization) UpdateMemberTerms(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// RemoveMember ...
func (b *Organization) RemoveMember(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListMembers ...
func (b *Organization) ListMembers(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// SearchMember ...
func (b *Organization) SearchMember(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// BatchAddOUMember ...
func (b *Organization) BatchAddOUMember(ctx context.Context, org string, input tpl.OrganizationBatchAddOUMemberInput) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Organization.BatchAddOUMember(ctx, org, input); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// RemoveOUMember ...
func (b *Organization) RemoveOUMember(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListOUMembers ...
func (b *Organization) ListOUMembers(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListOUDescendantMembers ...
func (b *Organization) ListOUDescendantMembers(ctx context.Context, org string, ex tpl.Extensions) (*tpl.SuccessResponseType, error) {
	return nil, nil
}
