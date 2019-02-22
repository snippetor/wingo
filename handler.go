package wingo

type Handler func(Context)

type HandlersChain interface {
	Append(...Handler)
	Set(...Handler)

	Next()
	Skip()
	Reset()
	Stop()
	IsStopped() bool
	Handlers() []Handler
}

type handlersChain struct {
	handlers     []Handler
	currentIndex int
	context      Context
}

func (c *handlersChain) Append(handlers ...Handler) {
	c.handlers = append(c.handlers, handlers...)
}

func (c *handlersChain) Set(handlers ...Handler) {
	c.handlers = handlers
}

func (c *handlersChain) Next() {
	if c.IsStopped() {
		return
	}
	c.handlers[c.currentIndex](c.context)
	c.currentIndex += 1
}

func (c *handlersChain) Skip() {
	c.currentIndex += 1
}

func (c *handlersChain) Reset() {
	c.currentIndex = 0
}

func (c *handlersChain) Stop() {
	c.currentIndex = len(c.handlers)
}

func (c *handlersChain) IsStopped() bool {
	return c.currentIndex >= len(c.handlers)
}

func (c *handlersChain) Handlers() []Handler {
	return c.handlers
}
