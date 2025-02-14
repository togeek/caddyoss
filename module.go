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

	//fmt.Println(r.Host)
	fmt.Println(m.Options["param1"])

	fmt.Println("==========================")
	fmt.Println(m.Options["param2"])

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
			fmt.Println("param1")
			m.Options["param1"] = d.RemainingArgs()
		case "param2":
			fmt.Println("param2")
			m.Options["param2"] = d.RemainingArgs()
		default:
			return d.Errf("unrecognized subdirective '%s'", d.Val())

		}
	}

	//for d.Next() {
	//	//if !d.Args(&m.Name) {
	//	//	return d.ArgErr()
	//	//}
	//
	//	// Expect an opening curly brace
	//	if !d.NextBlock(0) {
	//		return d.ArgErr()
	//	}
	//
	//	// Parse key-value pairs inside the curly braces
	//	for d.Next() {
	//		if d.NextArg() {
	//			key := d.Val()
	//			if !d.NextArg() {
	//				return d.ArgErr()
	//			}
	//			valueStr := d.Val()
	//
	//			// Attempt to parse the value as an int, float, or bool
	//			if intVal, err := strconv.Atoi(valueStr); err == nil {
	//				m.Options[key] = intVal
	//			} else if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
	//				m.Options[key] = floatVal
	//			} else if boolVal, err := strconv.ParseBool(valueStr); err == nil {
	//				m.Options[key] = boolVal
	//			} else {
	//				m.Options[key] = valueStr // Treat as string if parsing fails
	//			}
	//		}
	//
	//		if !d.NextLine() {
	//			break
	//		}
	//	}
	//	return nil
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
