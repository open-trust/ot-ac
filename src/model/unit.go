package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/open-trust/ot-ac/src/service/dgraph"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"
	"github.com/teambition/gear"
)

// Unit ...
type Unit struct {
	*Model
}

// BatchAdd ...
func (m *Unit) BatchAdd(ctx context.Context, tenant tpl.Tenant, units []tpl.Target, parent *tpl.Target, scope *tpl.Target) error {
	nqs := make([]*dgraph.Nquads, 0, len(units)*2)
	parentUID, _, scopeUID, err := m.acquireUnitObjectScope(ctx, tenant, parent, nil, scope, 0)
	if err != nil {
		return err
	}
	uks := make([]string, 0, len(units))

	for _, unit := range units {
		create := &dgraph.Nquads{
			UKkey: "OTAC.U.UK",
			UKval: util.HashBase64(tenant.Tenant, unit.Type, unit.ID),
			Type:  "OTACUnit",
			KV: map[string]interface{}{
				"OTAC.U-T":    util.FormatUID(tenant.UID),
				"OTAC.status": 0,
				"OTAC.UId":    unit.ID,
				"OTAC.UType":  unit.Type,
			},
		}

		update := &dgraph.Nquads{
			KV: map[string]interface{}{},
		}
		if parentUID != "" {
			create.KV["OTAC.U-Us"] = util.FormatUID(parentUID)
			update.KV["OTAC.U-Us"] = util.FormatUID(parentUID)
		}
		if scopeUID != "" {
			create.KV["OTAC.U-Scs"] = util.FormatUID(scopeUID)
			update.KV["OTAC.U-Scs"] = util.FormatUID(scopeUID)
		}

		nqs = append(nqs, create)
		uks = append(uks, util.FormatStr(create.UKval))
		if parentUID != "" || scopeUID != "" {
			nqs = append(nqs, update)
		}
	}

	if parentUID != "" {
		checkCyclic := fmt.Sprintf(`
			var(func: uid(%s), first: 1) @recurse(loop: false) {
				uids as uid
				OTAC.U-Us
			}
			CyclicData(func: uid(uids), first: 1) @filter(eq(OTAC.U.UK, [%s])) {
				CyclicUID as uid
				targetType: OTAC.UType
				targetId: OTAC.UId
			}
		`, util.FormatUID(parentUID), strings.Join(uks, ", "))
		return m.Model.BatchAddOrUpdate(ctx, nqs, checkCyclic)
	}
	if scopeUID != "" {
		return m.Model.BatchAddOrUpdate(ctx, nqs, "")
	}

	_, err = m.Model.BatchAdd(ctx, nqs)
	return err
}

// AddFromOrg 从组织服务的 Org 创建管理单元，当检测到将形成环时会返回 400 错误
func (m *Unit) AddFromOrg(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, org string, parent *tpl.Target, scope *tpl.Target) error {
	orgUID, _, err := m.acquireOrgOU(ctx, org, "", 0)
	if err != nil {
		return err
	}
	if orgUID == "" {
		return gear.ErrBadRequest.WithMsgf("Organization(%s) not found", org)
	}
	parentUID, _, scopeUID, err := m.acquireUnitObjectScope(ctx, tenant, parent, nil, scope, 0)
	if err != nil {
		return err
	}
	nqs := make([]*dgraph.Nquads, 0, 2)
	create := &dgraph.Nquads{
		UKkey: "OTAC.U.UK",
		UKval: util.HashBase64(tenant.Tenant, unit.Type, unit.ID),
		Type:  "OTACUnit",
		KV: map[string]interface{}{
			"OTAC.U-T":    util.FormatUID(tenant.UID),
			"OTAC.status": 0,
			"OTAC.UId":    unit.ID,
			"OTAC.UType":  unit.Type,
		},
	}

	update := &dgraph.Nquads{
		KV: map[string]interface{}{},
	}
	create.KV["OTAC.U-Orgs"] = util.FormatUID(orgUID)
	update.KV["OTAC.U-Orgs"] = util.FormatUID(orgUID)
	if parentUID != "" {
		create.KV["OTAC.U-Us"] = util.FormatUID(parentUID)
		update.KV["OTAC.U-Us"] = util.FormatUID(parentUID)
	}
	if scopeUID != "" {
		create.KV["OTAC.U-Scs"] = util.FormatUID(scopeUID)
		update.KV["OTAC.U-Scs"] = util.FormatUID(scopeUID)
	}
	nqs = append(nqs, create, update)
	if parentUID != "" {
		checkCyclic := fmt.Sprintf(`
			var(func: uid(%s), first: 1) @recurse(loop: false) {
				uids as uid
				OTAC.U-Us
			}
			CyclicData(func: uid(uids), first: 1) @filter(eq(OTAC.U.UK, %s)) {
				CyclicUID as uid
				targetType: OTAC.UType
				targetId: OTAC.UId
			}
		`, util.FormatUID(parentUID), util.FormatStr(create.UKval))
		return m.Model.BatchAddOrUpdate(ctx, nqs, checkCyclic)
	}
	return m.Model.BatchAddOrUpdate(ctx, nqs, "")
}

