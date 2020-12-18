package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
)

// Permission ...
type Permission struct {
	ms *model.Models
}

// BatchAdd 批量添加权限
func (b *Permission) BatchAdd(ctx context.Context, tenant tpl.Tenant, permissions []string) (
	*tpl.SuccessResponseType, error) {
	if err := b.ms.Permission.BatchAdd(ctx, tenant, permissions); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// Delete 删除权限
func (b *Permission) Delete(ctx context.Context, tenant tpl.Tenant, permission string) (
	*tpl.SuccessResponseType, error) {
	if err := b.ms.Permission.Delete(ctx, tenant, permission); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// List 列出该系统当前指定资源类型的权限，当 resource 为空时列出所有权限
func (b *Permission) List(ctx context.Context, tenant tpl.Tenant, resources []string, pg tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	data, err := b.ms.Permission.List(ctx, tenant, resources, pg.PageSize, pg.Skip, pg.PageToken)
	if err != nil {
		return nil, err
	}
	res := &tpl.SuccessResponseType{Result: data, NextToken: ""}
	if len(data) >= pg.PageSize {
		res.NextToken = data[len(data)-1].UID
	}
	return res, nil
}
