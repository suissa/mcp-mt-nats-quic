package moqt

const (
	ErrQUICConnectionFailed = "transport.quic.connection_failed"
	ErrMOQTSessionFailed    = "transport.moqt.session_failed"
	ErrTrackSubscribeFailed = "transport.moqt.track_subscribe_failed"
	ErrPublishFailed        = "transport.moqt.publish_failed"
	ErrMCPInitializeFailed  = "mcp.initialize_failed"
	ErrMCPToolNotFound      = "mcp.tool_not_found"
	ErrMCPToolExecution     = "mcp.tool_execution_failed"
	ErrMTLSRequired         = "security.mtls_required"
	ErrDPoPRequired         = "security.dpop_required"
	ErrDPoPInvalid          = "security.dpop_invalid"
	ErrUnauthorized         = "security.unauthorized"
)

type TransportError struct {
	Category, Message string
	Err               error
}

func (e *TransportError) Error() string {
	if e.Err != nil {
		return e.Category + ": " + e.Message + ": " + e.Err.Error()
	}
	return e.Category + ": " + e.Message
}
func (e *TransportError) Unwrap() error { return e.Err }
