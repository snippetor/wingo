package wingo

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"sync/atomic"
)

var lastCapturedContextID uint32

func ContextLog(c *Context) {
	c.LogD(">>>[%s] %s %s", string(c.Method()), string(c.Path()), string(c.Request.Body()))
	c.Next()
	c.LogD("<<<[%s] %s %s", string(c.Method()), string(c.Path()), string(c.Response.Body()))
}

type Context struct {
	*fasthttp.RequestCtx
	id          uint32
	chain       *HandlersChain
	RouteTester *RouteTester
}

func (c *Context) Route(method, path string) {
	if c.RouteTester != nil {
		c.RouteTester.Method = method
		c.RouteTester.Path = path
		panic(TestRouteError(0))
	}
}

func (c *Context) RouteGet(path string) {
	c.Route("GET", path)
}

func (c *Context) RoutePost(path string) {
	c.Route("POST", path)
}

func (c *Context) RoutePut(path string) {
	c.Route("PUT", path)
}

func (c *Context) RouteDelete(path string) {
	c.Route("DELETE", path)
}

func (c *Context) RouteOptions(path string) {
	c.Route("OPTIONS", path)
}

func (c *Context) Id() uint32 {
	if c.id == 0 {
		// set the id here.
		forward := atomic.AddUint32(&lastCapturedContextID, 1)
		c.id = forward
	}
	return c.id
}

func (c *Context) Proceed(h Handler) bool {
	beforeIdx := c.chain.currentIndex
	h(c)
	if c.chain.currentIndex > beforeIdx && !c.IsStopped() {
		return true
	}
	return false
}

func (c *Context) Next() {
	c.chain.Next()
}

func (c *Context) Stop() {
	c.chain.Stop()
}

func (c *Context) IsStopped() bool {
	return c.chain.IsStopped()
}

func (c *Context) RequestBody(body interface{}) {
	err := globalCodec.Unmarshal(c.Request.Body(), body)
	if err != nil {
		panic(err)
	}
}

func (c *Context) Success(body interface{}) {
	c.Response.SetStatusCode(fasthttp.StatusOK)
	bs, err := globalCodec.Marshal(body)
	if err != nil {
		panic(err)
	}
	c.Response.SetBody(bs)
}

func (c *Context) LogE(format string, v ...interface{}) {
	Log.E(fmt.Sprintf("[%010v] ", c.Id())+format, v...)
}

func (c *Context) LogD(format string, v ...interface{}) {
	Log.D(fmt.Sprintf("[%010v] ", c.Id())+format, v...)
}

func (c *Context) LogW(format string, v ...interface{}) {
	Log.W(fmt.Sprintf("[%010v] ", c.Id())+format, v...)
}

func (c *Context) LogI(format string, v ...interface{}) {
	Log.I(fmt.Sprintf("[%010v] ", c.Id())+format, v...)
}
