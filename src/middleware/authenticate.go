package middleware

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/conf"
	"github.com/open-trust/ot-ac/src/logging"
	"github.com/open-trust/ot-ac/src/tpl"
	otgo "github.com/open-trust/ot-go-lib"
	"github.com/teambition/gear"
)

type contextKey int

const (
	authKey contextKey = iota
	tenantKey
)

// VerifyAdmin 验证请求者管理员身份，如果验证失败，则返回 401 的 gear.HTTPError
func VerifyAdmin(ctx *gear.Context) error {
	token := otgo.ExtractTokenFromHeader(ctx.Req.Header)
	if token == "" {
		return gear.ErrUnauthorized.WithMsg("invalid authorization token")
	}

	vid, err := conf.OT.ParseOTVID(ctx, token)
	if err != nil {
		return gear.ErrUnauthorized.From(err)
	}
	if !vid.ID.Equal(conf.OT.OTID) {
		return gear.ErrForbidden.WithMsgf("%s is not admin", vid.ID.String())
	}

	ctx.SetAny(authKey, vid)
	logging.AccessLogger.SetTo(ctx, "subject", vid.ID.String())
	return nil
}

// VerifyTenant 验证请求者租户身份，如果验证失败，则返回 401 的 gear.HTTPError
func VerifyTenant(ctx *gear.Context) error {
	token := otgo.ExtractTokenFromHeader(ctx.Req.Header)
	if token == "" {
		return gear.ErrUnauthorized.WithMsg("invalid authorization token")
	}

	vid, err := conf.OT.ParseOTVID(ctx, token)
	if err != nil {
		return gear.ErrUnauthorized.From(err)
	}

	ctx.SetAny(authKey, vid)
	logging.AccessLogger.SetTo(ctx, "subject", vid.ID.String())

	blls, err := bll.FromCtx(ctx)
	if err != nil {
		return err
	}
	tenant, err := blls.Models.Tenant.Get(ctx, vid.ID)
	if err != nil {
		return gear.ErrUnauthorized.WithMsgf("invalid tenant: %s", err.Error())
	}
	logging.AccessLogger.SetTo(ctx, "tenant", tenant.Tenant)
	if tenant.Status < 0 {
		return gear.ErrForbidden.WithMsgf("%s is forbidden", vid.ID.String())
	}
	ctx.SetAny(tenantKey, tenant)
	return nil
}

// VidFromCtx ...
func VidFromCtx(ctx *gear.Context) (*otgo.OTVID, error) {
	val, err := ctx.Any(authKey)
	if err != nil {
		return nil, err
	}
	vid, ok := val.(*otgo.OTVID)
	if !ok {
		return nil, gear.ErrUnauthorized.WithMsg("OTVID not exist")
	}
	return vid, nil
}

// TenantFromCtx ...
func TenantFromCtx(ctx *gear.Context) (*tpl.Tenant, error) {
	val, err := ctx.Any(tenantKey)
	if err != nil {
		return nil, err
	}
	tenant, ok := val.(*tpl.Tenant)
	if !ok {
		return nil, gear.ErrUnauthorized.WithMsg("tenant not exist")
	}
	return tenant, nil
}
