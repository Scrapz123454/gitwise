package llm

import (
	"os"
	"testing"
)

func TestCache(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitwise-cache-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cache := NewCache(dir)

	// Miss
	_, ok := cache.Get("prompt1")
	if ok {
		t.Error("expected cache miss, got hit")
	}

	// Set + Hit
	if err := cache.Set("prompt1", "response1"); err != nil {
		t.Fatalf("cache.Set failed: %v", err)
	}

	val, ok := cache.Get("prompt1")
	if !ok {
		t.Error("expected cache hit, got miss")
	}
	if val != "response1" {
		t.Errorf("cache.Get = %q, want %q", val, "response1")
	}

	// Different key is a miss
	_, ok = cache.Get("prompt2")
	if ok {
		t.Error("expected cache miss for different key")
	}

	// Overwrite
	if err := cache.Set("prompt1", "updated"); err != nil {
		t.Fatalf("cache.Set overwrite failed: %v", err)
	}
	val, ok = cache.Get("prompt1")
	if !ok || val != "updated" {
		t.Errorf("cache overwrite failed: ok=%v val=%q", ok, val)
	}
}
