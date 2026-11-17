package jwt

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGenTokenAndParseToken(t *testing.T) {
	viper.Set("auth.jwt_secret", "unit-test-signing-key-not-for-deployment")
	t.Cleanup(func() { viper.Set("auth.jwt_secret", nil) })
	viper.Set("auth.jwt_expire", 1)
	token, err := GenToken(42, "alice")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if token == "" {
		t.Fatal("generated token is empty")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("parse generated token: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "alice" {
		t.Fatalf("unexpected claims: user_id=%d username=%q", claims.UserID, claims.Username)
	}
}

func TestParseTokenRejectsMalformedToken(t *testing.T) {
	viper.Set("auth.jwt_secret", "unit-test-signing-key-not-for-deployment")
	t.Cleanup(func() { viper.Set("auth.jwt_secret", nil) })
	if _, err := ParseToken("not-a-jwt"); err == nil {
		t.Fatal("expected malformed token to be rejected")
	}
}

func TestGenTokenRequiresSigningSecret(t *testing.T) {
	viper.Set("auth.jwt_secret", "")
	t.Cleanup(func() { viper.Set("auth.jwt_secret", nil) })
	if _, err := GenToken(42, "alice"); err == nil {
		t.Fatal("expected missing signing secret to be rejected")
	}
}
