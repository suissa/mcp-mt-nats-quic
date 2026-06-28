package moqt

import "testing"

func TestTrackNamespaceGeneration(t *testing.T) {
	ns := ToolsNamespace("s1")
	if got, want := ns.String(), "(mcp, s1, tools)"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if ClientControlTrack() != "client-to-server" || ServerControlTrack() != "server-to-client" {
		t.Fatal("bad control track names")
	}
}

func TestControlChannelsShareEventVocabulary(t *testing.T) {
	payload := []byte(`{"jsonrpc":"2.0","method":"initialize"}`)
	clientObj := ControlObject("s1", ControlChannelClientToServer, EventInitialize, 1, 0, payload)
	serverObj := ControlObject("s1", ControlChannelServerToClient, EventInitialize, 1, 0, payload)

	if clientObj.Track == serverObj.Track {
		t.Fatal("control channels should remain distinct logical tracks")
	}
	if clientObj.Event != serverObj.Event {
		t.Fatalf("events should use same nomenclature: %q != %q", clientObj.Event, serverObj.Event)
	}
	if string(clientObj.Payload) != string(serverObj.Payload) {
		t.Fatal("both channels must be able to carry the same payload")
	}
}
