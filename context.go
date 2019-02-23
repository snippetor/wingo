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
	fasthttp.RequestCtx
	id    uint32
	chain *HandlersChain
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
	globalCodec.Unmarshal(c.Request.Body(), body)
}

func (c *Context) ResponseBody(body interface{}) {
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
