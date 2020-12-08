package tpl

import "github.com/teambition/gear"

// BatchAddPermissionsInput ...
type BatchAddPermissionsInput struct {
	Permissions []string `json:"permissions"`
}

// Validate 实现 gear.BodyTemplate
func (t *BatchAddPermissionsInput) Validate() error {
	// OTID UnmarshalText method will validate
	if len(t.Permissions) == 0 {
		return gear.ErrBadRequest.WithMsg("empty permissions")
	}
	for _, permission := range t.Permissions {
		if err := CheckPermission(permission); err != nil {
			return err
		}
	}
	return nil
}

// ListPermissionsInput ...
type ListPermissionsInput struct {
	Pagination
	ResourcesInput
}

// Validate 实现 gear.BodyTemplate
func (t *ListPermissionsInput) Validate() error {
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	if err := t.ResourcesInput.Validate(); err != nil {
		return err
	}
	return nil
}

// DeletePermissionInput ...
type DeletePermissionInput struct {
	Permission string `json:"permission"`
}

// Validate 实现 gear.BodyTemplate
func (t *DeletePermissionInput) Validate() error {
	if err := CheckPermission(t.Permission); err != nil {
		return err
	}
	return nil
}
