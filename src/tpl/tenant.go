package tpl

import (
	"github.com/open-trust/ot-ac/src/conf"
	otgo "github.com/open-trust/ot-go-lib"
	"github.com/teambition/gear"
)

// Tenant ...
type Tenant struct {
	UID    string `json:"uid,omitempty"`
	Tenant string `json:"tenant"`
	Status int    `json:"status"`
}

// TenantAddInput ...
type TenantAddInput struct {
	Tenant otgo.OTID `json:"tenant"`
	Status int       `json:"status"`
}

// Validate 实现 gear.BodyTemplate
func (t *TenantAddInput) Validate() error {
	// OTID UnmarshalText method will validate
	if !t.Tenant.MemberOf(conf.OT.TrustDomain) {
		return gear.ErrBadRequest.WithMsg("tenant %s is not a member of %s", t.Tenant.String(), conf.OT.TrustDomain.String())
	}
	if t.Status < -1 {
		return gear.ErrBadRequest.WithMsgf("invalid tenant status %d", t.Status)
	}
	return nil
}
