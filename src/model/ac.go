package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"

	daggo "github.com/open-trust/dag-go"
)

// AC ...
type AC struct {
	*Model
}

// V ...
type V struct {
	UID         string                    `json:"uid"`
	Typ         string                    `json:"type"`
	Permissions []tpl.ACPermissionPayload `json:"permissions"`
}

// ID ...
func (v *V) ID() string {
	return v.UID
}

// Type ...
func (v *V) Type() string {
	return v.Typ
}

func getIDsFromDAG(dag *daggo.DAG, typ string) []string {
	vs := dag.Vertices(typ)
	vs.Sort()
	return vs.IDs()
}

type jsonCheckUnitOutput struct {
	UID     string                `json:"uid"`
	Parents []jsonCheckUnitOutput `json:"OTAC.U-Us"`
}

type jsonRawPermissionsOutput struct {
	UID         string                   `json:"uid"`
	ID          string                   `json:"targetId"`
	Type        string                   `json:"targetType"`
	Permissions []map[string]interface{} `json:"permissions"`
}

func rawToPermissions(input jsonRawPermissionsOutput) []tpl.ACPermissionPayload {
	data := make([]tpl.ACPermissionPayload, 0, len(input.Permissions))
	for _, raw := range input.Permissions {
		p := tpl.ACPermissionPayload{
			Target: tpl.Target{Type: input.Type, ID: input.ID},
		}
		p.Permission = raw["permission"].(string)
		p.Extensions = make(map[string]interface{})
		for k, v := range raw {
			if strings.HasPrefix(k, "permissions|") {
				p.Extensions[k[12:]] = v
			}
		}
		data = append(data, p)
	}
	return data
}

func rawsToPermissions(input []jsonRawPermissionsOutput) []tpl.ACPermissionPayload {
	data := make([]tpl.ACPermissionPayload, 0, len(input))
	for _, target := range input {
		data = append(data, rawToPermissions(target)...)
	}
	return data
}

// CheckUnit 检查请求主体到指定管理单元有没有指定权限
func (m *AC) CheckUnit(ctx context.Context, tenant tpl.Tenant, subject string,
	unit tpl.Target, permissions []string, withOrganization bool) (interface{}, error) {
	unitUID, _, _, err := m.acquireUnitObjectScope(ctx, tenant, &unit, nil, nil, 0)
	if err != nil {
		return nil, err
	}

	dag, err := m.getUnitsDAG(ctx, subject, tenant.UID, withOrganization)
	if err != nil {
		return nil, err
	}

	dag = dag.CloseDAG(&V{UID: subject, Typ: "Subject"}, &V{UID: unitUID, Typ: "Unit"})
	unitUIDs := getIDsFromDAG(dag, "Unit")
	if respondDetail(ctx) {
		return m.checkUnitPermissionsWithDetail(ctx, tenant.UID, unitUIDs, permissions)
	}
	return m.checkUnitPermissions(ctx, tenant.UID, unitUIDs, permissions)
}

