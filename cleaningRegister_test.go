package cleaningRegister

import (
	"testing"
	"time"
)

func TestPutAndGet(t *testing.T) {
	creg := New(1*time.Minute, nil, nil)

	key, val := "key1", "val1"
	creg.Put(key, val)
	out, ok := creg.Get(key)

	if !ok {
		t.Fatal("OK should be true for existing key")
	}

	if out != val {
		t.Fatal("Got", out, " expected val1")
	}
}

func TestPutAndPop(t *testing.T) {
	creg := New(1*time.Minute, nil, nil)

	key, val := "key1", "val1"
	creg.Put(key, val)
	out, ok := creg.Pop(key)

	if !ok {
		t.Fatal("OK should be true for existing key")
	}

	if out != val {
		t.Fatal("Got", out, " expected val1")
	}

	out, ok = creg.Get(key)

	if ok || out != nil {
		t.Fatal("Key should be gone after Pop")
	}
}

func TestGetValues(t *testing.T) {
	creg := New(1*time.Minute, nil, nil)

	key1, val1 := "key1", "val1"
	key2, val2 := "key2", "val2"
	creg.Put(key1, val1)
	creg.Put(key2, val2)

	vals := creg.Copy()

	if len(vals) != 2 {
		t.Fatal("GetValues should return two items")
	}

	if vals[key1] != val1 || vals[key2] != val2 {
		t.Fatal("Copy did not return a true copy. Actual: ", vals)
	}

}
