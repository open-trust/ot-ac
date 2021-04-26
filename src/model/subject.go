package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-trust/ot-ac/src/service/dgraph"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"
)

// Subject ...
type Subject struct {
	*Model
}

// List ...
func (m *Subject) List(ctx context.Context, pageSize, skip int, uidToken string) ([]tpl.Subject, error) {
	q := fmt.Sprintf(`query {
		result(func: eq(dgraph.type, "OTACSubject"), first: %d, offset: %d, after: %s) {
			uid
			status: OTAC.status
			subject: OTAC.Sub
		}
	}`, pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.Subject, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// BatchAdd ...
func (m *Subject) BatchAdd(ctx context.Context, input []string) error {
	nqs := make([]*dgraph.Nquads, 0, len(input))

	for _, sub := range input {
		nqs = append(nqs, &dgraph.Nquads{
			UKkey: "OTAC.Sub",
			UKval: sub,
			Type:  "OTACSubject",
			KV: map[string]interface{}{
				"OTAC.status": 0,
			},
		})
	}

	_, err := m.Model.BatchAdd(ctx, nqs)
	return err
}

// AcquireUIDsOrAdd ...
func (m *Subject) AcquireUIDsOrAdd(ctx context.Context, input []string) ([]tpl.Subject, error) {
	size := len(input)
	q := fmt.Sprintf(`query {
		result(func: eq(OTAC.Sub, [%s]), first: %d) {
			uid
			subject: OTAC.Sub
		}
	}`, strings.Join(util.FormatStrs(input), ", "), size)
	out := make([]tpl.Subject, 0, size)
	if err := m.Model.List(ctx, q, nil, &out); err != nil {
		return nil, err
	}

	if len(out) < size {
		nqs := make([]*dgraph.Nquads, 0, size-len(out))
		newSubjects := make([]string, 0, size-len(out))
		for _, sub := range input {
			if tpl.GetSubjectUID(out, sub) == "" {
				newSubjects = append(newSubjects, sub)
				nqs = append(nqs, &dgraph.Nquads{
					UKkey: "OTAC.Sub",
					UKval: sub,
					Type:  "OTACSubject",
					KV: map[string]interface{}{
						"OTAC.status": 0,
					},
				})
			}
		}

		_, err := m.Model.BatchAdd(ctx, nqs)
		if err != nil {
			return nil, err
		}

		q = fmt.Sprintf(`query {
			result(func: eq(OTAC.Sub, [%s]), first: %d) {
				uid
				subject: OTAC.Sub
			}
		}`, strings.Join(util.FormatStrs(newSubjects), ", "), len(newSubjects))
		newOut := make([]tpl.Subject, 0, len(newSubjects))
		if err := m.Model.List(ctx, q, nil, &out); err != nil {
			return nil, err
		}
		out = append(out, newOut...)
	}
	return out, nil
}

// Update ...
func (m *Subject) Update(ctx context.Context, input tpl.Subject) error {
	update := &dgraph.Nquads{
		UKkey: "OTAC.Sub",
		UKval: input.Sub,
		Type:  "OTACSubject",
		KV: map[string]interface{}{
			"OTAC.status": input.Status,
		},
	}

	return m.Model.Update(ctx, update, "")
}

// // ListUnits ...
// func (m *Subject) ListUnits(ctx context.Context, subjectID otgo.OTID) (*tpl.Subject, error) {
// 	q := fmt.Sprintf(`query {
// 		result(func: eq(OTAC.S.OTID, %s)) {
// 			uid
// 			status: OTAC.status
// 			otid: OTAC.S.OTID
// 			units: ~OTAC.U-Ss @facets(kind: kind, name: name, level: level) {
// 				uid
// 				status: OTAC.status
// 				unitId: OTAC.UId
// 				unitType: OTAC.UType
// 				hasPermissions: OTAC.U-Ps @facets(extention: ex) {
// 					uid
// 					objectType: OTAC.OType
// 					operation: OTAC.OP
// 				}
// 			}
// 		}
// 	`, util.FormatStr(subjectID.String()))
// 	res := tpl.Subject{
// 		Units: make([]*tpl.SubjectUnit, 0),
// 	}
// 	if err := m.Model.Get(ctx, q, nil, &res); err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }
