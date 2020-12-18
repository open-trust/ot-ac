package model

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/open-trust/ot-ac/src/service/dgraph"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"
)

// Scope ...
type Scope struct {
	*Model
}

// Add 创建范围约束
func (m *Scope) Add(ctx context.Context, tenant tpl.Tenant, input tpl.Scope) (bool, error) {
	nq := &dgraph.Nquads{
		UKkey: "OTAC.Sc.UK",
		UKval: util.HashBase64(tenant.Tenant, input.TargetType, input.TargetID),
		Type:  "OTACScope",
		KV: map[string]interface{}{
			"OTAC.Sc-T":   util.FormatUID(tenant.UID),
			"OTAC.status": input.Status,
			"OTAC.ScId":   input.TargetID,
			"OTAC.ScType": input.TargetType,
		},
	}

	return m.Model.Add(ctx, nq)
}

// UpdateStatus 更新范围约束的状态，-1 表示停用
func (m *Scope) UpdateStatus(ctx context.Context, tenant tpl.Tenant, scope tpl.Target, status int) error {
	update := &dgraph.Nquads{
		UKkey: "OTAC.Sc.UK",
		UKval: util.HashBase64(tenant.Tenant, scope.Type, scope.ID),
		Type:  "OTACScope",
		KV: map[string]interface{}{
			"OTAC.status": status,
		},
	}
	return m.Model.Update(ctx, update, "")
}

