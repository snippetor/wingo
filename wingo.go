package wingo

import (
	"github.com/snippetor/logger"
	"github.com/valyala/fasthttp"
	"strconv"
	"sync"
)

var (
	config            *log.Config
	Log               log.Logger
	ctxPool           *sync.Pool
	globalMiddleWares []Handler
)

func init() {
	config = &log.Config{
		Level:              log.Info,
		OutputType:         log.Console,
		LogFileRollingType: log.RollingNone,
		LogFileOutputDir:   ".",
		LogFileName:        "wingo",
	}
	ctxPool = &sync.Pool{}
	ctxPool.New = func() interface{} {
		return &Context{}
	}
}

func GetLogConfig() *log.Config {
	return config
}

func SetLogConfig(c *log.Config) {
	config = c
	if Log != nil {
		Log.Close()
		Log = log.NewLogger(c)
	}
}

func Use(middlewares ... Handler) {
	globalMiddleWares = append(globalMiddleWares, middlewares...)
}

type AppEngine struct {
	config *log.Config
	router *Router
}

func Default() *AppEngine {
	Log = log.NewLogger(config)
	Log.SetPrefixes("[Wingo]")
	return &AppEngine{router: &Router{}}
}

func (a *AppEngine) Run(port int) {
	err := fasthttp.ListenAndServe(":"+strconv.Itoa(port), func(req *fasthttp.RequestCtx) {
		ctx := ctxPool.Get().(*Context)
		defer ctxPool.Put(ctx)
		ctx.id = 0
		ctx.RequestCtx = *req
		if ctx.chain == nil {
			ctx.chain = &HandlersChain{}
		}
		ctx.chain.Set(a.router.getRequestHandlers(string(req.Method()), string(req.Path()))...)
		ctx.chain.Next()
	})
	if err != nil {
		panic(err)
	}
}
