package main

import (
	"github.com/snippetor/wingo"
)

type TestController struct {
}

func (c *TestController) Route(router *wingo.Router) {
	router.Post("/hello", c.Test)
}

func (c *TestController) Test(ctx *wingo.Context) {
	ctx.Success(&map[string]string{"Status": "OK"})
}

func main() {
	wingo.Use(wingo.Recover, wingo.ContextLog, wingo.Latency)
	app := wingo.Default()
	app.Group("v1").Group("test").Route(&TestController{})
	app.Run(8080)
}
