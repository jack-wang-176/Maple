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
	//添加动态路由匹配
	//这里我们需要的是模糊匹配
	//map结构体对应的精确匹配难以使用
	paramChild *node
}

// matchInfo 这里必须考虑一下拿到动态路由或者其他之后的返回
// 直接全部展开返回过于丑陋
type matchInfo struct {
	node *node
	Info map[string]string
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
		//if 和 else 里面逻辑是一样的
		//只是针对不同情况的创建
		if seg[0] == ':' {
			if treeRoute.paramChild == nil {
				treeRoute.paramChild = treeRoute.childOrCreateParam(seg)
			}
			treeRoute = treeRoute.paramChild
		} else {
			//总的来捉这里有两种情况。这里我写了不是动态路由的情况
			//正常节点情况
			if treeRoute.children[seg] == nil {
				//这里childOrCreateNode里面主要是一般节点的注册
				treeRoute.children[seg] = treeRoute.childOrCreateNode(seg)
			}
			//这里不断向下移动
			treeRoute = treeRoute.children[seg]
		}
	}
	//这里是树的顶部
	//我们在树的注册时，一般习惯在根部加上相应的函数
	treeRoute.handler = handler
}

// childOrCreateNode
// 这里是创建childrenNode
// 针对普遍路由节点
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

// childOrCreateParam 针对动态匹配节点
// 因为设计的原因这里我将param节点和正常节点的创建趋同化了
func (n *node) childOrCreateParam(seg string) *node {

	if n.paramChild == nil {
		newNode := &node{
			path:     seg,
			children: make(map[string]*node),
		}
		return newNode
	}
	return n.paramChild
}

// FindRoute 可以看作是AddRoute的逆过程
// 这里要保证找到对应map后存入context
func (r *Route) findRoute(path string, method string) (matchInfo, bool) {
	treeRoute, ok := r.tree[method]
	if !ok {
		return matchInfo{
			node: nil,
		}, false
	}
	//这里treeRoute是map匹配返回的节点
	//如果是根节点说明那个node直接就是我们想要的
	if path == "/" {
		return matchInfo{
			node: treeRoute,
		}, true
	}
	trim := strings.Trim(path, "/")
	segS := strings.Split(trim, "/")
	//这里和FindRoute不同的是不需要考虑到根节点
	//因为进行到这一步说明对应的树一定存在
	//现在就是循环递归seg的步骤
	paramInfo := make(map[string]string)
	for _, seg := range segS {
		theNode, is, isParam := treeRoute.detectNodeType(seg)

		if !is {
			return matchInfo{}, false
		}
		if isParam {
			//当动态路由的时候
			//传入 user/100
			//去掉：
			key := theNode.path[1:]
			//这里我们要将这里的模糊储存存在matchInfo中
			//尽管本身这里的输入是模糊的
			//但是储存在map里面最后会落到context里面
			//进行了严格匹配
			paramInfo[key] = seg

		}
		treeRoute = theNode
	}
	//返回到树上最顶端的节点
	return matchInfo{
		node: treeRoute,
		Info: paramInfo,
	}, true
}

// detectNodeType 这个是供给FindRoute使用 第一个bool是节点是否存在，第二个bool是是否是动态节点
func (n *node) detectNodeType(seg string) (*node, bool, bool) {
	if n.children != nil {
		theNode, ok := n.children[seg]
		//这里需要考虑动态路由和静态路由都存在的情况
		//所以说这里的if判断应该是嵌套逻辑
		if !ok {
			if n.paramChild != nil {
				return n.paramChild, true, true
			}
			return nil, false, false
		}
		return theNode, true, false
	}
	//考虑只有动态节点的状态
	if n.paramChild != nil {
		return n.paramChild, true, true
	}
	//既不存在静态路由也不存在动态路由
	return nil, false, false
}
