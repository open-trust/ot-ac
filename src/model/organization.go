package model

import (
	"context"
	"fmt"

	"github.com/open-trust/ot-ac/src/service/dgraph"
	"github.com/open-trust/ot-ac/src/tpl"
	"github.com/open-trust/ot-ac/src/util"
	"github.com/teambition/gear"
)

// Organization ...
type Organization struct {
	*Model
}

// AddOrg ...
func (m *Organization) AddOrg(ctx context.Context, org string) (bool, error) {
	nq := &dgraph.Nquads{
		UKkey: "OTAC.Org",
		UKval: org,
		Type:  "OTACOrg",
		KV: map[string]interface{}{
			"OTAC.status": 0,
		},
	}

	return m.Model.Add(ctx, nq)
}

// UpdateOrgStatus ...
func (m *Organization) UpdateOrgStatus(ctx context.Context, org string, status int) error {
	update := &dgraph.Nquads{
		UKkey: "OTAC.Org",
		UKval: org,
		Type:  "OTACOrg",
		KV: map[string]interface{}{
			"OTAC.status": status,
		},
	}

	return m.Model.Update(ctx, update, "")
}

// DeleteOrg ...
func (m *Organization) DeleteOrg(ctx context.Context, org string) error {
	return nil
}

// ListOrgs ...
func (m *Organization) ListOrgs(ctx context.Context, pageSize, skip int, uidToken string) ([]tpl.Organization, error) {
	q := fmt.Sprintf(`query {
		result(func: eq(dgraph.type, "OTACOrg"), first: %d, offset: %d, after: %s) {
			uid
			organization: OTAC.Org
			status: OTAC.status
		}
	}`, pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.Organization, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListSubjectOrgs ...
func (m *Organization) ListSubjectOrgs(ctx context.Context, subject string, pageSize, skip int, uidToken string) ([]tpl.Organization, error) {
	q := fmt.Sprintf(`query {
		var(func: eq(OTAC.Sub, %s), first: 1) @filter(ge(OTAC.status, %d)) {
			~OTAC.M-S @filter(ge(OTAC.status, %d)) {
				orgUIDs as OTAC.M-Org @filter(ge(OTAC.status, %d))
			}
		}
		result(func: uid(orgUIDs), first: %d, offset: %d, after: %s) {
			uid
			organization: OTAC.Org
			status: OTAC.status
		}
	}`, util.FormatStr(subject), 0, 0, 0, pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.Organization, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// AddOU ...
func (m *Organization) AddOU(ctx context.Context, org string, input tpl.OrganizationAddOUInput) (bool, error) {
	orgUID, parentUID, err := m.acquireOrgOU(ctx, org, input.Parent, 0)
	if err != nil {
		return false, err
	}
	nq := &dgraph.Nquads{
		UKkey: "OTAC.OU.UK",
		UKval: util.HashBase64(org, input.OU),
		Type:  "OTACOU",
		KV: map[string]interface{}{
			"OTAC.status":   0,
			"OTAC.OU":       input.OU,
			"OTAC.OU-Org":   util.FormatUID(orgUID),
			"OTAC.OU.terms": input.Terms,
		},
	}
	if parentUID != "" {
		nq.KV["OTAC.OU-OU"] = util.FormatUID(parentUID)
	}

	return m.Model.Add(ctx, nq)
}

// UpdateOUParent ...
func (m *Organization) UpdateOUParent(ctx context.Context, org string, input tpl.OrganizationUpdateOUParentInput) error {
	_, parentUID, err := m.acquireOrgOU(ctx, org, input.Parent, 0)
	if err != nil {
		return err
	}
	_, ouUID, err := m.acquireOrgOU(ctx, org, input.OU, 0)
	if err != nil {
		return err
	}
	checkCyclic := fmt.Sprintf(`
		var(func: uid(%s), first: 1) @recurse(loop: false) {
			uids as uid
			OTAC.OU-OU
		}
		CyclicData(func: uid(uids), first: 1) @filter(uid(%s)) {
			CyclicUID as uid
			ou: OTAC.OU
		}
	`, util.FormatUID(parentUID), util.FormatUID(ouUID))

	nq := &dgraph.Nquads{
		ID: util.FormatUID(ouUID),
		KV: map[string]interface{}{
			"OTAC.OU-OU": util.FormatUID(parentUID),
		},
	}
	return m.Model.Update(ctx, nq, checkCyclic)
}

// UpdateOUStatus ...
func (m *Organization) UpdateOUStatus(ctx context.Context, org string) error {
	return nil
}

// UpdateOUTerms ...
func (m *Organization) UpdateOUTerms(ctx context.Context, org string) error {
	return nil
}

// DeleteOU ...
func (m *Organization) DeleteOU(ctx context.Context, org string) error {
	return nil
}

// ListOUs ...
func (m *Organization) ListOUs(ctx context.Context, org, parent string, pageSize, skip int, uidToken string) ([]tpl.OU, error) {
	q := ""
	if parent == "" {
		q = fmt.Sprintf(`query {
			var(func: eq(OTAC.Org, %s), first: 1) {
				orgUID as uid
			}
			result(func: eq(dgraph.type, "OTACOU"), first: %d, offset: %d, after: %s) @filter(uid_in(OTAC.OU-Org, uid(orgUID)) AND NOT has(OTAC.OU-OU)) @normalize {
				uid: uid
				ou: OTAC.OU
				status: OTAC.status
			}
		}`, util.FormatStr(org), pageSize, skip, util.FormatUID(uidToken))
	} else {
		uk := util.HashBase64(org, parent)
		q = fmt.Sprintf(`query {
			result(func: eq(OTAC.OU.UK, %s), first: 1) @normalize {
				~OTAC.OU-OU (first: %d, offset: %d, after: %s) {
					uid: uid
					ou: OTAC.OU
					status: OTAC.status
					OTAC.OU-OU {
						parent: OTAC.OU
					}
				}
			}
		}`, util.FormatStr(uk), pageSize, skip, util.FormatUID(uidToken))
	}
	res := make([]tpl.OU, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListSubjectOUs ...
func (m *Organization) ListSubjectOUs(ctx context.Context, subject, org string, pageSize, skip int, uidToken string) ([]tpl.OU, error) {
	q := fmt.Sprintf(`query {
		var(func: eq(OTAC.Org, %s), first: 1) @filter(ge(OTAC.status, %d)) {
			orgUID as uid
		}
		var(func: eq(OTAC.Sub, %s), first: 1) @filter(ge(OTAC.status, %d)) {
			~OTAC.M-S @filter(uid_in(OTAC.M-Org, uid(orgUID)) AND ge(OTAC.status, %d)) {
				ouUIDs as ~OTAC.OU-Ms @filter(ge(OTAC.status, %d))
			}
		}
		result(func: uid(ouUIDs), first: %d, offset: %d, after: %s) @normalize {
			uid: uid
			ou: OTAC.OU
			status: OTAC.status
			OTAC.OU-OU {
				parent: OTAC.OU
			}
		}
	}`, util.FormatStr(org), 0, util.FormatStr(subject), 0, 0, 0, pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.OU, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// SearchOUs ...
func (m *Organization) SearchOUs(ctx context.Context, org, term string, pageSize, skip int, uidToken string) ([]tpl.OU, error) {
	q := fmt.Sprintf(`query {
		var(func: eq(OTAC.Org, %s), first: 1) {
			orgUID as uid
		}
		result(func: eq(dgraph.type, "OTACOU"), first: %d, offset: %d, after: %s) @filter(uid_in(OTAC.OU-Org, uid(orgUID)) AND allofterms(OTAC.OU.terms, %s)) @normalize {
			uid: uid
			ou: OTAC.OU
			status: OTAC.status
			OTAC.OU-OU {
				parent: OTAC.OU
			}
		}
	}`, util.FormatStr(org), pageSize, skip, util.FormatUID(uidToken), util.FormatStr(term))
	res := make([]tpl.OU, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// BatchAddMember ...
func (m *Organization) BatchAddMember(ctx context.Context, org string, input tpl.OrganizationBatchAddMemberInput) error {
	orgUID, _, err := m.acquireOrgOU(ctx, org, "", 0)
	if err != nil {
		return err
	}
	nqs := make([]*dgraph.Nquads, 0, len(input.Subjects))

	for _, sub := range input.Subjects {
		if sub.UID == "" {
			return gear.ErrBadRequest.WithMsgf("Subject(%s) not found", sub.Sub)
		}
		nqs = append(nqs, &dgraph.Nquads{
			UKkey: "OTAC.M.UK",
			UKval: util.HashBase64(org, sub.Sub),
			Type:  "OTACMember",
			KV: map[string]interface{}{
				"OTAC.status":  0,
				"OTAC.M-S":     util.FormatUID(sub.UID),
				"OTAC.M-Org":   util.FormatUID(orgUID),
				"OTAC.M.terms": sub.Terms,
			},
		})
	}

	_, err = m.Model.BatchAdd(ctx, nqs)
	return err
}

// UpdateMemberStatus ...
func (m *Organization) UpdateMemberStatus(ctx context.Context, org string) error {
	return nil
}

// UpdateMemberTerms ...
func (m *Organization) UpdateMemberTerms(ctx context.Context, org string) error {
	return nil
}

// RemoveMember ...
func (m *Organization) RemoveMember(ctx context.Context, org string) error {
	return nil
}

// ListMembers ...
func (m *Organization) ListMembers(ctx context.Context, org string, pageSize, skip int, uidToken string) ([]tpl.Member, error) {
	q := fmt.Sprintf(`query {
		result(func: eq(OTAC.Org, %s), first: 1) @normalize {
			~OTAC.M-Org (first: %d, offset: %d, after: %s) {
				uid: uid
				status: OTAC.status
				OTAC.M-S {
					subject: OTAC.Sub
				}
			}
		}
	}`, util.FormatStr(org), pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.Member, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// SearchMember ...
func (m *Organization) SearchMember(ctx context.Context, org, term string, pageSize, skip int, uidToken string) ([]tpl.Member, error) {
	q := fmt.Sprintf(`query {
		result(func: eq(OTAC.Org, %s), first: 1) @normalize {
			~OTAC.M-Org @filter(allofterms(OTAC.M.terms, %s)) (first: %d, offset: %d, after: %s) {
				uid: uid
				status: OTAC.status
				OTAC.M-S {
					subject: OTAC.Sub
				}
			}
		}
	}`, util.FormatStr(org), util.FormatStr(term), pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.Member, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// BatchAddOUMember ...
func (m *Organization) BatchAddOUMember(ctx context.Context, org string, input tpl.OrganizationBatchAddOUMemberInput) error {
	_, ouUID, err := m.acquireOrgOU(ctx, org, input.OU, 0)
	if err != nil {
		return err
	}
	memberUIDs, err := m.acquireOrgMembers(ctx, org, input.Subjects, 0)
	if err != nil {
		return err
	}
	if len(memberUIDs) == 0 {
		return nil
	}
	nq := &dgraph.Nquads{
		ID: util.FormatUID(ouUID),
		KV: map[string]interface{}{
			"OTAC.OU-Ms": util.FormatUIDs(memberUIDs),
		},
	}
	return m.Model.Update(ctx, nq, "")
}

// RemoveOUMember ...
func (m *Organization) RemoveOUMember(ctx context.Context, org string) error {
	return nil
}

// ListOUMembers ...
func (m *Organization) ListOUMembers(ctx context.Context, org, ou string, pageSize, skip int, uidToken string) ([]tpl.Member, error) {
	uk := util.HashBase64(org, ou)
	q := fmt.Sprintf(`query {
		result(func: eq(OTAC.OU.UK, %s), first: 1) @normalize {
			OTAC.OU-Ms (first: %d, offset: %d, after: %s) {
				uid: uid
				status: OTAC.status
				OTAC.M-S {
					subject: OTAC.Sub
				}
			}
		}
	}`, util.FormatStr(uk), pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.Member, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListOUDescendantMembers ...
func (m *Organization) ListOUDescendantMembers(ctx context.Context, org, ou string, pageSize, skip int, uidToken string) ([]tpl.Member, error) {
	uk := util.HashBase64(org, ou)
	q := fmt.Sprintf(`query {
		var(func: eq(OTAC.OU.UK, %s), first: 1) @recurse(loop: false) {
			ouUIDs as uid
			~OTAC.OU-OU
		}
		var(func: uid(ouUIDs)) {
			mUIDs as OTAC.OU-Ms
		 }
		result(func: uid(mUIDs), first: %d, offset: %d, after: %s) @normalize {
			uid: uid
			status: OTAC.status
			OTAC.M-S {
				subject: OTAC.Sub
			}
		}
	}`, util.FormatStr(uk), pageSize, skip, util.FormatUID(uidToken))
	res := make([]tpl.Member, 0, pageSize)
	if err := m.Model.List(ctx, q, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}
