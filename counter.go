package main

import "sync"

type SafeCounter struct {
	n   int
	mux sync.Mutex
}

func (c *SafeCounter) Inc() {
	c.mux.Lock()
	c.n++
	c.mux.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.n
}
