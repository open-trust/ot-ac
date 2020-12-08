package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/schema"
	"github.com/open-trust/ot-ac/src/tpl"
)

// Unit ...
type Unit struct {
	ms *model.Models
}

// BatchAdd 批量添加管理单元，当检测到将形成环时会返回 400 错误
func (b *Unit) BatchAdd(ctx context.Context, tenant *schema.Tenant, units []*tpl.Target, parent *tpl.Target,
	scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AssignParent 建立管理单元与父级管理单元的关系，当检测到将形成环时会返回 400 错误
func (b *Unit) AssignParent(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, parent *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AssignScope 建立管理单元与范围约束的关系
func (b *Unit) AssignScope(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, scope *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AssignObject 建立管理单元与资源对象的关系
func (b *Unit) AssignObject(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, object *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ClearParent 清除管理单元与父级对象的关系
func (b *Unit) ClearParent(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, parent *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ClearScope 清除管理单元与范围约束的关系
func (b *Unit) ClearScope(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, scope *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ClearObject 清除管理单元与资源对象的关系
func (b *Unit) ClearObject(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, object *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// Delete 删除管理单元及其所有子孙管理单元和链接关系
func (b *Unit) Delete(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// UpdateStatus 更新管理单元的状态，-1 表示停用
func (b *Unit) UpdateStatus(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, status int) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AddSubjects 管理单元批量添加请求主体，当请求主体不存在时会自动创建
func (b *Unit) AddSubjects(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, subjects []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ClearSubjects 管理单元批量移除请求主体
func (b *Unit) ClearSubjects(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, subjects []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// UpdatePermissions 给管理单元添加权限，权限必须预先存在
func (b *Unit) UpdatePermissions(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, permissions []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// OverridePermissions 覆盖管理单元的权限，权限必须预先存在，当 permissions 为空时会清空权限
func (b *Unit) OverridePermissions(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, permissions []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ClearPermissions 移除管理单元的权限
func (b *Unit) ClearPermissions(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, permissions []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListChildren 列出管理单元的指定目标类型的子级管理单元，不包含 status 为 -1 的节点
func (b *Unit) ListChildren(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, targetType string,
	pg *tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListDescendant 列出管理单元的指定目标类型的所有子孙管理单元，不包含 status 为 -1 的管理单元
// depth 定义对 targetType 类型管理单元的递归查询深度，而不是指定 unit 到 targetType 类型管理单元的深度，默认对 targetType 类型管理单元查到底
func (b *Unit) ListDescendant(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, targetType string, depth int,
	pg *tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListPermissions 列出管理单元的直属权限
func (b *Unit) ListPermissions(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, pg *tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListSubjects 列出管理单元的直属请求主体，不包含 status 为 -1 的请求主体
func (b *Unit) ListSubjects(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, pg *tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListDescendantSubjects 列出管理单元及子孙管理单元下所有的请求主体，不包含 status 为 -1 的请求主体
func (b *Unit) ListDescendantSubjects(ctx context.Context, tenant *schema.Tenant, unit *tpl.Target, pg *tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// GetDAG 根据 start 和 ends 找出一个 DAG，其中 start 可以为 Subject 或 Unit，ends 为 0 到多个 Unit，不包含 status 为 -1 的节点
func (b *Unit) GetDAG(ctx context.Context, tenant *schema.Tenant, start *tpl.Target, ends []*tpl.Target, pg *tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}