func (m *AC) getUnitsDAG(ctx context.Context, subject, tenantUID string, withOrganization bool) (*daggo.DAG, error) {
	sub := util.FormatStr(subject)
	tt := util.FormatUID(tenantUID)
	q := ""
	switch {
	case withOrganization:
		q = fmt.Sprintf(`query {
			var(func: eq(OTAC.Sub, %s), first: 1) @filter(ge(OTAC.status, 0)) {
				~OTAC.U-Ss @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0))  {
					unitUIDs0 as uid
				}
				~OTAC.M-S @filter(ge(OTAC.status, 0)) {
					ouUIDs as ~OTAC.OU-Ms @filter(ge(OTAC.status, 0))
					~OTAC.U-Ms @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0))  {
						unitUIDs1 as uid
					}
				}
			}
			var(func: uid(ouUIDs)) @recurse(loop: false) {
				ouUIDs1 as uid
				OTAC.OU-OU @filter(ge(OTAC.status, 0))
			}
			var(func: uid(ouUIDs1)) {
				~OTAC.U-OUs @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0))  {
					unitUIDs2 as uid
				}
				OTAC.OU-Org @filter(ge(OTAC.status, 0)) {
					~OTAC.U-Orgs @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0))  {
						unitUIDs3 as uid
					}
				}
			}
			result(func: uid(unitUIDs0, unitUIDs1, unitUIDs2, unitUIDs3)) @recurse(loop: false) {
				uid
				OTAC.U-Us @filter(ge(OTAC.status, 0))
			}
		}`, sub, tt, tt, tt, tt)
	default:
		q = fmt.Sprintf(`query {
			var(func: eq(OTAC.Sub, %s), first: 1) @filter(ge(OTAC.status, 0)) {
				~OTAC.U-Ss @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0))  {
					unitUIDs as uid
				}
			}
			result(func: uid(unitUIDs)) @recurse(loop: false) {
				uid
				OTAC.U-Us @filter(ge(OTAC.status, 0))
			}
		}`, sub, tt)
	}

	data := make([]jsonCheckUnitOutput, 0, 10)
	if err := m.Model.List(ctx, q, nil, &data); err != nil {
		return nil, err
	}
	dag := daggo.New()
	var iterator func(start daggo.Vertice, parents []jsonCheckUnitOutput) error
	iterator = func(start daggo.Vertice, parents []jsonCheckUnitOutput) error {
		for _, v := range parents {
			node := &V{UID: v.UID, Typ: "Unit"}
			err := dag.AddEdge(start, node, 0)
			if err != nil {
				return err
			}
			if len(v.Parents) > 0 {
				if err := iterator(node, v.Parents); err != nil {
					return err
				}
			}
		}
		return nil
	}
	s := &V{UID: subject, Typ: "Subject"}
	if err := iterator(s, data); err != nil {
		return nil, err
	}
	return dag, nil
}

func (m *AC) checkUnitPermissions(ctx context.Context, tenantUID string, unitUIDs, permissions []string) (bool, error) {
	if len(unitUIDs) == 0 {
		return false, nil
	}
	q := fmt.Sprintf(`query {
		var(func: uid(%s)) @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0)) {
			uids as OTAC.U-Ps @filter(uid_in(OTAC.P-T, %s) AND eq(OTAC.P, [%s]))
		}
		result(func: uid(uids), first: 1) { uid }
	}`, strings.Join(util.FormatUIDs(unitUIDs), ", "), util.FormatUID(tenantUID), util.FormatUID(tenantUID), strings.Join(util.FormatStrs(permissions), ", "))
	data := make([]jsonUID, 0)
	if err := m.Model.List(ctx, q, nil, &data); err != nil {
		return false, err
	}
	return (len(data) > 0), nil
}

func (m *AC) checkUnitPermissionsWithDetail(ctx context.Context, tenantUID string, unitUIDs, permissions []string) ([]tpl.ACPermissionPayload, error) {
	if len(unitUIDs) == 0 {
		return make([]tpl.ACPermissionPayload, 0), nil
	}
	q := fmt.Sprintf(`query {
		result(func: uid(%s)) @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0)) {
			targetType: OTAC.UType
			targetId: OTAC.UId
			permissions: OTAC.U-Ps @filter(uid_in(OTAC.P-T, %s) AND eq(OTAC.P, [%s])) @facets {
				permission: OTAC.P
			}
		}
	}`, strings.Join(util.FormatUIDs(unitUIDs), ", "), util.FormatUID(tenantUID), util.FormatUID(tenantUID), strings.Join(util.FormatStrs(permissions), ", "))
	data := make([]jsonRawPermissionsOutput, 0)
	if err := m.Model.List(ctx, q, nil, &data); err != nil {
		return nil, err
	}
	return rawsToPermissions(data), nil
}

