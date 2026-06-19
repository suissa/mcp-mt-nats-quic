package security

import "testing"

func TestMTLSConfigGeneration(t *testing.T) {
	cfg, err := BuildTLSConfig(TLSConfig{RequireMTLS: true, AllowInsecure: true})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ClientAuth == 0 || len(cfg.Certificates) == 0 {
		t.Fatalf("bad tls config")
	}
}
