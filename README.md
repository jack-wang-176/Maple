#Maple Web Framework

![Go Version](https://img.shields.io/github/go-mod/go-version/jack-wang-176/Maple)
![License](https://img.shields.io/badge/license-MIT-green)

> **Maple**是一个基于Go语言从零开始构建的轻量级Web框架。
> A light Web Framework built from scratch in Go

本项目旨在通过手写核心逻辑，深入理解Web框架的底层原理，核心特性包括**前缀树路由**和**上下文封装**

* **前缀树路由**
* 采用Trie树结构实现高效的路由匹配
* **链式中间件调用**
  *使用Next()方法实现洋葱模型
* **路由分组**:
  *支持层级分组管理，共享前缀和中间件。
* **Context**
  *封装了 `http.Request` 和`http.RespponseWriter`
  *提供了几种快捷方法
* **易于扩展**
  *`HttpService`实现了标准库`http.Handler`接口

##快速开始

###1，安装

```bash
go get [github.com/jack-wang-176/Maple](https://github.com/jack-wang-176/Maple)