package tpl

import (
	"github.com/open-trust/ot-ac/src/conf"
	otgo "github.com/open-trust/ot-go-lib"
	"github.com/teambition/gear"
)

// AddTenantInput ...
type AddTenantInput struct {
	OTID   otgo.OTID `json:"tenant"`
	Status int       `json:"status"`
}

// Validate 实现 gear.BodyTemplate
func (t *AddTenantInput) Validate() error {
	// OTID UnmarshalText method will validate
	if !t.OTID.MemberOf(conf.OT.TrustDomain) {
		return gear.ErrBadRequest.WithMsg("tenant %s is not a member of %s", t.OTID.String(), conf.OT.TrustDomain.String())
	}
	if t.Status < -1 {
		return gear.ErrBadRequest.WithMsgf("invalid tenant status %d", t.Status)
	}
	return nil
}
