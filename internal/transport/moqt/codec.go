package moqt

import (
	"encoding/json"
	"fmt"
)

type JSONRPCEnvelope struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   any             `json:"error,omitempty"`
}

type MOQTObject struct {
	Namespace TrackNamespace
	Track     TrackName
	Event     EventName
	GroupID   uint64
	ObjectID  uint64
	Payload   []byte
}

func EncodeJSONRPC(v any) ([]byte, error) { return json.Marshal(v) }
func DecodeJSONRPC(payload []byte) (*JSONRPCEnvelope, error) {
	var env JSONRPCEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		return nil, fmt.Errorf("decode JSON-RPC: %w", err)
	}
	if env.JSONRPC != "2.0" {
		return nil, fmt.Errorf("invalid JSON-RPC version %q", env.JSONRPC)
	}
	return &env, nil
}

func ControlObject(sessionID, channel string, event EventName, groupID, objectID uint64, payload []byte) MOQTObject {
	return MOQTObject{Namespace: ControlNamespace(sessionID), Track: ControlTrack(channel), Event: event, GroupID: groupID, ObjectID: objectID, Payload: payload}
}

func ToolCallObject(sessionID, toolName string, invocationID uint64, payload []byte) MOQTObject {
	return MOQTObject{Namespace: ToolsNamespace(sessionID), Track: ToolTrack(toolName), Event: EventToolsCall, GroupID: invocationID, ObjectID: 0, Payload: payload}
}

func ResultObject(sessionID, toolName string, invocationID, objectID uint64, payload []byte) MOQTObject {
	return MOQTObject{Namespace: ToolsNamespace(sessionID), Track: ToolTrack(toolName), Event: EventToolsResult, GroupID: invocationID, ObjectID: objectID, Payload: payload}
}
