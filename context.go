package wingo

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"sync/atomic"
)

var lastCapturedContextID uint32

func ContextLog(c *Context) {
	c.LogD(">>>[%s] %s %s", c.MethodString(), c.PathString(), Bytes2String(c.Request.Body()))
	c.Next()
	c.LogD("<<<[%s] %s %s", c.MethodString(), c.PathString(), Bytes2String(c.Response.Body()))
}

type Context struct {
	*fasthttp.RequestCtx
	id           uint32
	chain        *HandlersChain
	RouteTester  *RouteTester
	TokenPayload *TokenPayload
}

func (c *Context) Route(method, path string, needAuth bool) {
	if c.RouteTester != nil {
		c.RouteTester.Method = method
		c.RouteTester.Path = path
		c.RouteTester.NeedAuth = needAuth
		panic(TestRouteError(0))
	}
}

func (c *Context) RouteGet(path string, needAuth bool) {
	c.Route("GET", path, needAuth)
}

func (c *Context) RoutePost(path string, needAuth bool) {
	c.Route("POST", path, needAuth)
}

func (c *Context) RoutePut(path string, needAuth bool) {
	c.Route("PUT", path, needAuth)
}

func (c *Context) RouteDelete(path string, needAuth bool) {
	c.Route("DELETE", path, needAuth)
}

func (c *Context) RouteOptions(path string, needAuth bool) {
	c.Route("OPTIONS", path, needAuth)
}

func (c *Context) RouteHead(path string, needAuth bool) {
	c.Route("HEAD", path, needAuth)
}

func (c *Context) RouteConnect(path string, needAuth bool) {
	c.Route("CONNECT", path, needAuth)
}

func (c *Context) RouteTrace(path string, needAuth bool) {
	c.Route("TRACE", path, needAuth)
}

func (c *Context) RoutePatch(path string, needAuth bool) {
	c.Route("PATCH", path, needAuth)
}

func (c *Context) Id() uint32 {
	if c.id == 0 {
		// set the id here.
		forward := atomic.AddUint32(&lastCapturedContextID, 1)
		c.id = forward
	}
	return c.id
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

func (c *Context) SuccessAuth(tokenPayload *TokenPayload, body interface{}) {
	c.Response.SetStatusCode(fasthttp.StatusOK)
	bs, err := globalCodec.Marshal(body)
	if err != nil {
		panic(err)
	}
	c.Response.Header.Set("Authorization", marshalToken(tokenPayload))
	c.Response.SetBody(bs)
}

func (c *Context) MethodString() string {
	bytes := c.Method()
	return Bytes2String(bytes)
}

func (c *Context) PathString() string {
	bytes := c.Path()
	return Bytes2String(bytes)
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
