// do not build if using custom types
// +build !lrumap_custom

package lrumap_test

import (
	"fmt"
	"strconv"
	"time"

	"github.com/db47h/lrumap"
)

// removeCallback will be called upon entry removal from the lrumap.
func removeCallback(v lrumap.Wrapper) {
	fmt.Printf("Removed entry with key: %q, value: %v\n", v.Key(), v.Unwrap())
}

// A simple showcase were we store ints with string keys.
func Example() {
	// make a new lru map with a maximum of 10 entries.
	m, err := lrumap.New(10, lrumap.RemoveFunc(removeCallback))
	if err != nil {
		panic(err)
	}
	// and fill it
	for i := 0; i < m.Cap(); i++ {
		m.Set(strconv.Itoa(i), i)
		// the sleep is for testing purposes only
		// so that we have a different timestamp for every entry
		time.Sleep(10 * time.Millisecond)
	}

	// xyz does not exists, Get returns nil
	v := m.Get("xyz")
	if v != nil {
		panic("found unexpected entry")
	}

	// "0" exists, it will be refreshed and pushed back, "1" should now be the LRU entry)
	v = m.Get("0")
	if v == nil {
		panic("entry 0 does not exist")
	}

	// this should trigger removal of "1"
	m.Set("11", 11)

	// now update 2, should trigger removal of old "2"
	m.Set("2", 222)
	v, _ = m.GetWithDefault("2", func(key lrumap.Key) (lrumap.Value, error) {
		panic("here, we should not be called")
	})
	if v.Unwrap().(int) != 222 {
		panic("Got " + strconv.Itoa(v.Unwrap().(int)))
	}

	// Try to get "12". Will create a new one and delete "3"
	v, _ = m.GetWithDefault("12", func(key lrumap.Key) (lrumap.Value, error) {
		return 12, nil
	})
	if v.Unwrap().(int) != 12 {
		panic("Expected 12, got " + strconv.Itoa(v.Unwrap().(int)))
	}

	// manually delete "5"
	m.Delete("5")

	// Output:
	// Removed entry with key: "1", value: 1
	// Removed entry with key: "2", value: 2
	// Removed entry with key: "3", value: 3
	// Removed entry with key: "5", value: 5
}
