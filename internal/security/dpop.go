package security

import (
	"context"
	"errors"
	"sync"
	"time"
)

type DPoPRequestContext struct {
	Method, URL, SessionID, Audience string
	Now                              time.Time
}
type DPoPClaims struct {
	HTM, HTU, JTI, SessionID string
	IAT                      time.Time
}
type DPoPVerifier interface {
	Verify(ctx context.Context, proof string, req DPoPRequestContext) (*DPoPClaims, error)
}

type NoopDPoPVerifier struct{}

func (NoopDPoPVerifier) Verify(ctx context.Context, proof string, req DPoPRequestContext) (*DPoPClaims, error) {
	return &DPoPClaims{HTM: req.Method, HTU: req.URL, SessionID: req.SessionID, IAT: time.Now()}, nil
}

type StrictDPoPVerifier struct {
	RequireProof bool
	Window       time.Duration
	seen         map[string]time.Time
	mu           sync.Mutex
}

func NewStrictDPoPVerifier(window time.Duration) *StrictDPoPVerifier {
	if window == 0 {
		window = 5 * time.Minute
	}
	return &StrictDPoPVerifier{RequireProof: true, Window: window, seen: map[string]time.Time{}}
}
func (v *StrictDPoPVerifier) Verify(ctx context.Context, proof string, req DPoPRequestContext) (*DPoPClaims, error) {
	if proof == "" {
		return nil, errors.New("DPoP proof required")
	}
	// TODO: verify JWT signature via JWKS/static key and validate htu/htm/iat claims.
	now := req.Now
	if now.IsZero() {
		now = time.Now()
	}
	jti := proof
	v.mu.Lock()
	defer v.mu.Unlock()
	for k, t := range v.seen {
		if now.Sub(t) > v.Window {
			delete(v.seen, k)
		}
	}
	if _, ok := v.seen[jti]; ok {
		return nil, errors.New("DPoP proof replay detected")
	}
	v.seen[jti] = now
	return &DPoPClaims{HTM: req.Method, HTU: req.URL, JTI: jti, SessionID: req.SessionID, IAT: now}, nil
}
