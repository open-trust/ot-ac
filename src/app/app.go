package app

import (
	"log"
	"strings"
	"time"

	"github.com/teambition/compressible-go"
	"github.com/teambition/gear"

	"github.com/open-trust/ot-ac/src/conf"
	"github.com/open-trust/ot-ac/src/logging"
	"github.com/open-trust/ot-ac/src/util"
)

// New ...
func New() *gear.App {
	app := gear.New()

	app.Set(gear.SetEnv, conf.AppEnv)
	app.Set(gear.SetTrustedProxy, conf.Config.TrustedProxy)
	app.Set(gear.SetBodyParser, gear.DefaultBodyParser(2<<20)) // 2MB
	// ignore TLS handshake error
	app.Set(gear.SetLogger, log.New(gear.DefaultFilterWriter(), "", 0))
	app.Set(gear.SetCompress, compressible.WithThreshold(1024))
	app.Set(gear.SetTimeout, time.Second*5)
	app.Set(gear.SetRenderError, gear.RenderErrorResponse)
	app.Set(gear.SetParseError, func(err error) gear.HTTPError {
		msg := err.Error()
		if strings.Contains(msg, "already exists") {
			return gear.ErrConflict.WithMsg(msg)
		}
		return gear.ParseError(err)
	})
	if app.Env() != "testing" {
		app.Use(logging.WithAccessLogger)
	}

	err := util.DigInvoke(func(routers []*gear.Router) error {
		for _, router := range routers {
			app.UseHandler(router)
		}
		return nil
	})

	if err != nil {
		logging.Panicf("DigInvoke error: %v", err)
	}

	return app
}
