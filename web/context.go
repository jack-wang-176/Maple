package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Context 是一个proxy持有原始数据，解决数据传输问题
type Context struct {
	Resp    http.ResponseWriter
	Request *http.Request
	//这里放在context里面是无奈之举
	//这里不在这里放handle没法检测是否
	//所有函数都被执行
	index    int
	handlers []HandleFunc
	//封装的动态路由储存的信息
	param map[string]string
}

// Param 这里提供了一个简单的快捷方式来去读取出在Context里面储存的
// 动态匹配内容
func (ctx *Context) Param(key string) (string, error) {
	s, ok := ctx.param[key]
	if !ok {
		return "", errors.New("获取参数路径失败")
	}
	return s, nil
}

// HandleFunc 定义业务逻辑
type HandleFunc func(ctx *Context)

// JsonResp 将输入内容转化为json格返回给客户端
func (ctx *Context) JsonResp(val any) error {
	//这里运用marshal方法转化json格式
	marshal, err := json.Marshal(val)
	if err != nil {
		return errors.New("输入数据无法转化为json格式")
	}
	//这里要设置header告诉浏览器数据类型
	//必须要先写入header再进行数据的发送
	ctx.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	//这里是运用封装的response_writer方法将marshal返回
	//int是写入的输入字节数
	write, err := ctx.Resp.Write(marshal)

	if err != nil {
		return err
	}
	if len(marshal) != write {
		return errors.New("写入数据不等于预期")
	}
	return nil
}

// BindJson 将输入的val和json绑定在一块
func (ctx *Context) BindJson(val any) error {
	if val == nil {
		return errors.New("输入的val不能为空")
	}
	if ctx.Request.Body == nil {
		return errors.New("body不能为空")
	}
	//这里NewDecoder本身返回了decoder,所以这里可以链式调用
	//这里本身调用了json方法，val已经被转化为json格式
	err := json.NewDecoder(ctx.Request.Body).Decode(val)
	if err != nil {
		return err
	}
	return nil
}
