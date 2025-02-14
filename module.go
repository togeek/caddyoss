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
	Options map[string]interface{} `json:"options,omitempty"`

	w io.Writer
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.visitor_ip",
		New: func() caddy.Module { return new(Middleware) },
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
	for k, v := range m.Options {
		switch v.(type) {
		case string, int, bool, float64: // 允许的类型
			// 可以根据 key 做更细致的验证
			fmt.Printf("Option %s: %v\n", k, v)
		default:
			return fmt.Errorf("option %s has invalid type: %T", k, v)
		}
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	_, _ = m.w.Write([]byte(r.RemoteAddr + "\n"))

	if param1, ok := m.Options["param1"].(string); ok {
		fmt.Println("param1:", param1)
	} else {
		fmt.Println("param1 not found or not a string")
	}

	if param1, ok := m.Options["param1"].([]string); ok {
		fmt.Println("param1:", param1)
	} else {
		fmt.Println("param1 not found or not a []string")
	}

	_, _ = m.w.Write([]byte(r.RemoteAddr + "\n"))

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
//
//	func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
//		d.Next() // consume directive name
//
//		// require an argument
//		if !d.NextArg() {
//			return d.ArgErr()
//		}
//
//		fmt.Println(d.Val())
//
//		// store the argument
//		m.Output = d.Val()
//		return nil
//	}
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	m.Options = make(map[string]interface{})

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "param1":
			m.Options["param1"] = d.RemainingArgs()
			fmt.Println("param1")

		default:
			m.Options["param1"] = "sssss"
			fmt.Println("ssdddwsd3232")
			return d.Errf("unrecognized subdirective '%s'", d.Val())

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
