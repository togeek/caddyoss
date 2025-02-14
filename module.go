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
	httpcaddyfile.RegisterHandlerDirective("visitor_ip", parseCaddyfile)
}

// Middleware implements an HTTP handler that writes the
// visitor's IP address to a file or stream.
type Middleware struct {
	// The file or stream to write to. Can be "stdout"
	// or "stderr".
	//Output string `json:"output,omitempty"`
	Options Options `json:"options,omitempty"`

	//Param1 string `json:"Param1,omitempty"`
	//Param2 string `json:"Param2,omitempty"`
	w io.Writer
}

type Options struct {
	Param1 string `json:"Param1,omitempty"`
	Param2 string `json:"Param2,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.visitor_ip",
		New: func() caddy.Module { return new(Middleware) },
		//New: func() caddy.Module { return &Provider{new(cloudflare.Provider)} },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	//switch m.Output {
	//case "stdout":
	//	m.w = os.Stdout
	//case "stderr":
	//	m.w = os.Stderr
	//default:
	//	return fmt.Errorf("an output stream is required")
	//}

	m.w = os.Stdout

	//fmt.Println("visitor_ip middleware loaded")
	//fmt.Println(m.Output)
	return nil
}

// Validate implements caddy.Validator.
//
//	func (m *Middleware) Validate() error {
//		//if m.w == nil {
//		//	return fmt.Errorf("no writer")
//		//}
//		return nil
//	}
func (m *Middleware) Validate() error {
	//for k, v := range m.Options {
	//	fmt.Printf("Option %s: %v\n", k, v)
	//}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	fmt.Println("visitor_ip middleware loaded")
	//fmt.Println(m.Param1)
	//fmt.Println(m.Param2)

	//if param1, ok := m.Options["param1"]; ok {
	//	fmt.Println("param1:", param1)
	//} else {
	//	fmt.Println("param1 not found or not a string")
	//}
	//
	//if param2, ok := m.Options["param2"]; ok {
	//	fmt.Println("param2:", param2)
	//} else {
	//	fmt.Println("param2 not found or not a string")
	//}

	//if param1, ok := m.Options["param1"].([]string); ok {
	//	fmt.Println("param1:", param1)
	//} else {
	//	fmt.Println("param1 not found or not a []string")
	//}

	return next.ServeHTTP(w, r)
}

func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {

	d.Next() // consume directive name

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "param1":
			if d.NextArg() {
				m.Options.Param1 = d.Val()
			} else {
				return d.ArgErr()
			}
		case "zone_token":
			if d.NextArg() {
				m.Options.Param2 = d.Val()
			} else {
				return d.ArgErr()
			}
		default:
			return d.Errf("unrecognized subdirective '%s'", d.Val())
		}
	}
	if d.NextArg() {
		return d.Errf("unexpected argument '%s'", d.Val())
	}
	//if p.Provider.APIToken == "" {
	//	return d.Err("missing API token")
	//}
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
