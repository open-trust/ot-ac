package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-trust/ot-ac/src/service/dgraph"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"
	"github.com/teambition/gear"
)

// Object ...
type Object struct {
	*Model
}

// BatchAdd ...
func (m *Object) BatchAdd(ctx context.Context, tenant tpl.Tenant, objects []tpl.Target, parent *tpl.Target, scope *tpl.Target) error {
	nqs := make([]*dgraph.Nquads, 0, len(objects)*2)
	_, parentUID, scopeUID, err := m.acquireUnitObjectScope(ctx, tenant, nil, parent, scope, 0)
	if err != nil {
		return err
	}
	uks := make([]string, 0, len(objects))

	for _, obj := range objects {
		create := &dgraph.Nquads{
			UKkey: "OTAC.O.UK",
			UKval: util.HashBase64(tenant.Tenant, obj.Type, obj.ID),
			Type:  "OTACObject",
			KV: map[string]interface{}{
				"OTAC.O-T":   util.FormatUID(tenant.UID),
				"OTAC.OId":   obj.ID,
				"OTAC.OType": obj.Type,
			},
		}

		update := &dgraph.Nquads{
			KV: map[string]interface{}{},
		}
		if parentUID != "" {
			create.KV["OTAC.O-Os"] = util.FormatUID(parentUID)
			update.KV["OTAC.O-Os"] = util.FormatUID(parentUID)
		}
		if scopeUID != "" {
			create.KV["OTAC.O-Scs"] = util.FormatUID(scopeUID)
			update.KV["OTAC.O-Scs"] = util.FormatUID(scopeUID)
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
				OTAC.O-Os
			}
			CyclicData(func: uid(uids), first: 1) @filter(eq(OTAC.O.UK, [%s])) {
				CyclicUID as uid
				targetType: OTAC.OType
				targetId: OTAC.OId
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

// AddPermissions ...
func (m *Object) AddPermissions(ctx context.Context, tenant tpl.Tenant, object tpl.Target, permissions []string) error {
	_, objectUID, _, err := m.acquireUnitObjectScope(ctx, tenant, nil, &object, nil, 0)
	if err != nil {
		return err
	}
	if len(permissions) == 0 {
		return nil
	}

	ps, err := m.acquirePermissions(ctx, tenant, permissions)
	if err != nil {
		return err
	}

	uids := make([]string, 0, len(permissions))
	for _, p := range permissions {
		uid := tpl.GetPermissionUID(ps, p)
		if uid == "" {
			return gear.ErrBadRequest.WithMsgf("permission %s not found", util.FormatStr(p))
		}
		uids = append(uids, uid)
	}

	nq := &dgraph.Nquads{
		ID: util.FormatUID(objectUID),
		KV: map[string]interface{}{
			"OTAC.O-Ps": util.FormatUIDs(uids),
		},
	}
	return m.Model.Update(ctx, nq, "")
}
