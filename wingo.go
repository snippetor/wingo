package wingo

import (
	"fmt"
	"github.com/snippetor/logger"
	"github.com/valyala/fasthttp"
	"path"
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

func Use(middlewares ...Handler) {
	globalMiddleWares = append(globalMiddleWares, middlewares...)
}

type AppEngine struct {
	config       *log.Config
	router       *Router
	tempBasePath string
}

func Default() *AppEngine {
	Log = log.NewLogger(config)
	return &AppEngine{router: newRouter()}
}

func (a *AppEngine) Group(group string) *AppEngine {
	a.tempBasePath = path.Join(a.tempBasePath, group)
	return a
}

func (a *AppEngine) Route(ctrls ...Controller) *AppEngine {
	a.router.TempBasePath = path.Join("/", a.tempBasePath)
	for _, ctrl := range ctrls {
		if IsController(ctrl) {
			ctrl.Route(a.router)
		}
	}
	a.router.TempBasePath = ""
	a.tempBasePath = ""
	return a
}

func (a *AppEngine) Run(port int) {
	Log.I("Server run on :%d", port)
	err := fasthttp.ListenAndServe(":"+strconv.Itoa(port), func(req *fasthttp.RequestCtx) {
		fmt.Println(string(req.Path()))
		handlers := a.router.getRequestHandlers(string(req.Method()), string(req.Path()))
		if len(handlers) == 0 {
			req.Error("404 Not found ;(", fasthttp.StatusNotFound)
			return
		}
		ctx := ctxPool.Get().(*Context)
		defer ctxPool.Put(ctx)
		ctx.id = 0
		ctx.RequestCtx = req
		switch globalCodec.Name() {
		case "json":
			ctx.SetContentType("application/json")
		case "protobuf":
			ctx.SetContentType("application/protobuf")
		}
		if ctx.chain == nil {
			ctx.chain = &HandlersChain{context: ctx}
		}
		ctx.chain.Reset()
		ctx.chain.Set(handlers)
		ctx.chain.Next()
	})
	if err != nil {
		panic(err)
	}
}
