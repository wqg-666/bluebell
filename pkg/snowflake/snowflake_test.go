package snowflake

import "testing"

func TestInitRejectsInvalidStartTime(t *testing.T) {
	if err := Init("not-a-date", 1); err == nil {
		t.Fatal("expected invalid start time to fail")
	}
}

func TestInitAndGenID(t *testing.T) {
	if err := Init("2026-01-01", 1); err != nil {
		t.Fatalf("initialize snowflake: %v", err)
	}
	first := GenID()
	second := GenID()
	if first <= 0 || second <= 0 || first == second {
		t.Fatalf("expected distinct positive IDs, got %d and %d", first, second)
	}
}
