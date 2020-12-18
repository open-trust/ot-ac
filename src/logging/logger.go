package logging

import (
	"context"
	"fmt"
	"os"

	"github.com/open-trust/ot-ac/src/conf"
	"github.com/teambition/gear"
	gearLogging "github.com/teambition/gear/logging"
)

func init() {
	Logger.SetJSONLog()
	AccessLogger.SetJSONLog()

	// AccessLogger is not needed to set level.
	err := gearLogging.SetLoggerLevel(Logger, conf.Config.Logger.Level)
	if err != nil {
		Logger.Err(err)
	}
}

type contextKey int

const (
	logKey contextKey = iota
)

// AccessLogger is used for access log
var AccessLogger = gearLogging.New(os.Stdout)

// Logger is used for the server.
var Logger = gearLogging.New(os.Stderr)

// SrvLog returns a Log with kind of server.
func SrvLog(format string, args ...interface{}) gearLogging.Log {
	return gearLogging.Log{
		"kind":    "server",
		"message": fmt.Sprintf(format, args...),
	}
}

// Panicf produce a "Emergency" log into the Logger.
func Panicf(format string, args ...interface{}) {
	Logger.Panic(SrvLog(format, args...))
}

// Errf produce a "Error" log into the Logger.
func Errf(format string, args ...interface{}) {
	Logger.Err(SrvLog(format, args...))
}

// Warningf produce a "Warning" log into the Logger.
func Warningf(format string, args ...interface{}) {
	Logger.Warning(SrvLog(format, args...))
}

// Infof produce a "Informational" log into the Logger.
func Infof(format string, args ...interface{}) {
	Logger.Info(SrvLog(format, args...))
}

// Debugf produce a "Debug" log into the Logger.
func Debugf(format string, args ...interface{}) {
	Logger.Debug(SrvLog(format, args...))
}

// FromCtx retrieve the Log instance for the AccessLogger.
func FromCtx(ctx *gear.Context) gearLogging.Log {
	return AccessLogger.FromCtx(ctx)
}

// WithAccessLogger ...
func WithAccessLogger(ctx *gear.Context) error {
	if err := AccessLogger.Serve(ctx); err != nil {
		return err
	}
	log := AccessLogger.FromCtx(ctx)
	ctx.WithContext(context.WithValue(ctx.Context(), logKey, log))
	return nil
}

// SetKV ...
func SetKV(ctx context.Context, key string, val interface{}) {
	if v := ctx.Value(logKey); v != nil {
		v.(gearLogging.Log)[key] = val
	}
}
