package api

import (
	"strconv"

	"github.com/open-trust/ot-ac/src/bll"
	"github.com/teambition/gear"
)

// AC ..
type AC struct {
	blls *bll.Blls
}

// Serve ..
func (a *AC) Serve(ctx *gear.Context) error {
	action := ctx.Param("action")
	switch action {
	case "CheckUnit":
		return a.CheckUnit(ctx)
	case "CheckScope":
		return a.CheckScope(ctx)
	case "CheckObject":
		return a.CheckObject(ctx)
	case "ListPermissionsByUnit":
		return a.ListPermissionsByUnit(ctx)
	case "ListPermissionsByScope":
		return a.ListPermissionsByScope(ctx)
	case "ListPermissionsByObject":
		return a.ListPermissionsByObject(ctx)
	case "ListObject":
		return a.ListObject(ctx)
	case "SearchObject":
		return a.SearchObject(ctx)
	}
	return gear.ErrBadRequest.WithMsgf("unknown action %s", strconv.Quote(action))
}

// CheckUnit 检查请求主体到指定管理单元有没有指定权限
func (a *AC) CheckUnit(ctx *gear.Context) error {
	return nil
}

func (a *AC) CheckScope(ctx *gear.Context) error {
	return nil
}

func (a *AC) CheckObject(ctx *gear.Context) error {
	return nil
}

func (a *AC) ListPermissionsByUnit(ctx *gear.Context) error {
	return nil
}

func (a *AC) ListPermissionsByScope(ctx *gear.Context) error {
	return nil
}

func (a *AC) ListPermissionsByObject(ctx *gear.Context) error {
	return nil
}

func (a *AC) ListObject(ctx *gear.Context) error {
	return nil
}

func (a *AC) SearchObject(ctx *gear.Context) error {
	return nil
}
