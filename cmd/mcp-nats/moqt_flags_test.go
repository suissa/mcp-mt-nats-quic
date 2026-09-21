package main

import "testing"

func TestValidateConfigAcceptsMOQTTransportAliases(t *testing.T) {
	cfg := &Config{Transport: "quic", MOQTAddress: "127.0.0.1:9443"}
	if err := validateConfig(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Transport != "moqt-quic" {
		t.Fatalf("got %s", cfg.Transport)
	}
}

func TestValidateConfigParsesScalableChannelBackends(t *testing.T) {
	cfg := &Config{Transport: "moqt-quic", MOQTAddress: "127.0.0.1:9443", ChannelBackendsRaw: "nats,kafka,redpanda,rabbitmq,bullmq,redis-streams,grpc"}
	if err := validateConfig(cfg); err != nil {
		t.Fatal(err)
	}
	if len(cfg.ChannelBackends) != 7 {
		t.Fatalf("got %d channel backends", len(cfg.ChannelBackends))
	}
}
