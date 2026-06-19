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
