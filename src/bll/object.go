package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/schema"
	"github.com/open-trust/ot-ac/src/tpl"
)

// Object ...
type Object struct {
	ms *model.Models
}

// BatchAdd 批量添加资源对象，当检测到将形成环时会返回 400 错误
func (b *Object) BatchAdd(ctx context.Context, tenant *schema.Tenant, objects []*tpl.Target, parent *tpl.Target,
	scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AddWithUnit 添加资源对象，并同时添加对应的管理单元和建立连接关系，当检测到将形成环时会返回 400 错误
func (b *Object) AddWithUnit(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, parent *tpl.Target,
	scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AssignParent 建立资源对象与父级对象的关系，当检测到将会形成环时会返回 400 错误
func (b *Object) AssignParent(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, parent *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AssignScope 建立资源对象与范围约束的关系
func (b *Object) AssignScope(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, scope *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ClearParent 清除资源对象与父级对象的关系
func (b *Object) ClearParent(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, parent *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ClearScope 清除资源对象与范围约束的关系
func (b *Object) ClearScope(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, scope *tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// Delete 删除资源对象及其所有子孙资源对象和链接关系
func (b *Object) Delete(ctx context.Context, tenant *schema.Tenant, object *tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// UpdateTerms 更新资源对象的搜索关键词
func (b *Object) UpdateTerms(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, terms []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// UpdatePermissions 给资源对象添加可透传的权限，权限必须预先存在
func (b *Object) UpdatePermissions(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, permissions []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// OverridePermissions 覆盖资源对象可透传的权限，权限必须预先存在，当 permissions 为空时会清空权限
func (b *Object) OverridePermissions(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, permissions []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ClearPermissions 移除资源对象可透传的权限
func (b *Object) ClearPermissions(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, permissions []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListChildren 列出资源对象的指定目标类型的子级资源对象
func (b *Object) ListChildren(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, targetType string,
	pg *tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListDescendant 列出资源对象的所有指定目标类型的子孙资源对象
// depth 定义对 *tpl.TargetType 类型资源对象的递归查询深度，而不是指定 object 到 *tpl.TargetType 类型资源对象的深度，默认对 *tpl.TargetType 类型资源对象查到底
func (b *Object) ListDescendant(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, targetType string,
	depth int, pg *tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListPermissions 列出资源对象可透传的权限
func (b *Object) ListPermissions(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, pg *tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// GetDAG 根据 start 和 ends 找出一个 DAG，其中 start 为 Object，ends 为 0 到多个 Object
func (b *Object) GetDAG(ctx context.Context, tenant *schema.Tenant, start *tpl.Target, ends []*tpl.Target,
	pg *tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// Search 根据关键词在资源对象的所有指定类型的子孙资源对象中进行搜索，term 为空不匹配任何资源对象
func (b *Object) Search(ctx context.Context, tenant *schema.Tenant, object *tpl.Target, targetType string,
	term string, pg *tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}
