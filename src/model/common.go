package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/open-trust/ot-ac/src/service/dgraph"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"
	otgo "github.com/open-trust/ot-go-lib"
	"github.com/teambition/gear"
)

func init() {
	util.DigProvide(NewModels)
}

// Model ...
type Model struct {
	*dgraph.Dgraph
}

// Models ...
type Models struct {
	Model        *Model
	AC           *AC
	Object       *Object
	Organization *Organization
	Permission   *Permission
	Scope        *Scope
	Unit         *Unit
	Tenant       *Tenant
	Subject      *Subject
}

// NewModels ...
func NewModels(dg *dgraph.Dgraph) *Models {
	m := &Model{dg}
	return &Models{
		Model:        m,
		AC:           &AC{m},
		Object:       &Object{m},
		Organization: &Organization{m},
		Permission:   &Permission{m},
		Scope:        &Scope{m},
		Unit:         &Unit{m},
		Tenant:       &Tenant{m},
		Subject:      &Subject{m},
	}
}

type contextKey int

const (
	idempotentKey contextKey = iota
	respondDetailKey
)

// ContextWithPrefer ...
func ContextWithPrefer(ctx *gear.Context) context.Context {
	c := ctx.Context()
	idempotent := true
	respondDetail := false
	for _, prefer := range ctx.GetHeaders("Prefer") {
		prefer = strings.ToLower(prefer)
		if prefer == "respond-conflict" {
			idempotent = false
		}
		if prefer == "respond-detail" {
			respondDetail = true
		}
		// other prefer ...
	}
	if !idempotent {
		c = context.WithValue(c, idempotentKey, struct{}{})
	}
	if respondDetail {
		c = context.WithValue(c, respondDetailKey, struct{}{})
	}

	return c
}

func isIdempotent(ctx context.Context) bool {
	return ctx.Value(idempotentKey) == nil
}

func respondDetail(ctx context.Context) bool {
	return ctx.Value(respondDetailKey) != nil
}

type jsonUID struct {
	UID string `json:"uid"`
}

// Add ...
func (m *Model) Add(ctx context.Context, nq *dgraph.Nquads) (bool, error) {
	if nq.UKkey == "" || nq.UKval == "" || nq.Type == "" {
		return false, errors.New("UK and Type required for Create")
	}

	q := fmt.Sprintf(`query {
		result(func: eq(%s, %s), first: 1) {
			_uid as uid
		}
	}`, nq.UKkey, util.FormatStr(nq.UKval))

	nq.ID = "_:x"
	if _, ok := nq.KV[nq.UKkey]; !ok {
		nq.KV[nq.UKkey] = nq.UKval
	}
	data, err := nq.Bytes()
	if err != nil {
		return false, err
	}

	r := make([]*jsonUID, 0)
	out := &otgo.Response{Result: &r}
	err = m.Do(ctx, q, nil, out, &api.Mutation{
		Cond:      "@if(eq(len(_uid), 0))",
		SetNquads: data,
	})
	if err != nil {
		return false, err
	}
	if !isIdempotent(ctx) && len(r) > 0 {
		return false, gear.ErrConflict.WithMsgf("%s %s exists", nq.Type, r[0].UID)
	}
	return len(r) == 0, nil
}

// BatchAdd ...
func (m *Model) BatchAdd(ctx context.Context, nqs []*dgraph.Nquads) ([]string, error) {
	uids := make([]string, 0, len(nqs))
	if len(nqs) == 0 {
		return uids, nil
	}

	qs := make([]string, 0, len(nqs))
	mus := make([]*api.Mutation, 0, len(nqs))

	for i, nq := range nqs {
		if nq.UKkey == "" || nq.UKval == "" || nq.Type == "" {
			return nil, errors.New("UK and Type required for Create")
		}
		uidx := fmt.Sprintf("uid_%d", i)
		qs = append(qs, fmt.Sprintf("%s as q_%d(func: eq(%s, %s), first: 1)", uidx, i, nq.UKkey, util.FormatStr(nq.UKval)))
		nq.ID = "_:" + uidx
		if _, ok := nq.KV[nq.UKkey]; !ok {
			nq.KV[nq.UKkey] = nq.UKval
		}
		data, err := nq.Bytes()
		if err != nil {
			return nil, err
		}
		mus = append(mus, &api.Mutation{
			Cond:      fmt.Sprintf("@if(eq(len(%s), 0))", uidx),
			SetNquads: data,
		})
	}

	q := fmt.Sprintf(`query {
		%s
	}`, strings.Join(qs, "\n"))

	res, err := m.DoRaw(ctx, q, nil, mus...)
	if err != nil {
		return nil, err
	}
	for _, uid := range res.Uids {
		uids = append(uids, uid)
	}
	return uids, nil
}

