package tpl

import (
	"github.com/teambition/gear"
)

// Permission ...
type Permission struct {
	UID        string `json:"uid,omitempty"`
	Permission string `json:"permission"`
}

// GetPermissionUID ...
func GetPermissionUID(ps []Permission, p string) string {
	for _, v := range ps {
		if v.Permission == p {
			return v.UID
		}
	}
	return ""
}

// PermissionEx ...
type PermissionEx struct {
	Permission string     `json:"permission"`
	Extensions Extensions `json:"extensions"`
}

// Validate 实现 gear.BodyTemplate
func (t *PermissionEx) Validate() error {
	if err := CheckPermission(t.Permission); err != nil {
		return err
	}
	if len(t.Extensions) > 10 {
		return gear.ErrBadRequest.WithMsgf("too many extensions: %d", len(t.Extensions))
	}
	if err := t.Extensions.Validate(); err != nil {
		return err
	}
	return nil
}

// PermissionBatchAddInput ...
type PermissionBatchAddInput struct {
	Permissions []string `json:"permissions"`
}

// Validate 实现 gear.BodyTemplate
func (t *PermissionBatchAddInput) Validate() error {
	// OTID UnmarshalText method will validate
	if len(t.Permissions) == 0 {
		return gear.ErrBadRequest.WithMsg("empty permissions")
	}
	cr := make(checkRepetitive)
	for _, permission := range t.Permissions {
		if err := cr.Check(permission); err != nil {
			return err
		}
		if err := CheckPermission(permission); err != nil {
			return err
		}
	}
	return nil
}

// PermissionListInput ...
type PermissionListInput struct {
	Pagination
	ResourcesInput
}

// Validate 实现 gear.BodyTemplate
func (t *PermissionListInput) Validate() error {
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	if err := t.ResourcesInput.Validate(); err != nil {
		return err
	}
	return nil
}

// PermissionDeleteInput ...
type PermissionDeleteInput struct {
	Permission string `json:"permission"`
}

// Validate 实现 gear.BodyTemplate
func (t *PermissionDeleteInput) Validate() error {
	if err := CheckPermission(t.Permission); err != nil {
		return err
	}
	return nil
}
