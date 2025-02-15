package visitorip

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"regexp"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("oss", parseCaddyfile)
}

var paramReplacementPattern = regexp.MustCompile("\\{[a-zA-Z0-9_\\-.]+?}")

// Middleware implements an HTTP handler that replaces response contents based on regex
//
// Additional configuration is required in addition to adding a filter{} block. See
// Github page for instructions.
type Middleware struct {
	// Regex to specify which kind of response should we filter
	//ContentType string `json:"content_type"`
	//// Regex to specify which pattern to look up
	//SearchPattern string `json:"search_pattern"`
	//// A byte-array specifying the string used to replace matches
	//Replacement []byte `json:"replacement"`
	//
	//MaxSize int    `json:"max_size"`
	//Path    string `json:"path"`
	//
	//compiledContentTypeRegex *regexp.Regexp
	//compiledSearchRegex      *regexp.Regexp
	//compiledPathRegex        *regexp.Regexp

	Param1 string `json:"param1,omitempty"`
	Param2 string `json:"param2,omitempty"`
	logger *zap.Logger
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
	m.logger = ctx.Logger(m)
	m.logger.Debug(fmt.Sprintf("Param1: %s. Param2: %s",
		m.Param1,
		m.Param2))
	return nil
}

// Validate implements caddy.Validator.
func (m *Middleware) Validate() error {
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	m.logger.Info("filter middleware called")

	m.logger.Info(m.Param1)
	m.logger.Info(m.Param2)
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if !d.Next() {
		return d.Err("expected token following filter")
	}
	for d.NextBlock(0) {
		key := d.Val()
		var value string
		d.Args(&value)
		if d.NextArg() {
			return d.ArgErr()
		}
		switch key {
		case "param1":
			m.Param1 = value
		case "param2":
			m.Param2 = value

		default:
			return d.Err(fmt.Sprintf("invalid key for filter directive: %s", key))
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