type jsonCheckCyclic struct {
	CyclicData []json.RawMessage `json:"CyclicData"`
	Result     interface{}       `json:"result"`
}

// BatchAddOrUpdate ...
func (m *Model) BatchAddOrUpdate(ctx context.Context, nqs []*dgraph.Nquads, checkCyclic string) error {
	if len(nqs) == 0 {
		return nil
	}

	qs := make([]string, 0, 1+len(nqs)/2)
	if checkCyclic != "" {
		if !strings.Contains(checkCyclic, "CyclicData") && !strings.Contains(checkCyclic, "CyclicUID") {
			return gear.ErrInternalServerError.WithMsgf("invalid checkCyclic query: %s", util.FormatStr(checkCyclic))
		}
		qs = append(qs, checkCyclic)
	}
	mus := make([]*api.Mutation, 0, len(nqs))

	for i := 0; i < len(nqs); i += 2 {
		if nqs[i].UKkey == "" || nqs[i].UKval == "" || nqs[i].Type == "" {
			return errors.New("UK and Type required for Create")
		}
		uidx := fmt.Sprintf("uid_%d", i)
		qs = append(qs, fmt.Sprintf("%s as q_%d(func: eq(%s, %s), first: 1)", uidx, i, nqs[i].UKkey, util.FormatStr(nqs[i].UKval)))
		nqs[i].ID = "_:" + uidx
		if _, ok := nqs[i].KV[nqs[i].UKkey]; !ok {
			nqs[i].KV[nqs[i].UKkey] = nqs[i].UKval
		}
		data, err := nqs[i].Bytes()
		if err != nil {
			return err
		}
		nqs[i+1].ID = fmt.Sprintf("uid(%s)", uidx)
		updateData, err := nqs[i+1].Bytes()
		if err != nil {
			return err
		}

		mus = append(mus, &api.Mutation{
			Cond:      fmt.Sprintf("@if(eq(len(%s), 0))", uidx),
			SetNquads: data,
		})
		if checkCyclic != "" {
			mus = append(mus, &api.Mutation{
				Cond:      fmt.Sprintf("@if(eq(len(%s), 1) AND eq(len(CyclicUID), 0))", uidx),
				SetNquads: updateData,
			})
		} else {
			mus = append(mus, &api.Mutation{
				Cond:      fmt.Sprintf("@if(eq(len(%s), 1))", uidx),
				SetNquads: updateData,
			})
		}
	}

	q := fmt.Sprintf(`query {
		%s
	}`, strings.Join(qs, "\n"))

	out := &jsonCheckCyclic{}
	err := m.Do(ctx, q, nil, out, mus...)
	if err != nil {
		return err
	}
	if checkCyclic != "" && len(out.CyclicData) > 0 {
		return gear.ErrConflict.WithMsgf("cyclic graph will come into being: %s", string(out.CyclicData[0]))
	}
	return nil
}

