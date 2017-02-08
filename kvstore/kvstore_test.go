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

func TestGet_AfterSet(t *testing.T) {
	store := New()
	key := "id_123"
	setVal := "abc"
	store.Set(key, setVal)
	val := store.Get(key)
	if val != setVal {
		t.Errorf("Get(%s) returned %s, expected %s", key, val, setVal)
	}
}

func TestSet_AfterSet(t *testing.T) {
	store := New()
	key := "id_123"
	origValue := "original_value_123"
	store.Set(key, origValue)
	newVal := "abc"
	val := store.Set(key, newVal)
	if val != newVal {
		t.Errorf("Set(%s, %s) returned %s, expected %s", key, newVal, val, newVal)
	}
}

func TestTest_AfterSet(t *testing.T) {
	store := New()
	key := "id_123"
	origVal := "abc"
	store.Set(key, origVal)

	testSetVal := "newVal"
	val := store.TestSet(key, "someOtherVal", testSetVal)
	if val != origVal {
		t.Errorf("TestSet(%s, someOtherVal, %s) returned %s, expected %s", key, testSetVal, val, origVal)
	}
	val = store.TestSet(key, origVal, testSetVal)
	if val != testSetVal {
		t.Errorf("TestSet(%s, %s, %s) returned %s, expected %s", key, origVal, testSetVal, val, testSetVal)
	}
}
