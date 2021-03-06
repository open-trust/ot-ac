package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
)

// AC ...
type AC struct {
	ms *model.Models
}

// CheckUnit 检查请求主体到指定管理单元有没有指定权限
func (b *AC) CheckUnit(ctx context.Context, tenant tpl.Tenant, subject string,
	unit tpl.Target, permissions []string, withOrganization bool) (*tpl.SuccessResponseType, error) {
	res, err := b.ms.AC.CheckUnit(ctx, tenant, subject, unit, permissions, withOrganization)
	return &tpl.SuccessResponseType{Result: res}, err
}

// CheckScope 检查请求主体到指定范围约束有没有指定权限
func (b *AC) CheckScope(ctx context.Context, tenant tpl.Tenant, subject string,
	scope tpl.Target, permissions []string, withOrganization bool) (*tpl.SuccessResponseType, error) {
	res, err := b.ms.AC.CheckScope(ctx, tenant, subject, scope, permissions, withOrganization)
	return &tpl.SuccessResponseType{Result: res}, err
}

// CheckObject 检查请求主体通过 Scope 或 Unit-Object 的连接关系到指定资源对象有没有指定权限，如果 ignoreScope 为 true，则要求必须有 Unit-Object 的连接关系
func (b *AC) CheckObject(ctx context.Context, tenant tpl.Tenant, subject string,
	object tpl.Target, permissions []string, withOrganization, ignoreScope bool) (*tpl.SuccessResponseType, error) {
	res, err := b.ms.AC.CheckObject(ctx, tenant, subject, object, permissions, withOrganization, ignoreScope)
	return &tpl.SuccessResponseType{Result: res}, err
}

// ListPermissionsByUnit 列出请求主体到指定管理单元的符合 resource 的权限，如果未指定管理单元，则会查询请求主体能触达的所有管理单元，如果 resources 为空，则会列出所有触达的有效权限
func (b *AC) ListPermissionsByUnit(ctx context.Context, tenant tpl.Tenant, subject string,
	unit tpl.Target, resources []string, withOrganization bool, pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListPermissionsByScope 列出请求主体到指定范围约束的符合 resource 的权限，如果 resources 为空，则会列出所有触达的有效权限
func (b *AC) ListPermissionsByScope(ctx context.Context, tenant tpl.Tenant, subject string,
	scope tpl.Target, resources []string, withOrganization bool, pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListPermissionsByObject 列出请求主体到指定资源对象的符合 resource 的权限，如果 resources 为空，则会列出所有触达的有效权限
func (b *AC) ListPermissionsByObject(ctx context.Context, tenant tpl.Tenant, subject string,
	object tpl.Target, resources []string, withOrganization bool, ignoreScope bool, pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListObject 列出请求主体在指定资源对象中能触达的所有指定类型的子孙资源对象
// depth 定义对 targetType 类型资源对象的递归查询深度，而不是指定 object 到 targetType 类型资源对象的深度，默认对 targetType 类型资源对象查到底
func (b *AC) ListObject(ctx context.Context, tenant tpl.Tenant, subject string,
	object tpl.Target, permission, targetType string, withOrganization bool, ignoreScope bool, depth int, pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// SearchObject 根据关键词，在指定资源对象的子孙资源对象中，对请求主体能触达的所有指定类型的资源对象中进行搜索，term 为空不匹配任何资源对象
func (b *AC) SearchObject(ctx context.Context, tenant tpl.Tenant, subject string,
	object tpl.Target, permission, targetType, term string, withOrganization bool, ignoreScope bool, pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}
