package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/open-trust/ot-ac/src/service/dgraph"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"
)

// Permission ...
type Permission struct {
	*Model
}

// BatchAdd ...
func (m *Permission) BatchAdd(ctx context.Context, tenant tpl.Tenant, permissions []string) error {
	nqs := make([]*dgraph.Nquads, 0, len(permissions))

	for _, p := range permissions {
		nqs = append(nqs, &dgraph.Nquads{
			UKkey: "OTAC.P.UK",
			UKval: util.HashBase64(tenant.Tenant, p),
			Type:  "OTACPermission",
			KV: map[string]interface{}{
				"OTAC.P-T": util.FormatUID(tenant.UID),
				"OTAC.P":   p,
			},
		})
	}

	_, err := m.Model.BatchAdd(ctx, nqs)
	return err
}

// List ...
func (m *Permission) List(ctx context.Context, tenant tpl.Tenant, resources []string, pageSize, skip int, uidToken string) (
	[]*tpl.Permission, error) {
	q := fmt.Sprintf(`query {
		result(func: eq(dgraph.type, "OTACPermission"), first: %d, offset: %d, after: %s) @filter(uid_in(OTAC.P-T, %s)) {
			uid
			permission: OTAC.P
		}
	}`, pageSize, skip, util.FormatUID(uidToken), util.FormatUID(tenant.UID))
	if len(resources) > 0 {
		q = fmt.Sprintf(`query {
			result(func: eq(dgraph.type, "OTACPermission"), first: %d, offset: %d, after: %s) @filter(uid_in(OTAC.P-T, %s) AND regexp(OTAC.P, /^(%s)/)) {
				uid
				permission: OTAC.P
			}
		}`, pageSize, skip, util.FormatUID(uidToken), util.FormatUID(tenant.UID), strings.Join(resources, "|"))
	}
	res := make([]*tpl.Permission, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Delete ...
func (m *Permission) Delete(ctx context.Context, tenant tpl.Tenant, permission string) error {
	q := fmt.Sprintf(`query {
		permissionUid as var(func: eq(OTAC.P, %s)) @filter(uid_in(OTAC.P-T, %s))
		objectUids as var(func: has(OTAC.O-Ps)) @filter(uid_in(OTAC.O-Ps, uid(permissionUid)))
		unitsUids as var(func: has(OTAC.U-Ps)) @filter(uid_in(OTAC.U-Ps, uid(permissionUid)))
	}`, util.FormatStr(permission), util.FormatUID(tenant.UID))
	delPermission := &dgraph.Nquads{
		ID: "uid(permissionUid)",
		KV: map[string]interface{}{
			"*": "*",
		},
	}
	delPermissionData, err := delPermission.Bytes()
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

	return m.Do(ctx, q, nil, nil, &api.Mutation{
		Cond:      "@if(gt(len(permissionUid), 0))",
		DelNquads: delPermissionData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(objectUids), 0))",
		DelNquads: ddelObjectsData,
	}, &api.Mutation{
		Cond:      "@if(gt(len(unitsUids), 0))",
		DelNquads: delUnitsData,
	})
}
