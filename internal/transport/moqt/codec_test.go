package moqt

import "testing"

func TestJSONRPCEncodeDecode(t *testing.T) {
	b, err := EncodeJSONRPC(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "tools/list"})
	if err != nil {
		t.Fatal(err)
	}
	env, err := DecodeJSONRPC(b)
	if err != nil {
		t.Fatal(err)
	}
	if env.Method != "tools/list" {
		t.Fatalf("method=%s", env.Method)
	}
}
func TestToolCallAndResultObjectMapping(t *testing.T) {
	req := ToolCallObject("s1", "nats_info", 7, []byte(`{}`))
	if req.Namespace.String() != "(mcp, s1, tools)" || req.Track != "nats_info" || req.GroupID != 7 || req.ObjectID != 0 {
		t.Fatalf("bad req: %#v", req)
	}
	res := ResultObject("s1", "nats_info", 7, 1, []byte(`{"jsonrpc":"2.0"}`))
	if res.ObjectID != 1 {
		t.Fatal("bad result object")
	}
}
func TestTransportAliases(t *testing.T) {
	for _, n := range []string{"moqt-quic", "quic", "quicmq", "mcp-moqt"} {
		if NormalizeTransport(n) != "moqt-quic" {
			t.Fatalf("alias %s rejected", n)
		}
	}
}
