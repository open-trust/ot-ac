package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
	otgo "github.com/open-trust/ot-go-lib"
	"github.com/teambition/gear"
)

// Admin ...
type Admin struct {
	ms *model.Models
}

// AddTenant 添加租户
func (b *Admin) AddTenant(ctx context.Context, tenant otgo.OTID) (*tpl.SuccessResponseType, error) {
	ok, err := b.ms.Tenant.Add(ctx, tpl.Tenant{
		Tenant: tenant.String(),
		Status: 0,
	})
	if err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: ok}, nil
}

// UpdateTenantStatus 更新租户，-1 表示停用
func (b *Admin) UpdateTenantStatus(ctx context.Context, tenant otgo.OTID, status int) (*tpl.SuccessResponseType, error) {
	doc, err := b.ms.Tenant.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}
	data := tpl.Tenant{
		Tenant: tenant.String(),
		Status: status,
	}

	if status != doc.Status {
		if err := b.ms.Tenant.Update(ctx, data); err != nil {
			return nil, err
		}
	}
	return &tpl.SuccessResponseType{Result: data}, nil
}

// DeleteTenant 删除租户及其名下所有数据，status 必须为 -1 才能删除
func (b *Admin) DeleteTenant(ctx context.Context, tenant otgo.OTID) (*tpl.SuccessResponseType, error) {
	doc, err := b.ms.Tenant.Get(ctx, tenant)
	if err != nil {
		return &tpl.SuccessResponseType{Result: true}, nil
	}
	if doc.Status >= 0 {
		return nil, gear.ErrPreconditionRequired.WithMsgf("tenant %s should be disabled before deleting", tenant.String())
	}

	if err := b.ms.Tenant.Delete(ctx, tenant); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// ListTenants 列出系统中的租户
func (b *Admin) ListTenants(ctx context.Context, pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	data, err := b.ms.Tenant.List(ctx, pg.PageSize, pg.Skip, pg.PageToken)
	if err != nil {
		return nil, err
	}
	res := &tpl.SuccessResponseType{Result: data, NextToken: ""}
	if len(data) >= pg.PageSize {
		res.NextToken = data[len(data)-1].UID
	}
	return res, nil
}

// BatchAddSubjects 批量添加请求主体
func (b *Admin) BatchAddSubjects(ctx context.Context, subjects []string) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Subject.BatchAdd(ctx, subjects); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// UpdateSubjectStatus 更新请求主体，-1 表示停用
func (b *Admin) UpdateSubjectStatus(ctx context.Context, subject string, status int) (*tpl.SuccessResponseType, error) {
	data := tpl.Subject{
		Sub:    subject,
		Status: status,
	}
	if err := b.ms.Subject.Update(ctx, data); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: data}, nil
}

// ListSubjects 列出系统中的请求主体
func (b *Admin) ListSubjects(ctx context.Context, pg tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	data, err := b.ms.Subject.List(ctx, pg.PageSize, pg.Skip, pg.PageToken)
	if err != nil {
		return nil, err
	}
	res := &tpl.SuccessResponseType{Result: data, NextToken: ""}
	if len(data) >= pg.PageSize {
		res.NextToken = data[len(data)-1].UID
	}
	return res, nil
}
