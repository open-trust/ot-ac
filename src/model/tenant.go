package model

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/open-trust/ot-ac/src/service/dgraph"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"
	otgo "github.com/open-trust/ot-go-lib"
)

// Tenant ...
type Tenant struct {
	*Model
}

// Get ...
func (m *Tenant) Get(ctx context.Context, tenant otgo.OTID) (*tpl.Tenant, error) {
	res := tpl.Tenant{Tenant: tenant.String(), Status: -1}
	q := fmt.Sprintf(`query {
		result(func: eq(OTAC.T, %s), first: 1) {
			uid
			status: OTAC.status
		}
	}`, util.FormatStr(res.Tenant))

	if err := m.Model.Get(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// List ...
func (m *Tenant) List(ctx context.Context, pageSize, skip int, uidToken string) ([]tpl.Tenant, error) {
	q := fmt.Sprintf(`query {
		result(func: eq(dgraph.type, "OTACTenant"), first: %d, offset: %d, after: %s) {
			uid
			tenant: OTAC.T
			status: OTAC.status
		}
	}`, pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.Tenant, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Add ...
func (m *Tenant) Add(ctx context.Context, input tpl.Tenant) (bool, error) {
	nq := &dgraph.Nquads{
		UKkey: "OTAC.T",
		UKval: input.Tenant,
		Type:  "OTACTenant",
		KV: map[string]interface{}{
			"OTAC.status": input.Status,
		},
	}

	return m.Model.Add(ctx, nq)
}

// Update ...
func (m *Tenant) Update(ctx context.Context, input tpl.Tenant) error {
	update := &dgraph.Nquads{
		UKkey: "OTAC.T",
		UKval: input.Tenant,
		Type:  "OTACTenant",
		KV: map[string]interface{}{
			"OTAC.status": input.Status,
		},
	}

	return m.Model.Update(ctx, update, "")
}

// Delete ...
func (m *Tenant) Delete(ctx context.Context, tenant otgo.OTID) error {
	q := fmt.Sprintf(`query {
		tenantUid as var(func: eq(OTAC.T, %s), first: 1) @filter(lt(OTAC.status, 0))
		objectUids as var(func: has(OTAC.O-T)) @filter(uid_in(OTAC.O-T, uid(tenantUid)))
		unitsUids as var(func: has(OTAC.U-T)) @filter(uid_in(OTAC.U-T, uid(tenantUid)))
		permissionsUids as var(func: has(OTAC.P-T)) @filter(uid_in(OTAC.P-T, uid(tenantUid)))
		scopesUids as var(func: has(OTAC.Sc-T)) @filter(uid_in(OTAC.Sc-T, uid(tenantUid)))
	}`, util.FormatStr(tenant.String()))
	delTenant := &dgraph.Nquads{
		ID: "uid(tenantUid)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	delTenantData, err := delTenant.Bytes()
	if err != nil {
		return err
	}
	delObjects := &dgraph.Nquads{
		ID: "uid(objectUids)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	ddelObjectsData, err := delObjects.Bytes()
	if err != nil {
		return err
	}
	delUnits := &dgraph.Nquads{
		ID: "uid(unitsUids)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	delUnitsData, err := delUnits.Bytes()
	if err != nil {
		return err
	}
	delPermissions := &dgraph.Nquads{
		ID: "uid(permissionsUids)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	delPermissionsData, err := delPermissions.Bytes()
	if err != nil {
		return err
	}
	delScopes := &dgraph.Nquads{
		ID: "uid(scopesUids)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	delScopesData, err := delScopes.Bytes()
	if err != nil {
		return err
	}

	return m.Do(ctx, q, nil, nil, &api.Mutation{
		Cond:      "@if(gt(len(tenantUid), 0))",
		DelNquads: delTenantData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(objectUids), 0))",
		DelNquads: ddelObjectsData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(unitsUids), 0))",
		DelNquads: delUnitsData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(permissionsUids), 0))",
		DelNquads: delPermissionsData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(scopesUids), 0))",
		DelNquads: delScopesData,
	})
}
