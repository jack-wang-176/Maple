#Maple Web Framework

![Go Version](https://img.shields.io/github/go-mod/go-version/jack-wang-176/Maple)
![License](https://img.shields.io/badge/license-MIT-green)

> **Maple**是一个基于Go语言从零开始构建的轻量级Web框架。
> A light Web Framework built from scratch in Go

本项目旨在通过手写核心逻辑，深入理解Web框架的底层原理，核心抽象包括**Context**和**HttpService**

--------------------------------
##核心抽象和设计思路简介

###HttpService
* **面向用户核心接口**
  * 通过`NewHttpService`方法生成核心web引擎
  * 用户通过该接口实现web框架调用
  * 隐式调用`ServeHTTP`方法，实现业务逻辑拼装搭建
* **和http/net原生包可拆卸绑定**
  * 支持和原生包的拆卸，同时有保证一定满足 `http.Handler`接口
  * ```go
    var _ Service = &HttpService{}
    ```
* **Start方法实现了建立tcp监听**
* **提供GET，POST，DELETE，PUT快捷方式**
* **实现和route层的嵌入绑定**
  * ```go
    type HttpService struct {
	              *Route
        }
    ```
###Context
* **解决了业务函数之间无法进行数据传输问题**‘
  * `http.ResponseWriter`和 `*http.Request`的封装在用户层次实现cs数据交互
  * `index    int`和`handlers []HandleFunc`实现了多个业务函数的组装
  * `param map[string]string`供用户访问动态参数
* **封装了快捷方法**
  * `Param`快捷读取动态参数内容
  * `JsonResp`将输入内容快速转化为json格式快速返回给客户端
  * `BindJson`将输入内容和json格式进行绑定

###Router和Node
* **每一个http方法对应一颗路由树**
* **AddRoute函数实现路由注册**
  * GET等快捷方法套接在AddRoute方法上
  * 通过最开始的判断验证保证树的唯一性
  * 根节点和一般节点被区别对待（对应map逻辑不一样）
* **FindRoute方法实现了业务逻辑搭建**
  * 在NewService层中被调用
  * 定义了`matchInfo`结构体来优雅的兼容动态路径
* **函数功能拆分**
* **Node节点多种类型**

###RouterGroup
* **实现路由组的核心抽象**
  * 实际上是基于Router层的进一步包装，因为AddRouter实际上挂载到Router层，所以可以实现直接调用注册
* **封装多种快捷方式**
  * 封装了GET,DELETE,PUT,POST快捷方式，思路和Service层一致
* **`Group`方法实现优雅的调用**

-----------------------------------------

##快速开始

###1，安装

```bash
go get [github.com/jack-wang-176/Maple](https://github.com/jack-wang-176/Maple)