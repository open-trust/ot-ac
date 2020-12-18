package tpl

import "github.com/teambition/gear"

// Object ...
type Object struct {
	UID        string `json:"uid,omitempty"`
	TargetID   string `json:"targetId"`
	TargetType string `json:"targetType"`
	Terms      string `json:"terms,omitempty"`
}

// ObjectAddPermissionsInput ...
type ObjectAddPermissionsInput struct {
	Target
	Permissions []string `json:"permissions"`
}

// Validate 实现 gear.BodyTemplate
func (t *ObjectAddPermissionsInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if len(t.Permissions) == 0 {
		return gear.ErrBadRequest.WithMsgf("permissions empty")
	}
	if len(t.Permissions) > 100 {
		return gear.ErrBadRequest.WithMsgf("too many permissions: %d", len(t.Permissions))
	}
	for _, p := range t.Permissions {
		if err := CheckPermission(p); err != nil {
			return err
		}
	}
	return nil
}
