package pkg

import (
	gologmid "github.com/driftappdev/libpackage/logmid/logging-middleware"
	"net/http"
)

type LogMiddlewareConfig = gologmid.Config
type HTTPLogField = gologmid.Field
type HTTPLogger = gologmid.Logger

func NewHTTPLogMiddleware(cfg LogMiddlewareConfig) func(http.Handler) http.Handler {
	return gologmid.New(cfg)
}
func NewHTTPRecovery(logger HTTPLogger, extraFields ...HTTPLogField) func(http.Handler) http.Handler {
	return gologmid.Recovery(logger, extraFields...)
}
