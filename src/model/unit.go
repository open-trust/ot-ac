package model

import (
	"context"

	"github.com/open-trust/ot-ac/src/schema"
	"github.com/open-trust/ot-ac/src/tpl"
)

// Unit ...
type Unit struct {
	*Model
}

// Add ...
func (m *Unit) Add(ctx context.Context, tenant *schema.Tenant, inputs []*tpl.AddUnitInput) error {
	return nil
}
