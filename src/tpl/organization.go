package tpl

import (
	"github.com/teambition/gear"
)

// Organization ...
type Organization struct {
	UID    string `json:"uid,omitempty"`
	Org    string `json:"organization"`
	Status int    `json:"status"`
}

// OU ...
type OU struct {
	UID    string `json:"uid,omitempty"`
	OU     string `json:"ou"`
	Parent string `json:"parent,omitempty"`
	Status int    `json:"status"`
}

// Member ...
type Member struct {
	UID     string `json:"uid,omitempty"`
	Subject string `json:"subject"`
	Status  int    `json:"status"`
}

// OrganizationInput ...
type OrganizationInput struct {
	Org string `json:"organization"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationInput) Validate() error {
	if err := CheckSubject(t.Org); err != nil {
		return err
	}
	return nil
}

// OrganizationStatusInput ...
type OrganizationStatusInput struct {
	OrganizationInput
	Status int `json:"status"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationStatusInput) Validate() error {
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if t.Status < -1 {
		return gear.ErrBadRequest.WithMsgf("invalid Org status %d", t.Status)
	}
	return nil
}

// OrganizationAddOUInput ...
type OrganizationAddOUInput struct {
	OrganizationInput
	OU     string `json:"ou"`
	Parent string `json:"parent"`
	Terms  string `json:"terms"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationAddOUInput) Validate() error {
	if err := CheckSubject(t.OU); err != nil {
		return err
	}
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if t.Parent != "" {
		if err := CheckSubject(t.Parent); err != nil {
			return err
		}
	}
	if t.Terms != "" {
		if err := CheckTerm(t.Terms); err != nil {
			return err
		}
	}
	return nil
}

// OrganizationUpdateOUParentInput ...
type OrganizationUpdateOUParentInput struct {
	OrganizationInput
	OU     string `json:"ou"`
	Parent string `json:"parent"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationUpdateOUParentInput) Validate() error {
	if err := CheckSubject(t.OU); err != nil {
		return err
	}
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if err := CheckSubject(t.Parent); err != nil {
		return err
	}
	return nil
}

// OrganizationOUInput ...
type OrganizationOUInput struct {
	OrganizationInput
	OU string `json:"ou"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationOUInput) Validate() error {
	if err := CheckSubject(t.OU); err != nil {
		return err
	}
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	return nil
}

// OrganizationBatchAddMemberInput ...
type OrganizationBatchAddMemberInput struct {
	OrganizationInput
	Subjects []struct {
		Subject
		Terms string `json:"terms"`
	} `json:"subjects"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationBatchAddMemberInput) Validate() error {
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if len(t.Subjects) == 0 {
		return gear.ErrBadRequest.WithMsg("empty subjects")
	}
	if len(t.Subjects) > 1000 {
		return gear.ErrBadRequest.WithMsgf("too many subjects: %d", len(t.Subjects))
	}
	cr := make(checkRepetitive)
	for _, v := range t.Subjects {
		if err := cr.Check(v.Sub); err != nil {
			return err
		}
		if err := CheckSubject(v.Sub); err != nil {
			return err
		}
		if v.Terms != "" {
			if err := CheckTerm(v.Terms); err != nil {
				return err
			}
		}
	}
	return nil
}

// OrganizationUpdateMemberStatusInput ...
type OrganizationUpdateMemberStatusInput struct {
	OrganizationStatusInput
	Subject string `json:"subject"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationUpdateMemberStatusInput) Validate() error {
	if err := CheckSubject(t.Subject); err != nil {
		return err
	}
	if err := t.OrganizationStatusInput.Validate(); err != nil {
		return err
	}
	return nil
}

// OrganizationBatchAddOUMemberInput ...
type OrganizationBatchAddOUMemberInput struct {
	OrganizationOUInput
	SubjectsInput
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationBatchAddOUMemberInput) Validate() error {
	if err := t.OrganizationOUInput.Validate(); err != nil {
		return err
	}
	if err := t.SubjectsInput.Validate(); err != nil {
		return err
	}
	return nil
}

// OrganizationListInput ...
type OrganizationListInput struct {
	OrganizationInput
	Pagination
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationListInput) Validate() error {
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	return nil
}

// OrganizationSearchInput ...
type OrganizationSearchInput struct {
	OrganizationInput
	Pagination
	Term string `json:"term"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationSearchInput) Validate() error {
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if len(t.Term) < 3 {
		return gear.ErrBadRequest.WithMsgf("term' length too small %d", len(t.Term))
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	return nil
}

// OrganizationListSubjectOrgsInput ...
type OrganizationListSubjectOrgsInput struct {
	Pagination
	Subject string `json:"subject"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationListSubjectOrgsInput) Validate() error {
	if err := CheckSubject(t.Subject); err != nil {
		return err
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	return nil
}

// OrganizationListSubjectOUsInput ...
type OrganizationListSubjectOUsInput struct {
	OrganizationListSubjectOrgsInput
	OrganizationInput
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationListSubjectOUsInput) Validate() error {
	if err := t.OrganizationListSubjectOrgsInput.Validate(); err != nil {
		return err
	}
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	return nil
}

// OrganizationListOUsInput ...
type OrganizationListOUsInput struct {
	OrganizationInput
	Pagination
	Parent string `json:"parent"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationListOUsInput) Validate() error {
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	if t.Parent != "" {
		if err := CheckSubject(t.Parent); err != nil {
			return err
		}
	}
	return nil
}

// OrganizationListOUMembersInput ...
type OrganizationListOUMembersInput struct {
	OrganizationInput
	Pagination
	OU string `json:"ou"`
}

// Validate 实现 gear.BodyTemplate
func (t *OrganizationListOUMembersInput) Validate() error {
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	if err := CheckSubject(t.OU); err != nil {
		return err
	}
	return nil
}