// CheckScope 检查请求主体到指定范围约束有没有指定权限
func (m *AC) CheckScope(ctx context.Context, tenant tpl.Tenant, subject string,
	scope tpl.Target, permissions []string, withOrganization bool) (interface{}, error) {
	_, _, scopeUID, err := m.acquireUnitObjectScope(ctx, tenant, nil, nil, &scope, 0)
	if err != nil {
		return nil, err
	}

	dag, err := m.getUnitsDAG(ctx, subject, tenant.UID, withOrganization)
	if err != nil {
		return nil, err
	}

	unitUIDs := getIDsFromDAG(dag, "Unit")
	q := fmt.Sprintf(`query {
		result(func: uid(%s)) @filter(uid_in(OTAC.U-Scs, %s)) {
			uid
		}
	}`, strings.Join(util.FormatUIDs(unitUIDs), ", "), util.FormatUID(scopeUID))
	data := make([]jsonUID, 0)
	if err := m.Model.List(ctx, q, nil, &data); err != nil {
		return false, err
	}
	scopeNode := &V{UID: scopeUID, Typ: "Scope"}
	for _, v := range data {
		err := dag.AddEdge(&V{UID: v.UID, Typ: "Unit"}, scopeNode, 0)
		if err != nil {
			return nil, err
		}
	}

	dag = dag.CloseDAG(&V{UID: subject, Typ: "Subject"}, scopeNode)
	unitUIDs = getIDsFromDAG(dag, "Unit")

	if respondDetail(ctx) {
		return m.checkUnitPermissionsWithDetail(ctx, tenant.UID, unitUIDs, permissions)
	}
	return m.checkUnitPermissions(ctx, tenant.UID, unitUIDs, permissions)
}

// CheckObject 检查请求主体通过 Scope 或 Unit -> Object 的连接关系到指定资源对象有没有指定权限，如果 ignoreScope 为 true，则要求必须有 Unit -> Object 的连接关系
func (m *AC) CheckObject(ctx context.Context, tenant tpl.Tenant, subject string,
	object tpl.Target, permissions []string, withOrganization, ignoreScope bool) (interface{}, error) {
	_, objectUID, _, err := m.acquireUnitObjectScope(ctx, tenant, nil, &object, nil, 0)
	if err != nil {
		return nil, err
	}
	objectDAG, err := m.getObjectsDAG(ctx, objectUID)
	if err != nil {
		return nil, err
	}
	if objectDAG.Len() == 0 {
		if respondDetail(ctx) {
			return make([]tpl.ACPermissionPayload, 0), nil
		}
		return false, nil
	}

	q := `query {
		result(func: uid(%s)) @filter(uid_in(OTAC.O-T, %s)) {
			uid
			units: OTAC.O-Us @filter(ge(OTAC.status, 0)) {
				uid
			}
		}
	}`
	if !ignoreScope {
		q = `query {
			result(func: uid(%s)) @filter(uid_in(OTAC.O-T, %s)) {
				uid
				units: OTAC.O-Us @filter(ge(OTAC.status, 0)) {
					uid
				}
				scopes: OTAC.O-Scs @filter(ge(OTAC.status, 0)) {
					uid
				}
			}
		}`
	}
	q = fmt.Sprintf(q, strings.Join(util.FormatUIDs(getIDsFromDAG(objectDAG, "Object")), ", "), util.FormatUID(tenant.UID))
	data := make([]jsonCheckScopeOutput, 0)
	if err := m.Model.List(ctx, q, nil, &data); err != nil {
		return nil, err
	}

	for _, v := range data {
		start := &V{UID: v.UID, Typ: "Object"}
		for _, u := range v.Units {
			err := objectDAG.AddEdge(start, &V{UID: u.UID, Typ: "Unit"}, 0)
			if err != nil {
				return nil, err
			}
		}
		for _, s := range v.Scopes {
			err := objectDAG.AddEdge(start, &V{UID: s.UID, Typ: "Scope"}, 0)
			if err != nil {
				return nil, err
			}
		}
	}

	objectUnitUIDs := getIDsFromDAG(objectDAG, "Unit")
	scopeUIDs := getIDsFromDAG(objectDAG, "Scope")
	if len(objectUnitUIDs) == 0 && len(scopeUIDs) == 0 {
		if respondDetail(ctx) {
			return make([]tpl.ACPermissionPayload, 0), nil
		}
		return false, nil
	}

	unitDAG, err := m.getUnitsDAG(ctx, subject, tenant.UID, withOrganization)
	if err != nil {
		return nil, err
	}
	if unitDAG.Len() == 0 {
		if respondDetail(ctx) {
			return make([]tpl.ACPermissionPayload, 0), nil
		}
		return false, nil
	}

	if len(scopeUIDs) > 0 {
		q = fmt.Sprintf(`query {
			result(func: uid(%s)) @filter(uid_in(OTAC.U-Scs, [%s])) {
				uid
				scopes: OTAC.U-Scs @filter(ge(OTAC.status, 0)) {
					uid
				}
			}
		}`, strings.Join(util.FormatUIDs(getIDsFromDAG(unitDAG, "Unit")), ", "), strings.Join(util.FormatUIDs(scopeUIDs), ", "))
		data := make([]jsonCheckScopeOutput, 0)

		if err := m.Model.List(ctx, q, nil, &data); err != nil {
			return false, err
		}
		for _, v := range data {
			start := &V{UID: v.UID, Typ: "Unit"}
			for _, s := range v.Scopes {
				err := unitDAG.AddEdge(start, &V{UID: s.UID, Typ: "Scope"}, 0)
				if err != nil {
					return nil, err
				}
			}
		}
		scopeUIDs = getIDsFromDAG(unitDAG, "Scope")
	}

	if err = unitDAG.Merge(objectDAG.Reverse()); err != nil {
		return nil, err
	}

	dag := unitDAG.CloseDAG(&V{UID: subject, Typ: "Subject"}, &V{UID: objectUID, Typ: "Object"})
	if respondDetail(ctx) {
		return m.checkDAGPermissionsWithDetail(ctx, tenant.UID, dag, permissions)
	}
	return m.checkDAGPermissions(ctx, tenant.UID, dag, permissions)
}

