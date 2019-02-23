package wingo

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"runtime"
)

func Recover(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			if ctx.IsStopped() {
				return
			}

			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break

				}

				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}

			// when stack finishes
			logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.chain.CurrentHandlerName())
			logMessage += fmt.Sprintf("At Request: %s\n", fmt.Sprintf("%s %s %s", string(ctx.Path()), string(ctx.Method()), ctx.RemoteIP().String()))
			logMessage += fmt.Sprintf("Trace: %s\n", err)
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			ctx.LogE(logMessage)
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.ResponseBody(map[string]string{error: fmt.Sprintf("%s", err)})
			ctx.Stop()
		}
	}()

	ctx.Next()
}
