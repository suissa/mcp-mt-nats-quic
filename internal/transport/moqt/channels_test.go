package moqt

import "testing"

func TestParseChannelBackendsSupportsRequestedCarriers(t *testing.T) {
	backends, err := ParseChannelBackends("nats,quic,kafka,redpanda,rabbitmq,bullmq,redis-streams,grpc")
	if err != nil {
		t.Fatal(err)
	}
	want := []ChannelBackend{ChannelBackendNATS, ChannelBackendQUIC, ChannelBackendKafka, ChannelBackendRedpanda, ChannelBackendRabbitMQ, ChannelBackendBullMQ, ChannelBackendRedisStreams, ChannelBackendGRPC}
	if len(backends) != len(want) {
		t.Fatalf("got %d backends, want %d", len(backends), len(want))
	}
	for i := range want {
		if backends[i] != want[i] {
			t.Fatalf("backend[%d]=%q want %q", i, backends[i], want[i])
		}
	}
}

func TestParseChannelBackendsAliasesAndDedupes(t *testing.T) {
	backends, err := ParseChannelBackends("moqt-quic,amqp,redis,grpc,quic")
	if err != nil {
		t.Fatal(err)
	}
	want := []ChannelBackend{ChannelBackendQUIC, ChannelBackendRabbitMQ, ChannelBackendRedisStreams, ChannelBackendGRPC}
	if len(backends) != len(want) {
		t.Fatalf("got %#v want %#v", backends, want)
	}
	for i := range want {
		if backends[i] != want[i] {
			t.Fatalf("backend[%d]=%q want %q", i, backends[i], want[i])
		}
	}
}

func TestParseChannelBackendsRejectsUnsupported(t *testing.T) {
	if _, err := ParseChannelBackends("nats,unknown"); err == nil {
		t.Fatal("expected unsupported backend error")
	}
}
