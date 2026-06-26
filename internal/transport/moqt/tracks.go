package moqt

import "strings"

type TrackNamespace []string
type TrackName string
type EventName string

const (
	TrackControl       = "control"
	TrackResources     = "resources"
	TrackTools         = "tools"
	TrackPrompts       = "prompts"
	TrackNotifications = "notifications"
	TrackElicitation   = "elicitation"
	TrackLogs          = "logs"
)

const (
	ControlChannelClientToServer = "client-to-server"
	ControlChannelServerToClient = "server-to-client"
)

const (
	EventInitialize            EventName = "initialize"
	EventInitialized           EventName = "initialized"
	EventPing                  EventName = "ping"
	EventPong                  EventName = "pong"
	EventCapabilityNegotiation EventName = "capability.negotiation"
	EventSessionTeardown       EventName = "session.teardown"
	EventToolCancellation      EventName = "tool.cancellation"
	EventTransportError        EventName = "transport.error"
	EventToolsList             EventName = "tools.list"
	EventToolsCall             EventName = "tools.call"
	EventToolsProgress         EventName = "tools.progress"
	EventToolsResult           EventName = "tools.result"
	EventToolsError            EventName = "tools.error"
	EventResourcesRead         EventName = "resources.read"
	EventPromptsGet            EventName = "prompts.get"
	EventNotification          EventName = "notification"
)

var ExperimentalTracks = []string{TrackControl, TrackTools, TrackResources, TrackPrompts, TrackNotifications, TrackElicitation, TrackLogs}

var ControlEvents = []EventName{EventInitialize, EventInitialized, EventPing, EventPong, EventCapabilityNegotiation, EventSessionTeardown, EventToolCancellation, EventTransportError}

func Namespace(sessionID, track string) TrackNamespace {
	return TrackNamespace{"mcp", sessionID, track}
}
func ControlNamespace(sessionID string) TrackNamespace { return Namespace(sessionID, TrackControl) }
func ToolsNamespace(sessionID string) TrackNamespace   { return Namespace(sessionID, TrackTools) }
func (n TrackNamespace) String() string                { return "(" + strings.Join(n, ", ") + ")" }

// ControlTrack names one of the two concurrently-open logical control channels.
// Both control channels use the same EventName vocabulary and may carry the same
// JSON-RPC payload; the client decides which channel(s) it subscribes to or uses.
func ControlTrack(channel string) TrackName { return TrackName(channel) }
func ClientControlTrack() TrackName         { return ControlTrack(ControlChannelClientToServer) }
func ServerControlTrack() TrackName         { return ControlTrack(ControlChannelServerToClient) }
func ToolTrack(toolName string) TrackName   { return TrackName(toolName) }
