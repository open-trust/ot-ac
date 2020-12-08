package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/open-trust/ot-ac/src/service/dgraph"
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
	Model      *Model
	Permission *Permission
	Scope      *Scope
	Unit       *Unit
	Tenant     *Tenant
	Subject    *Subject
}

// NewModels ...
func NewModels(dg *dgraph.Dgraph) *Models {
	m := &Model{dg}
	return &Models{
		Model:      m,
		Permission: &Permission{m},
		Scope:      &Scope{m},
		Unit:       &Unit{m},
		Tenant:     &Tenant{m},
		Subject:    &Subject{m},
	}
}

type contextKey int

const (
	idempotentKey contextKey = iota
)

// ContextWithPrefer ...
func ContextWithPrefer(ctx *gear.Context) context.Context {
	c := ctx.Context()
	idempotent := true
	for _, prefer := range ctx.QueryAll("Prefer") {
		if prefer == "respond-conflict" {
			idempotent = false
		}
		// other prefer ...
	}
	if !idempotent {
		c = context.WithValue(c, idempotentKey, false)
	}

	return c
}

func isIdempotent(ctx context.Context) bool {
	return ctx.Value(idempotentKey) == nil
}

type jsonUID struct {
	UID string `json:"uid"`
}

// Add ...
func (m *Model) Add(ctx context.Context, nq *dgraph.Nquads) error {
	if nq.UKkey == "" || nq.UKval == "" || nq.Type == "" {
		return errors.New("UK and Type required for Create")
	}

	q := fmt.Sprintf(`query {
		result(func: eq(%s, %s), first: 1) {
			_uid as uid
		}
	}`, nq.UKkey, strconv.Quote(nq.UKval))

	nq.ID = "_:x"
	nq.KV[nq.UKkey] = nq.UKval
	data, err := nq.Bytes()
	if err != nil {
		return err
	}

	r := make([]*jsonUID, 0)
	out := &otgo.Response{Result: &r}
	err = m.Do(ctx, q, nil, out, &api.Mutation{
		Cond:      "@if(eq(len(_uid), 0))",
		SetNquads: data,
	})
	if err != nil {
		return err
	}
	if !isIdempotent(ctx) && len(r) > 0 {
		return gear.ErrConflict.WithMsgf("%s %s exists", nq.Type, r[0].UID)
	}
	return nil
}

// BatchAdd ...
func (m *Model) BatchAdd(ctx context.Context, nqs []*dgraph.Nquads) error {
	if len(nqs) == 0 {
		return nil
	}

	qs := make([]string, 0, len(nqs))
	mus := make([]*api.Mutation, 0, len(nqs))

	for i, nq := range nqs {
		if nq.UKkey == "" || nq.UKval == "" || nq.Type == "" {
			return errors.New("UK and Type required for Create")
		}
		uidx := fmt.Sprintf("uid_%d", i)
		qs = append(qs, fmt.Sprintf("%s as q_%d(func: eq(%s, %s), first: 1)", uidx, i, nq.UKkey, strconv.Quote(nq.UKval)))
		nq.ID = "_:" + uidx
		nq.KV[nq.UKkey] = nq.UKval
		data, err := nq.Bytes()
		if err != nil {
			return err
		}
		fmt.Println("==========", string(data))
		mus = append(mus, &api.Mutation{
			Cond:      fmt.Sprintf("@if(eq(len(%s), 0))", uidx),
			SetNquads: data,
		})
	}

	q := fmt.Sprintf(`query {
		%s
	}`, strings.Join(qs, "\n"))
	fmt.Println("==========", q)

	return m.Do(ctx, q, nil, nil, mus...)
}

// Update ...
func (m *Model) Update(ctx context.Context, nq *dgraph.Nquads) error {
	if nq.UKkey == "" || nq.UKval == "" {
		return errors.New("UK required for Update")
	}

	q := fmt.Sprintf(`query {
		result(func: eq(%s, %s), first: 1) {
			_uid as uid
		}
	}`, nq.UKkey, strconv.Quote(nq.UKval))

	nq.ID = "uid(_uid)"
	data, err := nq.Bytes()
	if err != nil {
		return err
	}

	r := make([]*jsonUID, 0)
	out := &otgo.Response{Result: &r}
	err = m.Do(ctx, q, nil, out, &api.Mutation{
		Cond:      "@if(eq(len(_uid), 1))",
		SetNquads: data,
	})
	if err != nil {
		return err
	}
	if len(r) == 0 {
		return gear.ErrNotFound.WithMsgf("no resources to update")
	}
	if len(r) > 1 {
		return gear.ErrUnprocessableEntity.WithMsgf("unexpected resources %v", r)
	}
	return nil
}

// CreateOrUpdate ...
func (m *Model) CreateOrUpdate(ctx context.Context, qs string, create, update *dgraph.Nquads) error {
	if create.UKkey == "" || create.UKval == "" || create.Type == "" {
		return errors.New("UK and Type required for CreateOrUpdate")
	}

	q := fmt.Sprintf(`query {
		result(func: eq(%s, %s), first: 1) {
			_uid as uid
		}
		%s
	}`, create.UKkey, strconv.Quote(create.UKval), qs)

	create.ID = "_:x"
	create.KV[create.UKkey] = create.UKval
	createData, err := create.Bytes()
	if err != nil {
		return err
	}

	update.ID = "uid(_uid)"
	updateData, err := update.Bytes()
	if err != nil {
		return err
	}

	r := make([]*jsonUID, 0)
	out := &otgo.Response{Result: &r}
	err = m.Do(ctx, q, nil, out, &api.Mutation{
		Cond:      "@if(eq(len(_uid), 0))",
		SetNquads: createData,
	}, &api.Mutation{
		Cond:      "@if(eq(len(_uid), 1))",
		SetNquads: updateData,
	})
	if err != nil {
		return err
	}
	if len(r) > 1 {
		return gear.ErrUnprocessableEntity.WithMsgf("unexpected resources: %v", r)
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
