package wingo

import (
	"reflect"
	"runtime"
)

type Handler func(*Context)

func (h Handler) Name() string {
	pc := reflect.ValueOf(h).Pointer()
	return runtime.FuncForPC(pc).Name()
}

type HandlersChain struct {
	handlers     []Handler
	currentIndex int
	context      *Context
}

func (c *HandlersChain) Append(handlers ...Handler) {
	c.handlers = append(c.handlers, handlers...)
}

func (c *HandlersChain) Set(handlers ...Handler) {
	c.handlers = handlers
}

func (c *HandlersChain) Next() {
	if c.IsStopped() {
		return
	}
	c.handlers[c.currentIndex](c.context)
	c.currentIndex += 1
}

func (c *HandlersChain) Skip() {
	c.currentIndex += 1
}

func (c *HandlersChain) Reset() {
	c.currentIndex = 0
}

func (c *HandlersChain) Stop() {
	c.currentIndex = len(c.handlers)
}

func (c *HandlersChain) IsStopped() bool {
	return c.currentIndex >= len(c.handlers)
}

func (c *HandlersChain) Handlers() []Handler {
	return c.handlers
}

func (c *HandlersChain) CurrentHandlerName() string {
	if c.currentIndex >= len(c.handlers) {
		return "unknown"
	}
	return c.handlers[c.currentIndex].Name()
}
