package tpl

import "github.com/teambition/gear"

// Unit ...
type Unit struct {
	UID        string `json:"uid,omitempty"`
	Status     int    `json:"status"`
	TargetID   string `json:"targetId"`
	TargetType string `json:"targetType"`
}

// UnitAddInput ...
type UnitAddInput struct {
	Target
	Parent *Target `json:"parent"`
	Scope  *Target `json:"scope"`
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAddInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if t.Parent != nil {
		if err := t.Parent.Validate(); err != nil {
			return err
		}
	}
	if t.Scope != nil {
		if err := t.Scope.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// UnitAddFromOrgInput ...
type UnitAddFromOrgInput struct {
	UnitAddInput
	OrganizationInput
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAddFromOrgInput) Validate() error {
	if err := t.UnitAddInput.Validate(); err != nil {
		return err
	}
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	return nil
}

// UnitAddFromOUInput ...
type UnitAddFromOUInput struct {
	UnitAddInput
	OrganizationOUInput
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAddFromOUInput) Validate() error {
	if err := t.UnitAddInput.Validate(); err != nil {
		return err
	}
	if err := t.OrganizationOUInput.Validate(); err != nil {
		return err
	}
	return nil
}

// UnitAddFromMembersInput ...
type UnitAddFromMembersInput struct {
	UnitAddInput
	OrganizationInput
	SubjectsInput
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAddFromMembersInput) Validate() error {
	if err := t.UnitAddInput.Validate(); err != nil {
		return err
	}
	if err := t.OrganizationInput.Validate(); err != nil {
		return err
	}
	if err := t.SubjectsInput.Validate(); err != nil {
		return err
	}
	return nil
}

// UnitAddSubjectsInput ...
type UnitAddSubjectsInput struct {
	Target
	SubjectsInput
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAddSubjectsInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if err := t.SubjectsInput.Validate(); err != nil {
		return err
	}
	return nil
}

// UnitAssignParentInput ...
type UnitAssignParentInput struct {
	Target
	Parent Target `json:"parent"`
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAssignParentInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if err := t.Parent.Validate(); err != nil {
		return err
	}
	return nil
}

// UnitAssignScopeInput ...
type UnitAssignScopeInput struct {
	Target
	Scope Target `json:"scope"`
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAssignScopeInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if err := t.Scope.Validate(); err != nil {
		return err
	}
	return nil
}

// UnitAssignObjectInput ...
type UnitAssignObjectInput struct {
	Target
	Object Target `json:"object"`
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAssignObjectInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if err := t.Object.Validate(); err != nil {
		return err
	}
	return nil
}

// UnitAddPermissionsInput ...
type UnitAddPermissionsInput struct {
	Target
	Permissions []PermissionEx `json:"permissions"`
}

// Validate 实现 gear.BodyTemplate
func (t *UnitAddPermissionsInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if len(t.Permissions) == 0 {
		return gear.ErrBadRequest.WithMsgf("permissions empty")
	}
	if len(t.Permissions) > 1000 {
		return gear.ErrBadRequest.WithMsgf("too many permissions: %d", len(t.Permissions))
	}
	cr := make(checkRepetitive)
	for _, p := range t.Permissions {
		if err := cr.Check(p.Permission); err != nil {
			return nil
		}
		if err := p.Validate(); err != nil {
			return err
		}
	}

	return nil
}
