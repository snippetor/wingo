package wingo

type RouterMethod struct {
	path       string
	method     string
	middleware []Handler
}

type Router struct {
	routes map[string]map[string][]Handler
}

func newRounter() *Router {
	return &Router{routes: make(map[string]map[string][]Handler)}
}

func (r *Router) buildHandler(h Handler) Handler {
	return func(ctx *Context) {
		if !ctx.Proceed(h) {
			ctx.Next()
		}
	}
}

func (r *Router) apply(handlers *[]Handler) bool {
	tmp := *handlers
	for i, h := range tmp {
		if h == nil {
			if len(tmp) == 1 {
				return false
			}
			continue
		}
		(*handlers)[i] = r.buildHandler(h)
	}
	return true
}

func (r *Router) handle(path string, method string, handlers ...Handler) {
	r.apply(&handlers)
	if _, found := r.routes[method]; !found {
		r.routes[method] = make(map[string][]Handler)
	}
	r.routes[method][path] = append(r.routes[method][path], handlers...)
}

func (r *Router) Get(path string, handlers ...Handler) {
}

func (r *Router) Post(path string, handlers ...Handler) {
}

func (r *Router) getRequestHandlers(method, path string) []Handler {
	var newHandlers []Handler
	if route, found := r.routes[method]; found {
		if handlers, found := route[path]; found {
			if len(handlers) > 0 {
				copy(newHandlers, globalMiddleWares)
				r.apply(&newHandlers)
				newHandlers = append(newHandlers, handlers...)
			}
		}
	}
	return newHandlers
}
