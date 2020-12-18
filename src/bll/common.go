package bll

import (
	"github.com/open-trust/ot-ac/src/model"
	"github.com/open-trust/ot-ac/src/util"
	"github.com/teambition/gear"
)

func init() {
	util.DigProvide(NewBlls)
}

// Blls ...
type Blls struct {
	Models       *model.Models
	AC           *AC
	Admin        *Admin
	Object       *Object
	Organization *Organization
	Permission   *Permission
	Scope        *Scope
	Unit         *Unit
}

// NewBlls ...
func NewBlls(models *model.Models) *Blls {
	return &Blls{
		Models:       models,
		AC:           &AC{models},
		Admin:        &Admin{models},
		Object:       &Object{models},
		Organization: &Organization{models},
		Permission:   &Permission{models},
		Scope:        &Scope{models},
		Unit:         &Unit{models},
	}
}

type contextKey int

const (
	bllsKey contextKey = iota
)

// IntoCtx ...
func IntoCtx(blls *Blls) func(ctx *gear.Context) error {
	return func(ctx *gear.Context) error {
		ctx.SetAny(bllsKey, blls)
		return nil
	}
}

// FromCtx ...
func FromCtx(ctx *gear.Context) (*Blls, error) {
	val, err := ctx.Any(bllsKey)
	if err != nil {
		return nil, err
	}
	blls, ok := val.(*Blls)
	if !ok {
		return nil, gear.ErrInternalServerError.WithMsg("Blls not exist")
	}
	return blls, nil
}
