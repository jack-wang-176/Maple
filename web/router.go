package web

import "strings"

// Route 这里这样分层次设计的好处
// 体现在route对应的是method方法
// node对应seg片段
// 这一点在FindRoute设计的时候也有体现
// 即是根节点要特殊处理
// 其他节点只需要递归即可
type Route struct {
	tree   map[string]*node
	method string
}
type node struct {
	path     string
	children map[string]*node
	//实现多个业务函数的输入
	handler []HandleFunc
}

func NewRoute() *Route {
	return &Route{
		tree: make(map[string]*node),
	}
}

// AddRoute
// 首先/user/login这个url
// / user login 分别处理
// 也就是说第一个参数路径要特殊处理。
// 还有就是说只有/的url也要特殊处理
// AddRoute 这里我们需要注意的是
// 每一次创建都是崭新的开始
// 一种方法只能用一个树
func (r *Route) AddRoute(path string, method string, handler ...HandleFunc) {
	if path == "" {
		panic("path不能为空")
	}
	if path[0] != '/' {
		panic("path必须以/开头")
	}
	if path[len(path)-1] == '/' && path[0] != '/' {
		panic("path不能以/结尾")
	}

	treeRoute, ok := r.tree[method]
	//如果没有该path对应的树直接初始化
	//根节点特殊处理
	if !ok {
		//创建根节点
		//他的路径是/
		treeRoute = &node{
			path:     "/",
			children: make(map[string]*node),
		}
		//这里必须存入
		//否则就在这个if语句中直接销毁
		//根节点特殊处理，有method作为参数路径名
		r.tree[method] = treeRoute
	}
	if path == "/" {
		//这里handler本身就是一个切片
		treeRoute.handler = handler
		//这里后面不需要
		return
	}
	//去掉/
	trim := strings.Trim(path, "/")
	segS := strings.Split(trim, "/")
	//这里构建了具体的一个路由树
	for _, seg := range segS {
		if treeRoute.children[seg] == nil {
			treeRoute.children[seg] = treeRoute.childOrCreateNode(seg)
		}
		//这里不断向下移动
		treeRoute = treeRoute.children[seg]
	}
	//这里是树的顶部
	//我们在树的注册时，一般习惯在根部加上相应的函数
	treeRoute.handler = handler
}

// childOrCreateNode
// 这里是创建childrenNode
func (n *node) childOrCreateNode(seg string) *node {
	theNode, ok := n.children[seg]
	//没有就
	if !ok {
		//handle 的问题在add里面进行调用
		newNode := &node{
			path:     seg,
			children: make(map[string]*node),
		}
		n.children[seg] = newNode
		return newNode
	}
	return theNode
}

// FindRoute 可以看作是AddRoute的逆过程
func (r *Route) findRoute(path string, method string) (*node, bool) {
	treeRoute, ok := r.tree[method]
	if !ok {
		return nil, false
	}
	//这里treeRoute是map匹配返回的节点
	//如果是根节点说明那个node直接就是我们想要的
	if path == "/" {
		return treeRoute, true
	}
	trim := strings.Trim(path, "/")
	segS := strings.Split(trim, "/")
	//这里和FindRoute不同的是不需要考虑到根节点
	//因为进行到这一步说明对应的树一定存在
	//现在就是循环递归seg的步骤
	for _, seg := range segS {
		child, ok := treeRoute.children[seg]
		if !ok {
			return nil, false
		}
		//递归遍历
		treeRoute = child
	}
	//返回到树上最顶端的节点
	return treeRoute, true
}
