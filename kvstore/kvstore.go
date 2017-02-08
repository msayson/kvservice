package kvstore

import (
	"sync"
)

// Data structure for values in the key-value store
type storeValue struct {
	value string
}

// Main data structure for key-value store
type KVStore struct {
	kvstore map[string]*storeValue // maps keys to values
	lock    *sync.RWMutex          // read/write mutex for safe concurrent access
}

func New() *KVStore {
	var store KVStore
	// Initialize key-value store
	store.kvstore = make(map[string]*storeValue)
	// Initialize read/write mutex
	store.lock = &sync.RWMutex{}
	return &store
}

func (store KVStore) Get(key string) string {
	// Acquire mutex for read access to kvstore
	store.lock.RLock()
	// Defer mutex unlock to function exit
	defer store.lock.RUnlock()

	// Look up and return store's value
	storeVal := store.lookup(key)
	return storeVal.value
}

func (store KVStore) Set(key string, value string) string {
	// Acquire mutex for exclusive access to kvstore
	store.lock.Lock()
	// Defer mutex unlock to function exit
	defer store.lock.Unlock()

	// Initialize entry and set to given value
	storeVal := store.lookup(key)
	storeVal.value = value
	return value
}

func (store KVStore) TestSet(key string, testVal string, setVal string) string {
	// Acquire mutex for exclusive access to kvstore
	store.lock.Lock()
	// Defer mutex unlock to function exit
	defer store.lock.Unlock()

	// Initialize and check value for key
	storeVal := store.lookup(key)

	// Execute the test-set
	if storeVal.value == testVal {
		storeVal.value = setVal
	}
	return storeVal.value
}

// Return the value associated with the given key
// Initializes its value to an empty string if this is the first access
func (store KVStore) lookup(key string) *storeValue {
	val := store.kvstore[key]
	if val == nil {
		// key used for the first time: create and initialize a storeValue
		val = &storeValue{
			value: "",
		}
		store.kvstore[key] = val
	}
	return val
}
