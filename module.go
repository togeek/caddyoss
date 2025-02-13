package visitorip

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"io"
	"net/http"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("oss_sign", parseCaddyfile)
}

// Middleware implements an HTTP handler that writes the
// visitor's IP address to a file or stream.
type Middleware struct {
	// The file or stream to write to. Can be "stdout"
	// or "stderr".
	//Parameters map[string]string

	Param1 string `json:"param1,omitempty"`
	Param2 string `json:"param2,omitempty"`
	Param3 string `json:"param3,omitempty"`
	//OSSConfig
	w io.Writer
}

//type OSSConfig struct {
//	AccessKeyID     string `json:"access_key_id,omitempty"`
//	AccessKeySecret string `json:"access_key_secret,omitempty"`
//	BucketName      string `json:"bucket_name,omitempty"`
//	Endpoint        string `json:"endpoint,omitempty"`
//	SignExpiry      int64  `json:"sign_expiry,omitempty"`
//}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.oss_sign",
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

	return nil
}

// Validate implements caddy.Validator.
func (m *Middleware) Validate() error {
	//if m.w == nil {
	//	return fmt.Errorf("no writer")
	//}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	_, _ = m.w.Write([]byte(r.RemoteAddr + "\n"))

	fmt.Println(r.Host)
	//fmt.Println(m.OSSConfig.AccessKeyID)
	//fmt.Println(m.Param1)
	//fmt.Println(m.Param2)
	//fmt.Println(m.Param3)

	_, _ = m.w.Write([]byte(r.URL.Host + "\n"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Param1: %s\n", m.Param1)
	fmt.Fprintf(w, "Param2: %d\n", m.Param2)
	fmt.Fprintf(w, "Param3: %t\n", m.Param3)

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	// 循环读取指令的参数
	for d.Next() {
		// 确保只有一个指令名称
		if !d.Args(&m.Param1, &m.Param2, &m.Param3) {
			return d.ArgErr() // 返回参数错误
		}
		m.Param2 = d.Val()

		// 转换参数类型 (如果需要)
		// 例如，假设 Param2 是一个整数：
		// if err != nil {
		//     return fmt.Errorf("invalid Param2: %v", err)
		// }

		// 例如，假设 Param3 是一个bool：
		// if err != nil {
		//     return fmt.Errorf("invalid Param3: %v", err)
		// }

		// 可以选择性地检查是否有剩余的参数（如果指令只能接受特定数量的参数）
		//if d.NextArg() {
		//	return d.ArgErr() // 返回额外的参数错误
		//}

		return nil // 解析成功
	}
	return fmt.Errorf("expected at least one argument")
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
