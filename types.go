package onion

import (
	"sync"
	"time"
)

var (
	globalLock = &sync.RWMutex{}
)

// String is the string value holder, with safe lock for concurrent reload
type String interface {
	String() string
}

// Int is the int value holder, with safe lock for concurrent reload
type Int interface {
	// Return the int value
	Int() int
	// Return int64 value
	Int64() int64
	// Return duration value
	Duration() time.Duration
}

// Bool is the bool value holder, with safe lock for concurrent reload
type Bool interface {
	// The bool value
	Bool() bool
}

// Float is the float value holder, with safe lock for concurrent reload
type Float interface {
	Float32() float32
	Float64() float64
}

type stringHolder struct {
	value *string
}

func (sh stringHolder) String() string {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return *sh.value
}

type intHolder struct {
	value *int64
}

func (ih intHolder) Int() int {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return int(*ih.value)
}

func (ih intHolder) Int64() int64 {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return *ih.value
}

func (ih intHolder) Duration() time.Duration {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return time.Duration(*ih.value)
}

type boolHolder struct {
	value *bool
}

func (bh boolHolder) Bool() bool {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return *bh.value
}

type floatHolder struct {
	value *float64
}

func (fh floatHolder) Float32() float32 {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return float32(*fh.value)
}

func (fh floatHolder) Float64() float64 {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return *fh.value
}
