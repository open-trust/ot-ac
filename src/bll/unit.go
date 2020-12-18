package bll

import (
	"context"

	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/tpl"
)

// Unit ...
type Unit struct {
	ms *model.Models
}

// BatchAdd 批量添加管理单元，当检测到将形成环时会返回 400 错误
func (b *Unit) BatchAdd(ctx context.Context, tenant tpl.Tenant, units []tpl.Target, parent, scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Unit.BatchAdd(ctx, tenant, units, parent, scope); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// AddFromOrg 从目录服务的 Org 创建管理单元，当检测到将形成环时会返回 400 错误
func (b *Unit) AddFromOrg(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, org string, parent *tpl.Target, scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Unit.AddFromOrg(ctx, tenant, unit, org, parent, scope); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// AddFromOU 从目录服务的 OU 创建管理单元，当检测到将形成环时会返回 400 错误
func (b *Unit) AddFromOU(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, org, ou string, parent *tpl.Target, scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Unit.AddFromOU(ctx, tenant, unit, org, ou, parent, scope); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// AddFromMembers 从目录服务的 OU 创建管理单元，当检测到将形成环时会返回 400 错误
func (b *Unit) AddFromMembers(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, org string, subjects []string, parent *tpl.Target, scope *tpl.Target) (*tpl.SuccessResponseType, error) {
	if err := b.ms.Unit.AddFromMembers(ctx, tenant, unit, org, subjects, parent, scope); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// AssignParent 建立管理单元与父级管理单元的关系，当检测到将形成环时会返回 400 错误
func (b *Unit) AssignParent(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, parent tpl.Target) (
	*tpl.SuccessResponseType, error) {
	if err := b.ms.Unit.AssignParent(ctx, tenant, unit, parent); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// AssignScope 建立管理单元与范围约束的关系
func (b *Unit) AssignScope(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, scope tpl.Target) (
	*tpl.SuccessResponseType, error) {
	if err := b.ms.Unit.AssignScope(ctx, tenant, unit, scope); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// AssignObject 建立管理单元与资源对象的关系
func (b *Unit) AssignObject(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, object tpl.Target) (
	*tpl.SuccessResponseType, error) {
	if err := b.ms.Unit.AssignObject(ctx, tenant, unit, object); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// RemoveParent 清除管理单元与父级对象的关系
func (b *Unit) RemoveParent(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, parent tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// RemoveScope 清除管理单元与范围约束的关系
func (b *Unit) RemoveScope(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, scope tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// RemoveObject 清除管理单元与资源对象的关系
func (b *Unit) RemoveObject(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, object tpl.Target) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// Delete 删除管理单元及其所有子孙管理单元和链接关系
func (b *Unit) Delete(ctx context.Context, tenant tpl.Tenant, unit tpl.Target) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// UpdateStatus 更新管理单元的状态，-1 表示停用
func (b *Unit) UpdateStatus(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, status int) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AddSubjects 管理单元批量添加请求主体，当请求主体不存在时会自动创建
func (b *Unit) AddSubjects(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, input []string) (
	*tpl.SuccessResponseType, error) {
	subjects, err := b.ms.Subject.AcquireUIDsOrAdd(ctx, input)
	if err != nil {
		return nil, err
	}
	subjectUIDs := make([]string, 0, len(subjects))
	for _, subject := range subjects {
		subjectUIDs = append(subjectUIDs, subject.UID)
	}
	if err := b.ms.Unit.AddSubjects(ctx, tenant, unit, subjectUIDs); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// RemoveSubjects 管理单元批量移除请求主体
func (b *Unit) RemoveSubjects(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, subjects []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// AddPermissions 给管理单元添加权限，权限必须预先存在
func (b *Unit) AddPermissions(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, permissions []tpl.PermissionEx) (
	*tpl.SuccessResponseType, error) {
	if err := b.ms.Unit.AddPermissions(ctx, tenant, unit, permissions); err != nil {
		return nil, err
	}
	return &tpl.SuccessResponseType{Result: true}, nil
}

// UpdatePermissions 覆盖管理单元的权限，权限必须预先存在，当 permissions 为空时会清空权限
func (b *Unit) UpdatePermissions(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, permissions []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// RemovePermissions 移除管理单元的权限
func (b *Unit) RemovePermissions(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, permissions []string) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListChildren 列出管理单元的指定目标类型的子级管理单元，不包含 status 为 -1 的节点
func (b *Unit) ListChildren(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, targetType string,
	pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListDescendant 列出管理单元的指定目标类型的所有子孙管理单元，不包含 status 为 -1 的管理单元
// depth 定义对 targetType 类型管理单元的递归查询深度，而不是指定 unit 到 targetType 类型管理单元的深度，默认对 targetType 类型管理单元查到底
func (b *Unit) ListDescendant(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, targetType string, depth int,
	pg tpl.Pagination) (*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListPermissions 列出管理单元的直属权限
func (b *Unit) ListPermissions(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, pg tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListSubjects 列出管理单元的直属请求主体，不包含 status 为 -1 的请求主体
func (b *Unit) ListSubjects(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, pg tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// ListDescendantSubjects 列出管理单元及子孙管理单元下所有的请求主体，不包含 status 为 -1 的请求主体
func (b *Unit) ListDescendantSubjects(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, pg tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}

// GetDAG 根据 start 和 ends 找出一个 DAG，其中 start 可以为 Subject 或 Unit，ends 为 0 到多个 Unit，不包含 status 为 -1 的节点
func (b *Unit) GetDAG(ctx context.Context, tenant tpl.Tenant, start tpl.Target, ends []tpl.Target, pg tpl.Pagination) (
	*tpl.SuccessResponseType, error) {
	return nil, nil
}
