package model

import (
	"context"
	"fmt"

	"github.com/open-trust/ot-ac/src/schema"
	"github.com/open-trust/ot-ac/src/service/dgraph"
)

// Subject ...
type Subject struct {
	*Model
}

// List ...
func (m *Subject) List(ctx context.Context, pageSize, skip int, uidToken string) ([]*schema.Subject, error) {
	q := fmt.Sprintf(`query {
		result(func: has(OTAC.Sub), first: %d, offset: %d, after: <%s>) {
			id: uid
			status: OTAC.status
			subject: OTAC.Sub
		}
	}`, pageSize, skip, uidToken)
	res := make([]*schema.Subject, 0, pageSize)
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

	return m.Model.BatchAdd(ctx, nqs)
}

// Update ...
func (m *Subject) Update(ctx context.Context, input *schema.Subject) error {
	update := &dgraph.Nquads{
		UKkey: "OTAC.Sub",
		UKval: input.Subject,
		Type:  "OTACSubject",
		KV: map[string]interface{}{
			"OTAC.status": input.Status,
		},
	}

	return m.Model.Update(ctx, update)
}

// // ListUnits ...
// func (m *Subject) ListUnits(ctx context.Context, subjectID otgo.OTID) (*schema.Subject, error) {
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
// 	`, strconv.Quote(subjectID.String()))
// 	res := schema.Subject{
// 		Units: make([]*schema.SubjectUnit, 0),
// 	}
// 	if err := m.Model.Get(ctx, q, nil, &res); err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }
