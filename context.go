package wingo

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

func contextLog(c Context) {

}

type Context interface {
	Next()
	Stop()
}

func NewContext() Context {
	c := &context{}
	c.chain = &handlersChain{context: c}
	return c
}

type context struct {
	chain HandlersChain
	ctx   *fasthttp.RequestCtx
}

func (c *context) Next() {
	c.chain.Next()
}

func (c *context) Stop() {
	c.chain.Stop()
}

func (c *context) IsStopped() bool {
	return c.chain.IsStopped()
}

func (c *context) RequestBody(body interface{}) {
	json.Unmarshal(c.ctx.Request.Body(), body)
}

func (c *context) ResponseOK(body interface{}) {
	c.ctx.Response.SetStatusCode(fasthttp.StatusOK)
	bs, err := c.Codec.Marshal(body)
	if err != nil {
		panic(err)
	}
	c.ctx.Response.SetBody(bs)
	c.LogD("<<< %s %s", string(c.RequestCtx.Path()), string(bs))
}

func (c *context) ResponseFailed(reason string) {
	c.ctx.Response.SetStatusCode(fasthttp.StatusOK)
	params := make(map[string]interface{})
	params["error"] = reason
	bs, err := c.Codec.Marshal(params)
	if err != nil {
		panic(err)
	}
	c.RequestCtx.Response.SetBody(bs)
	c.LogD("<<< %s %s", string(c.RequestCtx.Path()), string(bs))
}
