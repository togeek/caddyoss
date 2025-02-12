package caddyoss

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"net/http"
)

func init() {
	caddy.RegisterModule(MyCustomModule{})
}

type MyCustomModule struct{}

func (MyCustomModule) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.oss",
		New: func() caddy.Module { return new(MyCustomModule) },
	}
}

func (m *MyCustomModule) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	fmt.Fprintln(w, "Hello from oss!")
	return next.ServeHTTP(w, r)
}
