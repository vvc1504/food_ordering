package utils

import (
	"crypto/rand"
	"fmt"
)

// Call is a helper to wrap a synchronous function in an asynchronous channel-based response.
func Call[S any](fn func() S) <-chan S {
	stCh := make(chan S, 1)
	go func() {
		defer close(stCh)
		st := fn()
		stCh <- st
	}()
	return stCh
}

// Call is a helper to wrap a synchronous function in an asynchronous channel-based response.
func Call2[T any, S any](fn func() (T, S)) (<-chan T, <-chan S) {
	resCh := make(chan T, 1)
	stCh := make(chan S, 1)
	go func() {
		defer close(resCh)
		defer close(stCh)
		res, st := fn()
		resCh <- res
		stCh <- st
	}()
	return resCh, stCh
}

func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x-%x", b[0:4], b[4:8])
}
