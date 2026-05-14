package cacher_test

import (
	"testing"
	"time"

	"github.com/yourorg/patchwork/internal/cacher"
)

func TestSet_And_Get_Hit(t *testing.T) {
	c := cacher.New(5 * time.Second)
	c.Set("key1", "hello")
	v, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v.(string) != "hello" {
		t.Fatalf("expected 'hello', got %v", v)
	}
}

func TestGet_Miss_Unknown(t *testing.T) {
	c := cacher.New(5 * time.Second)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss for unknown key")
	}
}

func TestGet_Miss_Expired(t *testing.T) {
	c := cacher.New(1 * time.Millisecond)
	c.SetTTL("exp", 42, 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("exp")
	if ok {
		t.Fatal("expected cache miss after expiry")
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	c := cacher.New(5 * time.Second)
	c.Set("k", true)
	c.Delete("k")
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cacher.New(5 * time.Second)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected empty cache after flush, got %d", c.Len())
	}
}

func TestEvict_RemovesExpired(t *testing.T) {
	c := cacher.New(5 * time.Second)
	c.SetTTL("short", "x", 1*time.Millisecond)
	c.Set("long", "y")
	time.Sleep(5 * time.Millisecond)
	n := c.Evict()
	if n != 1 {
		t.Fatalf("expected 1 eviction, got %d", n)
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", c.Len())
	}
}

func TestLen_ReflectsCount(t *testing.T) {
	c := cacher.New(5 * time.Second)
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
	c.Set("x", 1)
	c.Set("y", 2)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestSetTTL_OverridesDefault(t *testing.T) {
	c := cacher.New(1 * time.Hour)
	c.SetTTL("quick", "v", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("quick")
	if ok {
		t.Fatal("expected entry to expire despite long default TTL")
	}
}
