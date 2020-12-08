package tpl

import (
	"github.com/teambition/gear"
)

// UpdateSubjectInput ...
type UpdateSubjectInput struct {
	Subject string `json:"subject"`
	Status  int    `json:"status"`
}

// Validate 实现 gear.BodyTemplate
func (t *UpdateSubjectInput) Validate() error {
	// OTID UnmarshalText method will validate
	if err := CheckSubject(t.Subject); err != nil {
		return err
	}
	if t.Status < -1 {
		return gear.ErrBadRequest.WithMsgf("invalid subject status %d", t.Status)
	}
	return nil
}

// BatchAddSubjectsInput ...
type BatchAddSubjectsInput struct {
	Subjects []string `json:"subjects"`
}

// Validate 实现 gear.BodyTemplate
func (t *BatchAddSubjectsInput) Validate() error {
	// OTID UnmarshalText method will validate
	if len(t.Subjects) == 0 {
		return gear.ErrBadRequest.WithMsg("empty subjects")
	}
	for _, subject := range t.Subjects {
		if err := CheckSubject(subject); err != nil {
			return err
		}
	}
	return nil
}
