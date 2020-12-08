package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/util"
)

func init() {
	util.DigProvide(NewAPIs)
}

// APIs ..
type APIs struct {
	Blls       *bll.Blls
	AC         *AC
	Admin      *Admin
	Healthz    *Healthz
	Object     *Object
	Permission *Permission
	Scope      *Scope
	Unit       *Unit
}

// NewAPIs ...
func NewAPIs(blls *bll.Blls) *APIs {
	return &APIs{
		Blls:       blls,
		AC:         &AC{blls: blls},
		Admin:      &Admin{blls: blls},
		Healthz:    &Healthz{blls: blls},
		Object:     &Object{blls: blls},
		Permission: &Permission{blls: blls},
		Scope:      &Scope{blls: blls},
		Unit:       &Unit{blls: blls},
	}
}
