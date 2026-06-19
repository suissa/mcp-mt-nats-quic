# Experimental MCP-Scalable-Channel MOQT/QUIC Transport

This package contains an experimental MCP-over-MOQT/QUIC transport adapter inspired by `draft-jennings-ai-mcp-over-moq-00`.

MCP JSON-RPC semantics remain unchanged. The adapter maps MCP control and tool messages to MOQT-style namespaces/tracks and delegates tool behavior to the existing MCP server core and existing NATS tools.

Current v0 status: the public adapter and security interfaces are isolated for replacement by a real MOQT library; the listener uses a TCP/TLS compatibility shim so local development and JSON-RPC behavior can be exercised while MOQT support matures.