// AddFromOU 从组织服务的 OU 创建管理单元，当检测到将形成环时会返回 400 错误
func (m *Unit) AddFromOU(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, org, ou string, parent *tpl.Target, scope *tpl.Target) error {
	_, ouUID, err := m.acquireOrgOU(ctx, org, ou, 0)
	if err != nil {
		return err
	}
	if ouUID == "" {
		return gear.ErrBadRequest.WithMsgf("OU(%s, %s) not found", org, ou)
	}
	parentUID, _, scopeUID, err := m.acquireUnitObjectScope(ctx, tenant, parent, nil, scope, 0)
	if err != nil {
		return err
	}
	nqs := make([]*dgraph.Nquads, 0, 2)
	create := &dgraph.Nquads{
		UKkey: "OTAC.U.UK",
		UKval: util.HashBase64(tenant.Tenant, unit.Type, unit.ID),
		Type:  "OTACUnit",
		KV: map[string]interface{}{
			"OTAC.U-T":    util.FormatUID(tenant.UID),
			"OTAC.status": 0,
			"OTAC.UId":    unit.ID,
			"OTAC.UType":  unit.Type,
		},
	}

	update := &dgraph.Nquads{
		KV: map[string]interface{}{},
	}
	create.KV["OTAC.U-OUs"] = util.FormatUID(ouUID)
	update.KV["OTAC.U-OUs"] = util.FormatUID(ouUID)
	if parentUID != "" {
		create.KV["OTAC.U-Us"] = util.FormatUID(parentUID)
		update.KV["OTAC.U-Us"] = util.FormatUID(parentUID)
	}
	if scopeUID != "" {
		create.KV["OTAC.U-Scs"] = util.FormatUID(scopeUID)
		update.KV["OTAC.U-Scs"] = util.FormatUID(scopeUID)
	}
	nqs = append(nqs, create, update)
	if parentUID != "" {
		checkCyclic := fmt.Sprintf(`
			var(func: uid(%s), first: 1) @recurse(loop: false) {
				uids as uid
				OTAC.U-Us
			}
			CyclicData(func: uid(uids), first: 1) @filter(eq(OTAC.U.UK, %s)) {
				CyclicUID as uid
				targetType: OTAC.UType
				targetId: OTAC.UId
			}
		`, util.FormatUID(parentUID), util.FormatStr(create.UKval))
		return m.Model.BatchAddOrUpdate(ctx, nqs, checkCyclic)
	}
	return m.Model.BatchAddOrUpdate(ctx, nqs, "")
}

// AddFromMembers 从组织服务的 Members 创建管理单元，当检测到将形成环时会返回 400 错误
func (m *Unit) AddFromMembers(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, org string, subjects []string, parent *tpl.Target, scope *tpl.Target) error {
	memberUIDs, err := m.acquireOrgMembers(ctx, org, subjects, 0)
	if err != nil {
		return err
	}
	if len(memberUIDs) != len(subjects) || len(subjects) == 0 {
		return gear.ErrBadRequest.WithMsgf("%d members not found for %#v", len(subjects)-len(memberUIDs), subjects)
	}
	parentUID, _, scopeUID, err := m.acquireUnitObjectScope(ctx, tenant, parent, nil, scope, 0)
	if err != nil {
		return err
	}
	nqs := make([]*dgraph.Nquads, 0, 2)
	create := &dgraph.Nquads{
		UKkey: "OTAC.U.UK",
		UKval: util.HashBase64(tenant.Tenant, unit.Type, unit.ID),
		Type:  "OTACUnit",
		KV: map[string]interface{}{
			"OTAC.U-T":    util.FormatUID(tenant.UID),
			"OTAC.status": 0,
			"OTAC.UId":    unit.ID,
			"OTAC.UType":  unit.Type,
		},
	}

	update := &dgraph.Nquads{
		KV: map[string]interface{}{},
	}
	create.KV["OTAC.U-Ms"] = util.FormatUIDs(memberUIDs)
	update.KV["OTAC.U-Ms"] = util.FormatUIDs(memberUIDs)
	if parentUID != "" {
		create.KV["OTAC.U-Us"] = util.FormatUID(parentUID)
		update.KV["OTAC.U-Us"] = util.FormatUID(parentUID)
	}
	if scopeUID != "" {
		create.KV["OTAC.U-Scs"] = util.FormatUID(scopeUID)
		update.KV["OTAC.U-Scs"] = util.FormatUID(scopeUID)
	}
	nqs = append(nqs, create, update)
	if parentUID != "" {
		checkCyclic := fmt.Sprintf(`
			var(func: uid(%s), first: 1) @recurse(loop: false) {
				uids as uid
				OTAC.U-Us
			}
			CyclicData(func: uid(uids), first: 1) @filter(eq(OTAC.U.UK, %s)) {
				CyclicUID as uid
				targetType: OTAC.UType
				targetId: OTAC.UId
			}
		`, util.FormatUID(parentUID), util.FormatStr(create.UKval))
		return m.Model.BatchAddOrUpdate(ctx, nqs, checkCyclic)
	}
	return m.Model.BatchAddOrUpdate(ctx, nqs, "")
}

