package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/teambition/gear"
)

// Healthz ..
type Healthz struct {
	blls *bll.Blls
}

// Check ..
func (a *Healthz) Check(ctx *gear.Context) error {
	h, err := a.blls.Models.Model.CheckHealth(ctx)
	if err != nil {
		return err
	}
	return ctx.OkJSON(h)
}
