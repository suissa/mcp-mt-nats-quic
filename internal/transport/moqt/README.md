# Experimental MCP-Scalable-Channel MOQT/QUIC Transport

This package contains an experimental MCP-over-MOQT/QUIC transport adapter inspired by `draft-jennings-ai-mcp-over-moq-00`.

MCP JSON-RPC semantics remain unchanged. The adapter maps MCP control and tool messages to MOQT-style namespaces/tracks and delegates tool behavior to the existing MCP server core and existing NATS tools.

Control uses two concurrently-open logical channels, `client-to-server` and `server-to-client`, under the same `(mcp, <session-id>, control)` namespace. Both channels share the same event-name vocabulary (`initialize`, `initialized`, `ping`, `pong`, `tools.call`, `tools.result`, etc.) and either channel can carry the same JSON-RPC payload when a client chooses to publish or subscribe that way.

Current v0 status: the public adapter and security interfaces are isolated for replacement by a real MOQT library; the listener uses a TCP/TLS compatibility shim so local development and JSON-RPC behavior can be exercised while MOQT support matures.
