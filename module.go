package visitorip

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"io"
	"net/http"
	"os"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("oss", parseCaddyfile)
}

// Middleware implements an HTTP handler that writes the
// visitor's IP address to a file or stream.
type Middleware struct {
	// The file or stream to write to. Can be "stdout"
	// or "stderr".
	Output string `json:"stdout,omitempty"`
	Param3 string `json:"param3,omitempty"`
	Param1 string `json:"param1,omitempty"`
	Param2 string `json:"param2,omitempty"`
	w      io.Writer
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.oss",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	switch m.Output {
	case "stdout":
		m.w = os.Stdout
	case "stderr":
		m.w = os.Stderr
	default:
		return fmt.Errorf("an output stream is required")
	}
	return nil
}

// Validate implements caddy.Validator.
func (m *Middleware) Validate() error {
	if m.w == nil {
		return fmt.Errorf("no writer")
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	_, _ = m.w.Write([]byte(r.RemoteAddr + "\n"))

	//fmt.Println(r.Host)
	fmt.Println(m.Output)
	fmt.Println(m.Param1)
	fmt.Println(m.Param2)
	fmt.Println(m.Param3)

	_, _ = m.w.Write([]byte(r.URL.Host + "\n"))

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	//d.Next() // consume directive name

	// require an argument
	//if !d.NextArg() {
	//	return d.ArgErr()
	//}

	//fmt.Println(d.Val())
	//
	//// store the argument
	//m.Output = d.Val()
	//d.Next()
	//m.Param1 = d.Val()
	//d.Next()
	//m.Param2 = d.Val()

	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "param1":
				if !d.NextArg() {
					return d.ArgErr()
				}
				m.Param1 = d.Val()
			case "param2":
				if !d.NextArg() {
					return d.ArgErr()
				}
				m.Param2 = d.Val()
			case "param3":
				if !d.NextArg() {
					return d.ArgErr()
				}
				m.Param3 = d.Val()
			default:
				return d.Errf("unknown directive: %s", d.Val())
			}
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)
