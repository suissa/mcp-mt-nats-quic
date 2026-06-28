package moqt

import (
	"context"
	"fmt"
	"strings"
)

// ChannelBackend identifies an experimental MCP-Scalable-Channel carrier.
// The MCP JSON-RPC payload and EventName vocabulary stay identical across all
// carriers; only the delivery adapter changes.
type ChannelBackend string

const (
	ChannelBackendNATS         ChannelBackend = "nats"
	ChannelBackendQUIC         ChannelBackend = "quic"
	ChannelBackendKafka        ChannelBackend = "kafka"
	ChannelBackendRedpanda     ChannelBackend = "redpanda"
	ChannelBackendRabbitMQ     ChannelBackend = "rabbitmq"
	ChannelBackendBullMQ       ChannelBackend = "bullmq"
	ChannelBackendRedisStreams ChannelBackend = "redis-streams"
	ChannelBackendGRPC         ChannelBackend = "grpc"
)

var SupportedChannelBackends = []ChannelBackend{
	ChannelBackendNATS,
	ChannelBackendQUIC,
	ChannelBackendKafka,
	ChannelBackendRedpanda,
	ChannelBackendRabbitMQ,
	ChannelBackendBullMQ,
	ChannelBackendRedisStreams,
	ChannelBackendGRPC,
}

// ChannelMessage is the backend-neutral representation used by future NATS,
// Kafka, Redpanda, RabbitMQ, BullMQ, Redis Streams, and gRPC adapters.
type ChannelMessage struct {
	Backend   ChannelBackend
	Namespace TrackNamespace
	Track     TrackName
	Event     EventName
	GroupID   uint64
	ObjectID  uint64
	Payload   []byte
	Metadata  map[string]string
}

// ChannelAdapter is the common interface for all scalable-channel carriers.
// Implementations must not change MCP JSON-RPC semantics or expose NATS tool
// internals directly to the client.
type ChannelAdapter interface {
	Backend() ChannelBackend
	Publish(ctx context.Context, msg ChannelMessage) error
	Subscribe(ctx context.Context, namespace TrackNamespace, track TrackName, events ...EventName) (<-chan ChannelMessage, error)
	Close(ctx context.Context) error
}

// Registry tracks configured backends without forcing optional dependencies into
// the binary until a concrete adapter is added.
type Registry struct {
	backends []ChannelBackend
}

func NewRegistry(backends []ChannelBackend) Registry {
	if len(backends) == 0 {
		backends = []ChannelBackend{ChannelBackendQUIC}
	}
	return Registry{backends: append([]ChannelBackend(nil), backends...)}
}

func (r Registry) Backends() []ChannelBackend { return append([]ChannelBackend(nil), r.backends...) }

func ParseChannelBackends(raw string) ([]ChannelBackend, error) {
	if strings.TrimSpace(raw) == "" {
		return []ChannelBackend{ChannelBackendQUIC}, nil
	}
	parts := strings.Split(raw, ",")
	backends := make([]ChannelBackend, 0, len(parts))
	seen := map[ChannelBackend]struct{}{}
	for _, part := range parts {
		backend, err := NormalizeChannelBackend(part)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[backend]; ok {
			continue
		}
		seen[backend] = struct{}{}
		backends = append(backends, backend)
	}
	return backends, nil
}

func NormalizeChannelBackend(raw string) (ChannelBackend, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "nats":
		return ChannelBackendNATS, nil
	case "quic", "moqt", "moqt-quic", "mcp-moqt":
		return ChannelBackendQUIC, nil
	case "kafka":
		return ChannelBackendKafka, nil
	case "redpanda":
		return ChannelBackendRedpanda, nil
	case "rabbitmq", "rabbit", "amqp":
		return ChannelBackendRabbitMQ, nil
	case "bullmq", "bull-mq", "bull":
		return ChannelBackendBullMQ, nil
	case "redis-streams", "redis_streams", "redisstreams", "redis":
		return ChannelBackendRedisStreams, nil
	case "grpc", "g-rpc":
		return ChannelBackendGRPC, nil
	default:
		return "", fmt.Errorf("unsupported scalable channel backend %q (supported: %s)", raw, SupportedChannelBackendsCSV())
	}
}

func SupportedChannelBackendsCSV() string {
	labels := make([]string, 0, len(SupportedChannelBackends))
	for _, backend := range SupportedChannelBackends {
		labels = append(labels, string(backend))
	}
	return strings.Join(labels, ",")
}