type jsonDAGPermissions struct {
	Units   []jsonRawPermissionsOutput `json:"units"`
	Objects []jsonRawPermissionsOutput `json:"objects"`
}

func (m *AC) checkDAGPermissions(ctx context.Context, tenantUID string, dag *daggo.DAG, permissions []string) (bool, error) {
	if dag.Len() == 0 {
		return false, nil
	}
	unitUIDs := getIDsFromDAG(dag, "Unit")
	objectUIDs := getIDsFromDAG(dag, "Object")
	fTenantUID := util.FormatUID(tenantUID)
	fPermissions := strings.Join(util.FormatStrs(permissions), ", ")
	q := fmt.Sprintf(`query {
		units(func: uid(%s)) @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0)) {
			uid
			permissions: OTAC.U-Ps @filter(uid_in(OTAC.P-T, %s) AND eq(OTAC.P, [%s])) {
				permission: OTAC.P
			}
		}
		objects(func: uid(%s)) @filter(uid_in(OTAC.O-T, %s)) {
			uid
			permissions: OTAC.O-Ps @filter(uid_in(OTAC.P-T, %s) AND eq(OTAC.P, [%s])) {
				permission: OTAC.P
			}
		}
	}`, strings.Join(util.FormatUIDs(unitUIDs), ", "), fTenantUID, fTenantUID, fPermissions,
		strings.Join(util.FormatUIDs(objectUIDs), ", "), fTenantUID, fTenantUID, fPermissions)
	data := &jsonDAGPermissions{}
	if err := m.Model.QueryBestEffort(ctx, q, nil, &data); err != nil {
		return false, err
	}
	if len(data.Units) == 0 && len(data.Objects) == 0 {
		return false, nil
	}
	for _, unit := range data.Units {
		v := dag.GetVertice("Unit", unit.UID)
		v.(*V).Permissions = rawToPermissions(unit)
	}
	for _, obj := range data.Objects {
		v := dag.GetVertice("Object", obj.UID)
		v.(*V).Permissions = rawToPermissions(obj)
	}
	sub := dag.StartingVertices()[0]
	ps := dag.Iterate(sub, nil, func(v daggo.Vertice, _ int, acc []interface{}) []interface{} {
		val := v.(*V)
		switch val.Typ {
		case "Unit":
			for _, p := range val.Permissions {
				acc = append(acc, p)
			}
			return acc
		case "Object":
			return removeACPermissionPayload(acc, val.Permissions)
		}
		return acc
	})
	return (len(ps) > 0), nil
}

