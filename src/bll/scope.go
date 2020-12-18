package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
)

// Scope ...
type Scope struct {
	ms *model.Models
}

// Add 创建范围约束
func (b *Scope) Add(ctx context.Context, tenant tpl.Tenant, scope tpl.Target) (*tpl.SuccessResponseType, error) {
	ok, err := b.ms.Scope.Add(ctx, tenant, tpl.Scope{
		Status:     0,
		TargetID:   scope.ID,
		TargetType: scope.Type,
	})
	if err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: ok}, nil
}

// Delete 删除范围约束
func (b *Scope) Delete(ctx context.Context, tenant tpl.Tenant, scope tpl.Target) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Scope.Delete(ctx, tenant, scope); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// DeleteAll 删除范围约束及范围内的所有 Unit 和 Object
func (b *Scope) DeleteAll(ctx context.Context, tenant tpl.Tenant, scope tpl.Target) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Scope.DeleteAll(ctx, tenant, scope); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// UpdateStatus 更新范围约束的状态，-1 表示停用
func (b *Scope) UpdateStatus(ctx context.Context, tenant tpl.Tenant, scope tpl.Target, status int) (*tpl.SuccessResponseType, error) {
	data := tpl.Scope{
		Status:     0,
		TargetID:   scope.ID,
		TargetType: scope.Type,
	}
	if err := b.ms.Scope.UpdateStatus(ctx, tenant, scope, status); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: data}, nil
}

// List 列出该系统当前所有指定目标类型的范围约束
func (b *Scope) List(ctx context.Context, tenant tpl.Tenant, targetType string, pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	data, err := b.ms.Scope.List(ctx, tenant, targetType, pg.PageSize, pg.Skip, pg.PageToken)
	if err != nil {
		return nil, err
	}
	res := &tpl.SuccessResponseType{Result: data, NextToken: ""}
	if len(data) >= pg.PageSize {
		res.NextToken = data[len(data)-1].UID
	}
	return res, nil
}

// ListUnits 列出范围约束下指定目标类型的直属的管理单元
func (b *Scope) ListUnits(ctx context.Context, tenant tpl.Tenant, scope tpl.Target, targetType string, pg tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	data, err := b.ms.Scope.ListUnits(ctx, tenant, scope, targetType, pg.PageSize, pg.Skip, pg.PageToken)
	if err != nil {
		return nil, err
	}
	res := &tpl.SuccessResponseType{Result: data, NextToken: ""}
	if len(data) >= pg.PageSize {
		res.NextToken = data[len(data)-1].UID
	}
	return res, nil
}

// ListObjects 列出范围约束下指定目标类型的直属的资源对象
func (b *Scope) ListObjects(ctx context.Context, tenant tpl.Tenant, scope tpl.Target, targetType string, pg tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	data, err := b.ms.Scope.ListObjects(ctx, tenant, scope, targetType, pg.PageSize, pg.Skip, pg.PageToken)
	if err != nil {
		return nil, err
	}
	res := &tpl.SuccessResponseType{Result: data, NextToken: ""}
	if len(data) >= pg.PageSize {
		res.NextToken = data[len(data)-1].UID
	}
	return res, nil
}
