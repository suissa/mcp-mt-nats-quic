package security

import "log/slog"

type Policy struct {
	RequireMTLS       bool
	RequireDPoP       bool
	AllowInsecureQUIC bool
	DPoPAudience      string
	DPoPJWKSURL       string
}

func (p Policy) DPoPVerifier() DPoPVerifier {
	if !p.RequireDPoP {
		return NoopDPoPVerifier{}
	}
	return NewStrictDPoPVerifier(0)
}
func (p Policy) WarnIfInsecure() {
	if !p.RequireMTLS && !p.RequireDPoP {
		slog.Warn("MCP-Scalable-Channel security disabled: both mTLS and DPoP are off; use only for local development")
	}
}
