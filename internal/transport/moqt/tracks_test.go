package moqt

import "testing"

func TestTrackNamespaceGeneration(t *testing.T) {
	ns := ToolsNamespace("s1")
	if got, want := ns.String(), "(mcp, s1, tools)"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if ControlTrack("client-to-server") != "client-to-server" {
		t.Fatal("bad control track")
	}
}