// AddSubjects ...
func (m *Unit) AddSubjects(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, subjectUIDs []string) error {
	unitUID, _, _, err := m.acquireUnitObjectScope(ctx, tenant, &unit, nil, nil, 0)
	if err != nil {
		return err
	}
	if len(subjectUIDs) == 0 {
		return nil
	}

	nq := &dgraph.Nquads{
		ID: util.FormatUID(unitUID),
		KV: map[string]interface{}{
			"OTAC.U-Ss": util.FormatUIDs(subjectUIDs),
		},
	}
	data, err := nq.Bytes()
	if err != nil {
		return err
	}

	return m.Do(ctx, "", nil, nil, &api.Mutation{
		SetNquads: data,
	})
}

// AddPermissions ...
func (m *Unit) AddPermissions(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, permissions []tpl.PermissionEx) error {
	unitUID, _, _, err := m.acquireUnitObjectScope(ctx, tenant, &unit, nil, nil, 0)
	if err != nil {
		return err
	}
	if len(permissions) == 0 {
		return nil
	}

	ss := make([]string, 0, len(permissions))
	for i, p := range permissions {
		ss[i] = p.Permission
	}
	ps, err := m.acquirePermissions(ctx, tenant, ss)
	if err != nil {
		return err
	}

	fs := make([]dgraph.WithFacets, len(permissions))
	for i, p := range permissions {
		uid := tpl.GetPermissionUID(ps, p.Permission)
		if uid == "" {
			return gear.ErrBadRequest.WithMsgf("permission %s not found", util.FormatStr(p.Permission))
		}
		fs[i] = dgraph.WithFacets{V: util.FormatUID(uid), KV: p.Extensions}
	}

	nq := &dgraph.Nquads{
		ID: util.FormatUID(unitUID),
		KV: map[string]interface{}{
			"OTAC.U-Ps": fs,
		},
	}
	data, err := nq.Bytes()
	if err != nil {
		return err
	}

	return m.Do(ctx, "", nil, nil, &api.Mutation{
		SetNquads: data,
	})
}

// AssignParent ...
func (m *Unit) AssignParent(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, parent tpl.Target) error {
	unitUID, _, _, err := m.acquireUnitObjectScope(ctx, tenant, &unit, nil, nil, 0)
	if err != nil {
		return err
	}
	parentUID, _, _, err := m.acquireUnitObjectScope(ctx, tenant, &parent, nil, nil, 0)
	if err != nil {
		return err
	}
	checkCyclic := fmt.Sprintf(`
		var(func: uid(%s), first: 1) @recurse(loop: false) {
			uids as uid
			OTAC.U-Us
		}
		CyclicData(func: uid(uids), first: 1) @filter(uid(%s)) {
			CyclicUID as uid
			targetType: OTAC.UType
			targetId: OTAC.UId
		}
	`, util.FormatUID(parentUID), util.FormatUID(unitUID))

	nq := &dgraph.Nquads{
		ID: util.FormatUID(unitUID),
		KV: map[string]interface{}{
			"OTAC.U-Us": util.FormatUID(parentUID),
		},
	}
	return m.Model.Update(ctx, nq, checkCyclic)
}

// AssignScope ...
func (m *Unit) AssignScope(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, scope tpl.Target) error {
	unitUID, _, scopeUID, err := m.acquireUnitObjectScope(ctx, tenant, &unit, nil, &scope, 0)
	if err != nil {
		return err
	}
	nq := &dgraph.Nquads{
		ID: util.FormatUID(unitUID),
		KV: map[string]interface{}{
			"OTAC.U-Scs": util.FormatUID(scopeUID),
		},
	}
	return m.Model.Update(ctx, nq, "")
}

// AssignObject 建立管理单元与资源对象的关系
func (m *Unit) AssignObject(ctx context.Context, tenant tpl.Tenant, unit tpl.Target, object tpl.Target) error {
	unitUID, objectUID, _, err := m.acquireUnitObjectScope(ctx, tenant, &unit, &object, nil, 0)
	if err != nil {
		return err
	}
	nq := &dgraph.Nquads{
		ID: util.FormatUID(objectUID),
		KV: map[string]interface{}{
			"OTAC.O-Us": util.FormatUID(unitUID),
		},
	}
	return m.Model.Update(ctx, nq, "")
}