// Update ...
func (m *Model) Update(ctx context.Context, nq *dgraph.Nquads, checkCyclic string) error {
	qs := make([]string, 0, 1)
	if checkCyclic != "" {
		if !strings.Contains(checkCyclic, "CyclicData") && !strings.Contains(checkCyclic, "CyclicUID") {
			return gear.ErrInternalServerError.WithMsgf("invalid checkCyclic query: %s", util.FormatStr(checkCyclic))
		}
		qs = append(qs, checkCyclic)
	}
	if nq.ID == "" {
		if nq.UKkey == "" || nq.UKval == "" {
			return errors.New("UK required for Update")
		}
		qs = append(qs, fmt.Sprintf(`
		  result(func: eq(%s, %s), first: 1) {
				_uid as uid
			}`, nq.UKkey, util.FormatStr(nq.UKval)))
		nq.ID = "uid(_uid)"
	}

	updateData, err := nq.Bytes()
	if err != nil {
		return err
	}

	conds := make([]string, 0, 1)
	if nq.ID == "" {
		conds = append(conds, "eq(len(_uid), 1)")
	}
	if checkCyclic != "" {
		conds = append(conds, "eq(len(CyclicUID), 0)")
	}

	cond := ""
	if len(conds) > 0 {
		cond = fmt.Sprintf("@if(%s)", strings.Join(conds, " AND "))
	}
	mu := &api.Mutation{
		Cond:      cond,
		SetNquads: updateData,
	}

	q := ""
	if len(qs) > 0 {
		q = fmt.Sprintf(`query {
			%s
		}`, strings.Join(qs, "\n"))
	}

	r := make([]jsonUID, 0)
	out := &jsonCheckCyclic{Result: &r}
	err = m.Do(ctx, q, nil, out, mu)
	if err != nil {
		return err
	}
	if checkCyclic != "" && len(out.CyclicData) > 0 {
		return gear.ErrConflict.WithMsgf("cyclic graph will come into being: %s", string(out.CyclicData[0]))
	}
	if nq.ID == "" {
		if len(r) == 0 {
			return gear.ErrNotFound.WithMsgf("no resources to update")
		}
		if len(r) > 1 {
			return gear.ErrUnprocessableEntity.WithMsgf("unexpected resources %v", r)
		}
	}
	return nil
}

// Get ...
func (m *Model) Get(ctx context.Context, query string, vars map[string]string, one interface{}) error {
	res := make([]json.RawMessage, 0, 1)
	out := &otgo.Response{Result: &res}
	err := m.Query(ctx, query, vars, out)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return gear.ErrNotFound.WithMsgf("resource not found")
	}
	if len(res) > 1 {
		return gear.ErrInternalServerError.WithMsgf("unexpected resources: %d", len(res))
	}
	return json.Unmarshal(res[0], one)
}

// List ...
func (m *Model) List(ctx context.Context, query string, vars map[string]string, slice interface{}) error {
	out := &otgo.Response{Result: slice}
	return m.Query(ctx, query, vars, out)
}

type jsonUnitObjectScope struct {
	Unit   []jsonUID `json:"unit"`
	Object []jsonUID `json:"object"`
	Scope  []jsonUID `json:"scope"`
}

func (m *Model) acquireUnitObjectScope(ctx context.Context, tenant tpl.Tenant, unit, object, scope *tpl.Target, status int) (string, string, string, error) {

	q := ""
	if unit != nil {
		uk := util.HashBase64(tenant.Tenant, unit.Type, unit.ID)
		q += fmt.Sprintf("\nunit(func: eq(OTAC.U.UK, %s), first: 1) @filter(ge(OTAC.status, %d)) { uid }", util.FormatStr(uk), status)
	}
	if object != nil {
		uk := util.HashBase64(tenant.Tenant, object.Type, object.ID)
		q += fmt.Sprintf("\nobject(func: eq(OTAC.O.UK, %s), first: 1) { uid }", util.FormatStr(uk))
	}
	if scope != nil {
		uk := util.HashBase64(tenant.Tenant, scope.Type, scope.ID)
		q += fmt.Sprintf("\nscope(func: eq(OTAC.Sc.UK, %s), first: 1) @filter(ge(OTAC.status, %d)) { uid }", util.FormatStr(uk), status)
	}
	if q == "" {
		return "", "", "", nil
	}
	q = fmt.Sprintf(`query {
		%s
	}`, q)
	res := &jsonUnitObjectScope{}
	err := m.Query(ctx, q, nil, res)
	if err != nil {
		return "", "", "", err
	}
	unitUID := ""
	objectUID := ""
	scopeUID := ""
	if unit != nil {
		if len(res.Unit) == 0 {
			return "", "", "", gear.ErrNotFound.WithMsgf("Unit(%s, %s) not found", unit.Type, unit.ID)
		}
		unitUID = res.Unit[0].UID
	}
	if object != nil {
		if len(res.Object) == 0 {
			return "", "", "", gear.ErrNotFound.WithMsgf("Object(%s, %s) not found", object.Type, object.ID)
		}
		objectUID = res.Object[0].UID
	}
	if scope != nil {
		if len(res.Scope) == 0 {
			return "", "", "", gear.ErrNotFound.WithMsgf("Scope(%s, %s) not found", scope.Type, scope.ID)
		}
		scopeUID = res.Scope[0].UID
	}
	return unitUID, objectUID, scopeUID, nil
}