func removeACPermissionPayload(acc []interface{}, allow []tpl.ACPermissionPayload) []interface{} {
	if len(allow) == 0 {
		return acc
	}

	exists := map[string]struct{}{}
	for _, p := range allow {
		exists[p.Permission] = struct{}{}
	}
	res := make([]interface{}, 0)
	for _, v := range acc {
		p := v.(tpl.ACPermissionPayload)
		if _, ok := exists[p.Permission]; ok {
			res = append(res, p)
		}
	}
	return res
}

func (m *AC) checkDAGPermissionsWithDetail(ctx context.Context, tenantUID string, dag *daggo.DAG, permissions []string) ([]tpl.ACPermissionPayload, error) {
	res := make([]tpl.ACPermissionPayload, 0)
	if dag.Len() == 0 {
		return res, nil
	}
	unitUIDs := getIDsFromDAG(dag, "Unit")
	objectUIDs := getIDsFromDAG(dag, "Object")
	fTenantUID := util.FormatUID(tenantUID)
	fPermissions := strings.Join(util.FormatStrs(permissions), ", ")
	q := fmt.Sprintf(`query {
		units(func: uid(%s)) @filter(uid_in(OTAC.U-T, %s) AND ge(OTAC.status, 0)) {
			uid
			targetType: OTAC.UType
			targetId: OTAC.UId
			permissions: OTAC.U-Ps @filter(uid_in(OTAC.P-T, %s) AND eq(OTAC.P, [%s])) @facets {
				permission: OTAC.P
			}
		}
		objects(func: uid(%s)) @filter(uid_in(OTAC.O-T, %s)) {
			uid
			targetType: OTAC.OType
			targetId: OTAC.OId
			permissions: OTAC.O-Ps @filter(uid_in(OTAC.P-T, %s) AND eq(OTAC.P, [%s])) {
				permission: OTAC.P
			}
		}
	}`, strings.Join(util.FormatUIDs(unitUIDs), ", "), fTenantUID, fTenantUID, fPermissions,
		strings.Join(util.FormatUIDs(objectUIDs), ", "), fTenantUID, fTenantUID, fPermissions)
	data := &jsonDAGPermissions{}
	if err := m.Model.QueryBestEffort(ctx, q, nil, &data); err != nil {
		return nil, err
	}
	if len(data.Units) == 0 && len(data.Objects) == 0 {
		return res, nil
	}
	for _, unit := range data.Units {
		v := dag.GetVertice("Unit", unit.UID)
		v.(*V).Permissions = rawToPermissions(unit)
	}
	for _, obj := range data.Objects {
		v := dag.GetVertice("Object", obj.UID)
		v.(*V).Permissions = rawToPermissions(obj)
	}
	sub := dag.StartingVertices()[0]
	ps := dag.Iterate(sub, nil, func(v daggo.Vertice, _ int, acc []interface{}) []interface{} {
		val := v.(*V)
		switch val.Typ {
		case "Unit":
			for _, p := range val.Permissions {
				acc = append(acc, p)
			}
			return acc
		case "Object":
			return removeACPermissionPayload(acc, val.Permissions)
		}
		return acc
	})
	for _, p := range ps {
		res = append(res, p.(tpl.ACPermissionPayload))
	}

	return res, nil
}

