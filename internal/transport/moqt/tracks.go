package moqt

import "strings"

type TrackNamespace []string
type TrackName string

const (
	TrackControl       = "control"
	TrackResources     = "resources"
	TrackTools         = "tools"
	TrackPrompts       = "prompts"
	TrackNotifications = "notifications"
	TrackElicitation   = "elicitation"
	TrackLogs          = "logs"
)

var ExperimentalTracks = []string{TrackControl, TrackTools, TrackResources, TrackPrompts, TrackNotifications, TrackElicitation, TrackLogs}

func Namespace(sessionID, track string) TrackNamespace {
	return TrackNamespace{"mcp", sessionID, track}
}
func ControlNamespace(sessionID string) TrackNamespace { return Namespace(sessionID, TrackControl) }
func ToolsNamespace(sessionID string) TrackNamespace   { return Namespace(sessionID, TrackTools) }
func (n TrackNamespace) String() string                { return "(" + strings.Join(n, ", ") + ")" }

func ControlTrack(direction string) TrackName { return TrackName(direction) }
func ToolTrack(toolName string) TrackName     { return TrackName(toolName) }
