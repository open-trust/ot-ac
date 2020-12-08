package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/schema"
	"github.com/open-trust/ot-ac/src/tpl"
)

// Scope ...
type Scope struct {
	ms *model.Models
}

// Add 创建范围约束
func (b *Scope) Add(ctx context.Context, tenant *schema.Tenant, scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// Delete 删除范围约束
func (b *Scope) Delete(ctx context.Context, tenant *schema.Tenant, scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// DeleteAll 删除范围约束及范围内的所有 Unit 和 Object
func (b *Scope) DeleteAll(ctx context.Context, tenant *schema.Tenant, scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// UpdateStatus 更新范围约束的状态，-1 表示停用
func (b *Scope) UpdateStatus(ctx context.Context, tenant *schema.Tenant, scope *tpl.Target, status int) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// List 列出该系统当前所有指定目标类型的范围约束
func (b *Scope) List(ctx context.Context, tenant *schema.Tenant, targetType string, pg *tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListUnits 列出范围约束下指定目标类型的直属的管理单元
func (b *Scope) ListUnits(ctx context.Context, tenant *schema.Tenant, scope *tpl.Target, targetType string, pg *tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListObjects 列出范围约束下指定目标类型的直属的资源对象
func (b *Scope) ListObjects(ctx context.Context, tenant *schema.Tenant, scope *tpl.Target, targetType string, pg *tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}
