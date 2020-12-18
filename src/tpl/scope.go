package tpl

import (
	"github.com/teambition/gear"
)

// Scope ...
type Scope struct {
	UID        string `json:"uid,omitempty"`
	Status     int    `json:"status"`
	TargetID   string `json:"targetId"`
	TargetType string `json:"targetType"`
}

// ScopeUpdateInput ...
type ScopeUpdateInput struct {
	Target
	Status int `json:"status"`
}

// Validate 实现 gear.BodyTemplate
func (t *ScopeUpdateInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if t.Status < -1 {
		return gear.ErrBadRequest.WithMsgf("invalid scope status %d", t.Status)
	}
	return nil
}

// ScopeListInput ...
type ScopeListInput struct {
	Pagination
	TargetType string `json:"targetType"`
}

// Validate 实现 gear.BodyTemplate
func (t *ScopeListInput) Validate() error {
	if err := CheckResource(t.TargetType); err != nil {
		return err
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	return nil
}

// ScopeListUnitObjectsInput ...
type ScopeListUnitObjectsInput struct {
	Pagination
	Scope      Target `json:"scope"`
	TargetType string `json:"targetType"`
}

// Validate 实现 gear.BodyTemplate
func (t *ScopeListUnitObjectsInput) Validate() error {
	if err := CheckResource(t.TargetType); err != nil {
		return err
	}
	if err := t.Scope.Validate(); err != nil {
		return err
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	return nil
}
