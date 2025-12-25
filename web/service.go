package web

import (
	"net"
	"net/http"
)

// Service 如果想换其他网络库只需要修改几行代码
// 可以快速地进行测试
// 直接创建空的service，满足相应方法即可
type Service interface {
	http.Handler
	Start(addr string) error
}

// NewHttpService 这个给用户提供了创建一个总引擎的方法
// 这里service层是封装了底层http包
// handler是这个底层包的接口
// 这里HttpService是我的中间件的核心接口
// 既和上层的接口封装到一块（利用go的特性必然包括并且可以拆卸）
// 也为用户使用提供了核心接口
func NewHttpService() *HttpService {
	httpService := &HttpService{
		Route: NewRoute(),
	}
	return httpService
}

// 这里是匿名变量为了实现编译器检查，HttpService必须实现Service接口
var _ Service = &HttpService{}

// HttpService 真正干活的结构体
// route提供了抽象层直接链接的核心接口
type HttpService struct {
	*Route
}

// ServeHTTP 请求进来，封装context，找对应路由，实现路由函数
// 实现了service层业务逻辑的实现
// 业务逻辑组装
func (httpService *HttpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//封装context
	ctx := &Context{
		Resp:    w,
		Request: r,
	}
	//对参数路径进行查询
	//对实际的业务逻辑进行组装
	route, ok := httpService.Route.findRoute(r.URL.Path, r.Method)
	//这里不能世界判断route.handler存在一个空指针保护
	//必须判定route不是空指针
	if !ok || route == nil || route.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("404 page not found"))
		return
	}
	//这个函数就是业务逻辑的实现，所以不用管注册，注册是注册的事情
	//route.handler[](ctx)
	//将handler封装到context层里
	ctx.handlers = route.handler
	ctx.index = -1
	//业务逻辑的执行被放在了这里
	ctx.Next()
}

// Start 创建真实的tcp接口监听
func (httpService *HttpService) Start(addr string) error {
	//创建端口号为addr的tcp端口
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	//这里要传入一个handle
	//httpService满足了handler接口
	return http.Serve(listen, httpService)
}

// AddRoute 作为面向用户的路由树创建
// 这里集成了router层的路由树
func (httpService *HttpService) AddRoute(path string, method string, handler ...HandleFunc) {
	//这里切片被当成一个整体
	//设计的时候要把切片重新打散
	httpService.Route.AddRoute(path, method, handler...)
}
func (httpService *HttpService) GET(path string, handler ...HandleFunc) {
	httpService.Route.AddRoute(path, http.MethodGet, handler...)
}
func (httpService *HttpService) POST(path string, handler ...HandleFunc) {
	httpService.Route.AddRoute(path, http.MethodPost, handler...)
}
func (httpService *HttpService) PUT(path string, handler ...HandleFunc) {
	httpService.Route.AddRoute(path, http.MethodPut, handler...)
}
func (httpService *HttpService) DELETE(path string, handler ...HandleFunc) {
	httpService.Route.AddRoute(path, http.MethodDelete, handler...)
}
