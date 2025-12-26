<div align="center">

# 🍁 Maple Web Framework

![Go Version](https://img.shields.io/github/go-mod/go-version/jack-wang-176/Maple)
![License](https://img.shields.io/badge/license-MIT-green)

> **Maple** 是一个基于 Go 语言从零开始构建的轻量级 Web 框架。
> <br>
> *A light Web Framework built from scratch in Go.*

<p align="left">
本项目旨在通过手写核心逻辑，深入理解 Web 框架的底层原理。核心抽象包括 <b>Context</b> 和 <b>HttpService</b>。
</p>

</div>

---

## 🛠 核心抽象与设计思路

### 1. HttpService (核心引擎)
这是面向用户的核心接口，负责通过 `NewHttpService` 方法生成核心 Web 引擎。

* **业务逻辑拼装**：隐式调用 `ServeHTTP` 方法，实现业务逻辑的拼装与搭建。
* **原生兼容性**：支持与 `net/http` 原生包的拆卸绑定，同时保证满足 `http.Handler` 接口：
    ```go
    var _ Service = &HttpService{}
    ```
* **TCP 监听**：`Start` 方法内部实现了建立 TCP 监听。
* **路由快捷方式**：提供 `GET`, `POST`, `DELETE`, `PUT` 等常用方法的注册入口。
* **路由层嵌入**：
    ```go
    type HttpService struct {
        *Route
    }
    ```

### 2. Context (上下文)
解决了业务函数之间无法进行数据传输的问题，是连接 User 层与 Framework 层的桥梁。

* **请求与响应封装**：对 `http.ResponseWriter` 和 `*http.Request` 进行封装，实现 CS 数据交互。
* **处理链组装**：通过 `index int` 和 `handlers []HandleFunc` 实现多个业务函数（中间件+处理函数）的有序调用。
* **参数访问**：提供 `param map[string]string` 供用户访问动态路由参数。
* **快捷方法**：
    * `Param`：快捷读取动态参数内容。
    * `JsonResp`：将输入内容快速转化为 JSON 格式并返回给客户端。
    * `BindJson`：将 Request Body 内容绑定到结构体。

### 3. Router 和 Node (路由树)
* **多树结构**：每一个 HTTP 方法（GET, POST 等）对应一棵独立的路由树（Trie）。
* **路由注册 (`AddRoute`)**：
    * GET 等快捷方式底层均调用 `AddRoute`。
    * 通过前置判断保证路由树的唯一性。
    * 根节点和普通节点逻辑分离，分别处理。
* **路由查找 (`FindRoute`)**：
    * 在 `NewService` 层被调用，用于匹配请求。
    * 定义了 `matchInfo` 结构体，优雅地兼容动态路径匹配。
* **节点设计**：`Node` 节点支持多种类型，功能拆分清晰。

### 4. RouterGroup (路由组)
* **逻辑分组**：基于 Router 层的进一步封装，`AddRouter` 实际上是挂载到 Router 层。
* **优雅调用**：
    * 支持 `Group` 方法创建子分组。
    * 封装了 `GET`, `DELETE`, `PUT`, `POST` 快捷方式，保持与 Service 层一致的体验。

---

## 🚀 快速开始 (Quick Start)

### 1. 安装

```bash
go get [github.com/jack-wang-176/Maple](https://github.com/jack-wang-176/Maple)
```

### 2. 使用示例

以下代码展示了如何启动一个 Maple 服务，处理 JSON 请求以及获取动态路由参数。

```go
package main

import (
	"net/http"
	"[github.com/jack-wang-176/Maple](https://github.com/jack-wang-176/Maple)" // 确保引用路径正确
)

// User 用于测试 BindJson
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// 1. 初始化引擎
	s := maple.NewHttpService()

	// 2. 注册基础路由
	s.GET("/", func(c *maple.Context) {
		c.JsonResp(http.