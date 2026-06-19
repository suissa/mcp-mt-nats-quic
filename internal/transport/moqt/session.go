package moqt

import "context"

type SessionState string

const (
	Disconnected           SessionState = "DISCONNECTED"
	QUICConnected                       = "QUIC_CONNECTED"
	MOQTEstablished                     = "MOQT_ESTABLISHED"
	TracksSubscribed                    = "TRACKS_SUBSCRIBED"
	MCPInitializing                     = "MCP_INITIALIZING"
	CapabilitiesNegotiated              = "CAPABILITIES_NEGOTIATED"
	MCPActive                           = "MCP_ACTIVE"
	SessionTerminating                  = "SESSION_TERMINATING"
)

type MOQTPublisher interface {
	PublishObject(ctx context.Context, namespace TrackNamespace, track TrackName, groupID uint64, objectID uint64, payload []byte) error
}
type MOQTSubscriber interface {
	Subscribe(ctx context.Context, namespace TrackNamespace, track TrackName) (<-chan MOQTObject, error)
}
type MOQTSession interface {
	ID() string
	Publisher() MOQTPublisher
	Subscriber() MOQTSubscriber
	Close(ctx context.Context) error
}
