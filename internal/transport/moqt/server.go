package moqt

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/server"
	mcpnats "github.com/sinadarbouy/mcp-nats"
	"github.com/sinadarbouy/mcp-nats/internal/security"
)

type Config struct {
	Address, CertFile, KeyFile, ClientCAFile, DPoPJWKSURL, DPoPAudience string
	RequireMTLS, RequireDPoP, AllowInsecureQUIC                         bool
	ChannelBackends                                                     []ChannelBackend
}

type Server struct {
	core      *server.MCPServer
	cfg       Config
	tlsConfig *tls.Config
	verifier  security.DPoPVerifier
	registry  Registry
	listener  net.Listener
	mu        sync.Mutex
	closed    bool
}

func NewServer(core *server.MCPServer, cfg Config) (*Server, error) {
	policy := security.Policy{RequireMTLS: cfg.RequireMTLS, RequireDPoP: cfg.RequireDPoP, AllowInsecureQUIC: cfg.AllowInsecureQUIC, DPoPAudience: cfg.DPoPAudience, DPoPJWKSURL: cfg.DPoPJWKSURL}
	policy.WarnIfInsecure()
	tlsCfg, err := security.BuildTLSConfig(security.TLSConfig{CertFile: cfg.CertFile, KeyFile: cfg.KeyFile, ClientCAFile: cfg.ClientCAFile, RequireMTLS: cfg.RequireMTLS, AllowInsecure: cfg.AllowInsecureQUIC})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrQUICConnectionFailed, err)
	}
	return &Server{core: core, cfg: cfg, tlsConfig: tlsCfg, verifier: policy.DPoPVerifier(), registry: NewRegistry(cfg.ChannelBackends)}, nil
}

func (s *Server) Listen(ctx context.Context) error {
	addr := s.cfg.Address
	if addr == "" {
		addr = "0.0.0.0:9443"
	}
	ln, err := tls.Listen("tcp", addr, s.tlsConfig) // TODO: replace TCP/TLS shim with real QUIC/MOQT library behind the interfaces.
	if err != nil {
		return &TransportError{Category: ErrQUICConnectionFailed, Message: "listen failed", Err: err}
	}
	s.listener = ln
	slog.Warn("Starting experimental MCP-Scalable-Channel MOQT/QUIC transport (TCP/TLS compatibility shim; wire format may change)", "address", addr, "channelBackends", s.registry.Backends())
	errCh := make(chan error, 1)
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				s.mu.Lock()
				closed := s.closed
				s.mu.Unlock()
				if closed {
					errCh <- nil
				} else {
					errCh <- err
				}
				return
			}
			go s.handleConn(ctx, c)
		}
	}()
	select {
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handleConn(ctx context.Context, c net.Conn) {
	defer c.Close()
	sessionID := uuid.NewString()
	state := MCPInitializing
	if tc, ok := c.(*tls.Conn); ok {
		_ = tc.Handshake()
		slog.Info("MOQT session established", "session", sessionID, "peerFingerprint", security.PeerFingerprint(tc.ConnectionState()))
	}
	scanner := bufio.NewScanner(c)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	enc := json.NewEncoder(c)
	baseCtx := mcpnats.ExtractNatsInfoFromEnv(ctx)
	for scanner.Scan() {
		raw := scanner.Bytes()
		env, err := DecodeJSONRPC(raw)
		if err != nil {
			_ = enc.Encode(errorResponse(nil, -32700, err.Error()))
			continue
		}
		if s.cfg.RequireDPoP {
			proof := ""
			var params map[string]any
			_ = json.Unmarshal(env.Params, &params)
			if v, ok := params["dpop"].(string); ok {
				proof = v
			}
			if _, err := s.verifier.Verify(baseCtx, proof, security.DPoPRequestContext{Method: env.Method, URL: "moqt://" + s.cfg.Address, SessionID: sessionID, Audience: s.cfg.DPoPAudience}); err != nil {
				_ = enc.Encode(errorResponse(env.ID, -32001, ErrDPoPInvalid+": "+err.Error()))
				continue
			}
		}
		resp := s.core.HandleMessage(baseCtx, append([]byte(nil), raw...))
		if resp == nil {
			continue
		}
		if env.Method == "initialize" {
			state = MCPActive
			slog.Info("MCP initialized over MOQT", "session", sessionID, "state", state)
		}
		_ = enc.Encode(resp)
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, net.ErrClosed) {
		slog.Warn("MOQT session read failed", "session", sessionID, "error", err)
	}
}

func errorResponse(id any, code int, msg string) map[string]any {
	return map[string]any{"jsonrpc": "2.0", "id": id, "error": map[string]any{"code": code, "message": msg}}
}

func IsTransport(name string) bool {
	switch name {
	case "moqt-quic", "quic", "quicmq", "mcp-moqt":
		return true
	default:
		return false
	}
}
func NormalizeTransport(name string) string {
	if IsTransport(name) {
		return "moqt-quic"
	}
	return name
}
