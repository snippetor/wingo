package wingo

import (
	"github.com/snippetor/logger"
	"github.com/valyala/fasthttp"
	"path"
	"reflect"
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

func (a *AppEngine) Route(ctrls ...interface{}) *AppEngine {
	ctx := &Context{RouteTester: &RouteTester{}}
	for _, ctrl := range ctrls {
		t := reflect.TypeOf(ctrl)
		for i := 0; i < t.NumMethod(); i++ {
			ctx.RouteTester.Method = ""
			ctx.RouteTester.Path = ""
			m := t.Method(i)
			var call = m.Func.Call
			if m.Type.NumIn() == 2 && m.Type.In(1) == reflect.TypeOf(ctx) {
				func() {
					defer func() {
						recover()
					}()
					call([]reflect.Value{reflect.ValueOf(ctrl), reflect.ValueOf(ctx)})
				}()
				if ctx.RouteTester.Method != "" && ctx.RouteTester.Path != "" {
					p := path.Join("/", a.tempBasePath, ctx.RouteTester.Path)
					a.router.handle(ctx.RouteTester.Method, p, func(c *Context) {
						call([]reflect.Value{reflect.ValueOf(ctrl), reflect.ValueOf(c)})
					})
				}
			}
		}
	}
	a.tempBasePath = ""
	return a
}

func (a *AppEngine) Run(port int) {
	Log.I("Server run on :%d", port)
	err := fasthttp.ListenAndServe(":"+strconv.Itoa(port), func(req *fasthttp.RequestCtx) {
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
		ctx.chain.Fire()
	})
	if err != nil {
		panic(err)
	}
}
