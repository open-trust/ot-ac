package tpl

import "github.com/teambition/gear"

// AddUnitInput ...
type AddUnitInput struct {
	ObjectType string `json:"objectType"`
	Operation  string `json:"operation"`
}

// Validate 实现 gear.BodyTemplate
func (t *AddUnitInput) Validate() error {
	if t.ObjectType == "" {
		return gear.ErrBadRequest.WithMsgf("objectType required")
	}
	if t.Operation == "" {
		return gear.ErrBadRequest.WithMsgf("operation required")
	}
	if len(t.ObjectType) > 64 {
		return gear.ErrBadRequest.WithMsgf("objectType size too large")
	}
	if len(t.Operation) > 64 {
		return gear.ErrBadRequest.WithMsgf("operation size too large")
	}
	return nil
}
