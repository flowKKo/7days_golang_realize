package singleflight

import "sync"

type call struct {
	// call represents the request which is currently
	// executing or has completed
	// a waitGroup object can wait the end of a set of goroutine
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group manages different requests of key (call)
type Group struct {
	mu sync.Mutex // protects m
	m  map[string]*call
}

// Do for the same key, no matter how many times Do is called
// function fn will be executed only once
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call) // delayed initialization
	}

	// if the request is ongoing, wait till the end of request
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)  // add lock before request
	g.m[key] = c // add call to g, informing that corresponding request is already ongoing
	g.mu.Unlock()

	c.val, c.err = fn() // call fn, send the requests
	c.wg.Done()         // the end of request

	g.mu.Lock()
	delete(g.m, key) // update g.ms
	g.mu.Unlock()

	return c.val, c.err // return the value
}