// List 列出该系统当前所有指定目标类型的范围约束
func (m *Scope) List(ctx context.Context, tenant tpl.Tenant, targetType string,
	pageSize, skip int, uidToken string) ([]*tpl.Scope, error) {
	q := fmt.Sprintf(`query {
		result(func: eq(OTAC.ScType, %s), first: %d, offset: %d, after: %s) @filter(uid_in(OTAC.Sc-T, %s)) {
			uid
			status: OTAC.status
			targetId: OTAC.ScId
			targetType: OTAC.ScType
		}
	}`, util.FormatStr(targetType), pageSize, skip, util.FormatUID(uidToken), util.FormatUID(tenant.UID))
	res := make([]*tpl.Scope, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListUnits 列出范围约束下指定目标类型的直属的管理单元
func (m *Scope) ListUnits(ctx context.Context, tenant tpl.Tenant, scope tpl.Target, targetType string,
	pageSize, skip int, uidToken string) ([]*tpl.Unit, error) {
	q := fmt.Sprintf(`query {
		scopeUid as var(func: eq(OTAC.ScId, %s), first: 1) @filter(eq(OTAC.ScType, %s) AND uid_in(OTAC.Sc-T, %s))
		result(func: eq(OTAC.UType, %s), first: %d, offset: %d, after: %s) @filter(uid_in(OTAC.U-Scs, uid(scopeUid))) {
			uid
			status: OTAC.status
			targetId: OTAC.UId
			targetType: OTAC.UType
		}
	}`, util.FormatStr(scope.ID), util.FormatStr(scope.Type), util.FormatUID(tenant.UID), util.FormatStr(targetType), pageSize, skip, util.FormatUID(uidToken))
	res := make([]*tpl.Unit, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListObjects 列出范围约束下指定目标类型的直属的资源对象
func (m *Scope) ListObjects(ctx context.Context, tenant tpl.Tenant, scope tpl.Target, targetType string,
	pageSize, skip int, uidToken string) ([]*tpl.Object, error) {
	q := fmt.Sprintf(`query {
			scopeUid as var(func: eq(OTAC.ScId, %s), first: 1) @filter(eq(OTAC.ScType, %s) AND uid_in(OTAC.Sc-T, %s))
			result(func: eq(OTAC.OType, %s), first: %d, offset: %d, after: %s) @filter(uid_in(OTAC.O-Scs, uid(scopeUid))) {
				uid
				status: OTAC.status
				targetId: OTAC.OId
				targetType: OTAC.OType
			}
		}`, util.FormatStr(scope.ID), util.FormatStr(scope.Type), util.FormatUID(tenant.UID), util.FormatStr(targetType), pageSize, skip, util.FormatUID(uidToken))
	res := make([]*tpl.Object, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Delete 删除范围约束
func (m *Scope) Delete(ctx context.Context, tenant tpl.Tenant, scope tpl.Target) error {
	q := fmt.Sprintf(`query {
		scopeUid as var(func: eq(OTAC.ScId, %s), first: 1) @filter(eq(OTAC.ScType, %s) AND uid_in(OTAC.Sc-T, %s))
		objectUids as var(func: has(OTAC.O-Scs)) @filter(uid_in(OTAC.O-Scs, uid(scopeUid)))
		unitsUids as var(func: has(OTAC.U-Scs)) @filter(uid_in(OTAC.U-Scs, uid(scopeUid)))
	}`, util.FormatStr(scope.ID), util.FormatStr(scope.Type), util.FormatUID(tenant.UID))
	delScope := &dgraph.Nquads{
		ID: "uid(scopeUid)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	delScopeData, err := delScope.Bytes()
	if err != nil {
		return err
	}
	delObjects := &dgraph.Nquads{
		ID: "uid(objectUids)",
		KV: map[string]interface{}{
			"OTAC.O-Scs": "uid(scopeUid)",
		},
	}
	ddelObjectsData, err := delObjects.Bytes()
	if err != nil {
		return err
	}
	delUnits := &dgraph.Nquads{
		ID: "uid(unitsUids)",
		KV: map[string]interface{}{
			"OTAC.U-Scs": "uid(scopeUid)",
		},
	}
	delUnitsData, err := delUnits.Bytes()
	if err != nil {
		return err
	}

	return m.Do(ctx, q, nil, nil, &api.Mutation{
		Cond:      "@if(gt(len(scopeUid), 0))",
		DelNquads: delScopeData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(objectUids), 0))",
		DelNquads: ddelObjectsData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(unitsUids), 0))",
		DelNquads: delUnitsData,
	})
}

// DeleteAll 删除范围约束及范围内的所有 Unit 和 Object
func (m *Scope) DeleteAll(ctx context.Context, tenant tpl.Tenant, scope tpl.Target) error {
	q := fmt.Sprintf(`query {
		scopeUid as var(func: eq(OTAC.ScId, %s), first: 1) @filter(eq(OTAC.ScType, %s) AND uid_in(OTAC.Sc-T, %s))
		objectUids as var(func: has(OTAC.O-Scs)) @filter(uid_in(OTAC.O-Scs, uid(scopeUid)))
		unitsUids as var(func: has(OTAC.U-Scs)) @filter(uid_in(OTAC.U-Scs, uid(scopeUid)))
		descendantObjectUids as var(func: uid(objectUids)) @recurse(loop: false) {
			~OTAC.O-Os
		}
		descendantUnitUids as var(func: uid(unitsUids)) @recurse(loop: false) {
			~OTAC.U-Us
		}
	}`, util.FormatStr(scope.ID), util.FormatStr(scope.Type), util.FormatUID(tenant.UID))
	delScope := &dgraph.Nquads{
		ID: "uid(scopeUid)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	delScopeData, err := delScope.Bytes()
	if err != nil {
		return err
	}
	delObjects := &dgraph.Nquads{
		ID: "uid(descendantObjectUids)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	ddelObjectsData, err := delObjects.Bytes()
	if err != nil {
		return err
	}
	delUnits := &dgraph.Nquads{
		ID: "uid(descendantUnitUids)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	delUnitsData, err := delUnits.Bytes()
	if err != nil {
		return err
	}

	return m.Do(ctx, q, nil, nil, &api.Mutation{
		Cond:      "@if(gt(len(scopeUid), 0))",
		DelNquads: delScopeData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(descendantObjectUids), 0))",
		DelNquads: ddelObjectsData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(descendantUnitUids), 0))",
		DelNquads: delUnitsData,
	})
}
