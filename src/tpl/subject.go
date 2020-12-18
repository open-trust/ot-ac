package tpl

import (
	"github.com/teambition/gear"
)

// Subject ...
type Subject struct {
	UID    string `json:"uid,omitempty"`
	Status int    `json:"status"`
	Sub    string `json:"subject"`
}

// GetSubjectUID ...
func GetSubjectUID(ss []Subject, sub string) string {
	for _, v := range ss {
		if v.Sub == sub {
			return v.UID
		}
	}
	return ""
}

// SubjectUpdateInput ...
type SubjectUpdateInput struct {
	Sub    string `json:"subject"`
	Status int    `json:"status"`
}

// Validate 实现 gear.BodyTemplate
func (t *SubjectUpdateInput) Validate() error {
	// OTID UnmarshalText method will validate
	if err := CheckSubject(t.Sub); err != nil {
		return err
	}
	if t.Status < -1 {
		return gear.ErrBadRequest.WithMsgf("invalid subject status %d", t.Status)
	}
	return nil
}
