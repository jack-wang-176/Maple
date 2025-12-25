package web

//总的来说，这里的设计思路是替代了原有路由
//使用了一个新的路由
//但是之前的属性会被继承

type RouterGroup struct {
	Handlers []HandleFunc
	Route    *Route
	prefix   string
}

//我在想把group和use方法分开有必要吗
//有必要
//如果不分开无法调用RouterGroup添加处理函数
//因为Route里面我们没有存handler
//所以必须要在这里来存

func (r *Route) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		prefix: prefix,
		Route:  r,
	}
}

// Group 这里可以直接在httpService层面实现对方法的调用
func (httpService *HttpService) Group(prefix string) *RouterGroup {
	return httpService.Route.Group(prefix)
}

// Next 处理多个业务函数应该用next处理
// 这样的话保证一个业务逻辑正确结束以后正确的业务逻辑可以进行
func (ctx *Context) Next() {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		//执行业务逻辑
		ctx.handlers[ctx.index](ctx)
	}
}

// Use 这里仅仅只是添加
// 实际起作用我们可以去在新去定义AddRoute方法
func (g *RouterGroup) Use(handlers ...HandleFunc) {
	g.Handlers = append(g.Handlers, handlers...)
}

func (g *RouterGroup) AddRoute(path string, method string, handlers ...HandleFunc) {
	finalPath := g.prefix + path
	finalhandler := append(g.Handlers, handlers...)
	g.Route.AddRoute(finalPath, method, finalhandler...)
}
func (g *RouterGroup) GET(path string, handler ...HandleFunc) {
	g.AddRoute(path, "GET", handler...)
}
func (g *RouterGroup) POST(path string, handler ...HandleFunc) {
	g.AddRoute(path, "POST", handler...)
}
func (g *RouterGroup) PUT(path string, handler ...HandleFunc) {
	g.AddRoute(path, "PUT", handler...)
}
func (g *RouterGroup) DELETE(path string, handler ...HandleFunc) {
	g.AddRoute(path, "DELETE", handler...)
}
