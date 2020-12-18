package tpl

// ACCheckPermissions ...
type ACCheckPermissions struct {
	PermissionBatchAddInput
	Subject          string `json:"subject"`
	WithOrganization bool   `json:"withOrganization"`
	IgnoreScope      bool   `json:"ignoreScope"` // 仅对 Object 权限检查有效
}

// Validate 实现 gear.BodyTemplate
func (t *ACCheckPermissions) Validate() error {
	if err := CheckSubject(t.Subject); err != nil {
		return err
	}
	if err := t.PermissionBatchAddInput.Validate(); err != nil {
		return err
	}
	return nil
}

// ACCheckPermissionsInput ...
type ACCheckPermissionsInput struct {
	Target
	ACCheckPermissions
}

// Validate 实现 gear.BodyTemplate
func (t *ACCheckPermissionsInput) Validate() error {
	if err := t.Target.Validate(); err != nil {
		return err
	}
	if err := t.ACCheckPermissions.Validate(); err != nil {
		return err
	}
	return nil
}

// ACPermissionPayload ...
type ACPermissionPayload struct {
	Target
	PermissionEx
}
