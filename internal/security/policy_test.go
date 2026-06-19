package security

import (
	"context"
	"testing"
)

func TestDPoPPolicyDisabledAllows(t *testing.T) {
	_, err := (Policy{}).DPoPVerifier().Verify(context.Background(), "", DPoPRequestContext{})
	if err != nil {
		t.Fatal(err)
	}
}
func TestStrictDPoPVerifierDenyAndReplay(t *testing.T) {
	v := NewStrictDPoPVerifier(0)
	if _, err := v.Verify(context.Background(), "", DPoPRequestContext{}); err == nil {
		t.Fatal("expected missing proof error")
	}
	if _, err := v.Verify(context.Background(), "jti", DPoPRequestContext{}); err != nil {
		t.Fatal(err)
	}
	if _, err := v.Verify(context.Background(), "jti", DPoPRequestContext{}); err == nil {
		t.Fatal("expected replay error")
	}
}