func (m *Model) acquirePermissions(ctx context.Context, tenant tpl.Tenant, permissions []string) ([]tpl.Permission, error) {
	uks := make([]string, len(permissions))
	for i, p := range permissions {
		uks[i] = util.HashBase64(tenant.Tenant, p)
	}
	q := fmt.Sprintf(`query {
		result(func: eq(OTAC.P.UK, [%s]), first: %d) {
			uid
			permission: OTAC.P
		}
	}`, strings.Join(util.FormatStrs(uks), ", "), len(uks))
	out := make([]tpl.Permission, 0, len(uks))
	if err := m.List(ctx, q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type jsonOrgOU struct {
	Org []jsonUID `json:"org"`
	OU  []jsonUID `json:"ou"`
}

func (m *Model) acquireOrgOU(ctx context.Context, org, ou string, status int) (string, string, error) {
	if org == "" {
		return "", "", gear.ErrBadRequest.WithMsg("organization required")
	}
	q := fmt.Sprintf("\norg(func: eq(OTAC.Org, %s), first: 1) @filter(ge(OTAC.status, %d)) { uid }", util.FormatStr(org), status)
	if ou != "" {
		uk := util.HashBase64(org, ou)
		q += fmt.Sprintf("\nou(func: eq(OTAC.OU.UK, %s), first: 1) @filter(ge(OTAC.status, %d)) { uid }", util.FormatStr(uk), status)
	}

	q = fmt.Sprintf(`query {
		%s
	}`, q)
	res := &jsonOrgOU{}
	err := m.Query(ctx, q, nil, res)
	if err != nil {
		return "", "", err
	}
	if len(res.Org) == 0 {
		return "", "", gear.ErrNotFound.WithMsgf("Organization(%s) not found", org)
	}
	orgUID := res.Org[0].UID
	ouUID := ""
	if ou != "" {
		if len(res.OU) == 0 {
			return "", "", gear.ErrNotFound.WithMsgf("OU(%s, %s) not found", org, ou)
		}
		ouUID = res.OU[0].UID
	}
	return orgUID, ouUID, nil
}

func (m *Model) acquireOrgMembers(ctx context.Context, org string, subjects []string, status int) ([]string, error) {
	if org == "" {
		return nil, gear.ErrBadRequest.WithMsg("organization required")
	}
	if len(subjects) == 0 {
		return make([]string, 0), nil
	}

	uks := make([]string, len(subjects))
	for i, sub := range subjects {
		uks[i] = util.HashBase64(org, sub)
	}

	q := fmt.Sprintf(`query {
		result(func: eq(OTAC.M.UK, [%s]), first: %d) @filter(ge(OTAC.status, %d)) { uid }
	}`, strings.Join(util.FormatStrs(uks), ", "), len(uks), status)

	res := make([]jsonUID, 0, len(uks))
	if err := m.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	uids := make([]string, len(res))
	for i, sub := range res {
		uids[i] = sub.UID
	}

	return uids, nil
}
