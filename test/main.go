package main

import (
	"github.com/snippetor/wingo"
)

type TestController struct {
}

func (c *TestController) Test(ctx *wingo.Context) {
	ctx.RouteGet("/hello", false)

	ctx.Response.Header.Set("Authorization", "eyJJZCI6MjAwMDEsIk5hbWUiOiJjYXJsIiwiQWRtaW4iOnRydWV9/Sr4iGbfjVxTWt2RKktOHUoA9a_6mStiROy3yXns7dTw")
	ctx.Success(&map[string]string{"Status": "OK"})
}

func (c *TestController) Test1(ctx *wingo.Context) {
	ctx.RouteGet("/hello1", false)

	ctx.Success(&map[string]string{"Status": "OK"})
}

type TestController1 struct {
}

func (c *TestController1) Test(ctx *wingo.Context) {
	ctx.RouteGet("/hello2", false)

	ctx.Success(&map[string]string{"Status": "OK"})
}

func (c *TestController1) Test1(ctx *wingo.Context) {
	ctx.RouteGet("/hello3", false)

	ctx.Success(&map[string]string{"Status": "OK"})
}

func main() {
	wingo.Use(wingo.Recover, wingo.ContextLog, wingo.Latency)
	app := wingo.Default()
	app.Namespace("v1").
		Group("test").Route(&TestController{}).
		Group("test1").Route(&TestController1{})
	app.Run(8080)
}