type jsonCheckScopeOutput struct {
	UID    string    `json:"uid"`
	Units  []jsonUID `json:"units"`
	Scopes []jsonUID `json:"scopes"`
}

type jsonCheckObjectOutput struct {
	UID     string                  `json:"uid"`
	Parents []jsonCheckObjectOutput `json:"OTAC.O-Os"`
}

func (m *AC) getObjectsDAG(ctx context.Context, objectUID string) (*daggo.DAG, error) {
	q := fmt.Sprintf(`query {
		result(func: uid(%s), first: 1) @recurse(loop: false) {
			uid
			OTAC.O-Os
		}
	}`, util.FormatUID(objectUID))

	data := make([]jsonCheckObjectOutput, 0, 10)
	if err := m.Model.List(ctx, q, nil, &data); err != nil {
		return nil, err
	}
	dag := daggo.New()
	if len(data) == 0 || len(data[0].Parents) == 0 {
		return dag, nil
	}

	var iterator func(start daggo.Vertice, parents []jsonCheckObjectOutput) error
	iterator = func(start daggo.Vertice, parents []jsonCheckObjectOutput) error {
		for _, v := range parents {
			node := &V{UID: v.UID, Typ: "Object"}
			err := dag.AddEdge(start, node, 0)
			if err != nil {
				return err
			}
			if len(v.Parents) > 0 {
				if err := iterator(node, v.Parents); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := iterator(&V{UID: objectUID, Typ: "Object"}, data[0].Parents); err != nil {
		return nil, err
	}
	return dag, nil
}

// ListPermissionsByUnit 列出请求主体到指定管理单元的符合 resource 的权限，如果未指定管理单元，则会查询请求主体能触达的所有管理单元，如果 resources 为空，则会列出所有触达的有效权限
func (m *AC) ListPermissionsByUnit(ctx context.Context, tenant tpl.Tenant, subject string, unit tpl.Target, resources []string, withOrganization bool) ([]tpl.ACPermissionPayload, error) {
}

// ListPermissionsByScope 列出请求主体到指定范围约束的符合 resource 的权限，如果 resources 为空，则会列出所有触达的有效权限
func (m *AC) ListPermissionsByScope(ctx context.Context, tenant tpl.Tenant, subject string, scope tpl.Target, resources []string, withOrganization bool) ([]tpl.ACPermissionPayload, error) {
}

// ListPermissionsByObject 列出请求主体到指定资源对象的符合 resource 的权限，如果 resources 为空，则会列出所有触达的有效权限
func (m *AC) ListPermissionsByObject(ctx context.Context, tenant tpl.Tenant, subject string, object tpl.Target, resources []string, withOrganization, ignoreScope bool) ([]tpl.ACPermissionPayload, error) {
}

// ListUnits 列出请求主体参与的指定类型的管理单元
func (m *AC) ListUnits(ctx context.Context, tenant tpl.Tenant, subject string, targetType string, withOrganization bool) {
}

// ListObjects 列出请求主体在指定资源对象中能触达的所有指定类型的子孙资源对象
// depth 定义对 targetType 类型资源对象的递归查询深度，而不是指定 object 到 targetType 类型资源对象的深度，默认对 targetType 类型资源对象查到底
func (m *AC) ListObjects(ctx context.Context, tenant tpl.Tenant, subject string, object tpl.Target, permission, targetType string, withOrganization, ignoreScope bool, depth int) {
}

// SearchObjects 根据关键词，在指定资源对象的子孙资源对象中，对请求主体能触达的所有指定类型的资源对象中进行搜索，term 为空不匹配任何资源对象
func (m *AC) SearchObjects(ctx context.Context, tenant tpl.Tenant, subject string, object tpl.Target, permission, targetType, term string, withOrganization, ignoreScope bool) {
}
