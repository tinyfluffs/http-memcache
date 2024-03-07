package cache

import (
	"bytes"
	"testing"
	"time"
)

func TestExpiryOkResult(t *testing.T) {
	key := "key"
	val := []byte("val")
	expiration := 100 * time.Millisecond
	gc := 250 * time.Millisecond

	mc := New(expiration, gc, 1)
	mc.Store(key, val)

	time.Sleep(expiration)

	_, expiry, ok := mc.Get(key)
	if !expiry.Before(time.Now()) {
		t.Error("expiry is not in the past")
	}

	if ok {
		t.Error("ok is true when key is in the past, should be false")
	}
}

func TestGet(t *testing.T) {
	key := "key"
	val := []byte("val")
	expiration := 100 * time.Millisecond
	gc := 250 * time.Millisecond

	mc := New(expiration, gc, 1)
	mc.Store(key, val)

	resultVal, expiry, ok := mc.Get(key)
	if !ok {
		t.Error("no value found")
	}

	if !bytes.Equal(val, resultVal) {
		t.Error("value data not equal")
	}

	if expiry.IsZero() {
		t.Error("zero time")
	}
}

func TestStore(t *testing.T) {
	key := "key"
	val := []byte("val")
	expiration := 100 * time.Millisecond
	gc := 250 * time.Millisecond

	mc := New(expiration, gc, 1)

	_, _, ok := mc.Get(key)
	if ok {
		t.Error("value found before expected")
	}

	mc.Store(key, val)

	_, _, ok = mc.Get(key)
	if !ok {
		t.Error("no value found")
	}
}

func TestExpiryGCCleanup(t *testing.T) {
	key := "key"
	val := []byte("val")
	expiration := 100 * time.Millisecond
	gc := 250 * time.Millisecond

	mc := New(expiration, gc, 1)
	mc.Store(key, val)

	time.Sleep(expiration)
	mc.GC()

	getVal, _, _ := mc.Get(key)
	if len(getVal) > 0 {
		t.Error("gc did not cleanup expired key")
	}
}

func TestExpiryGCOk(t *testing.T) {
	key := "key"
	val := []byte("val")
	expiration := 100 * time.Millisecond
	gc := 250 * time.Millisecond

	mc := New(expiration, gc, 1)
	mc.Store(key, val)

	time.Sleep(expiration)
	mc.GC()

	_, _, ok := mc.Get(key)
	if ok {
		t.Error("gc did not cleanup expired key")
	}
}

func TestHit(t *testing.T) {
	key := "key"
	val := []byte("val")
	expiration := 30 * time.Minute
	gc := 1 * time.Second

	mc := New(expiration, gc, 1)
	mc.Store(key, val)

	if _, _, ok := mc.Get(key); !ok {
		t.Error("key not found")
	}
}
