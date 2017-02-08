package kvstore

import (
	"testing"
)

func TestGet_UninitializedKey(t *testing.T) {
	store := New()
	valForUnusedKey := store.Get("unusedKey")
	if valForUnusedKey != "" {
		t.Errorf("Get(unusedKey) returned %s, expected empty string", valForUnusedKey)
	}
}

func TestSet_UninitializedKey(t *testing.T) {
	store := New()
	key := "1"
	setVal := "myNewVal_123"
	val := store.Set(key, setVal)
	if val != setVal {
		t.Errorf("Set(%s,%s) returned %s, expected %s", key, setVal, val, setVal)
	}
}

func TestTestSet_UninitializedKey(t *testing.T) {
	store := New()
	key := "1"
	testVal := "abc"
	setVal := "123"
	val := store.TestSet(key, testVal, setVal)
	if val != "" {
		t.Errorf("TestSet(%s,%s,%s) returned %s, expected empty string", key, testVal, setVal, val)
	}
}
