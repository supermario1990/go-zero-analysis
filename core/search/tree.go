// 使用二叉查找树来实现路由查找
package search

import "errors"

const (
	colon = ':'
	slash = '/'
)

var (
	ErrDupItem      = errors.New("duplicated item")
	ErrDupSlash     = errors.New("duplicated slash")
	ErrEmptyItem    = errors.New("empty item")
	ErrInvalidState = errors.New("search tree is in an invalid state")
	ErrNotFromRoot  = errors.New("path should start with /")

	NotFound Result
)

type (
	// 内部使用结构体
	innerResult struct {
		key   string
		value string
		named bool
		found bool
	}

	// 节点
	node struct {
		// 元素, handler
		item     interface{}
		// 孩子节点, key 为'/'分割的字符串， value为孩子节点
		// 0 存储一般节点/ 1 存储带冒号节点:
		children [2]map[string]*node
	}

	// 树
	Tree struct {
		root *node	// 根节点
	}

	// 外部使用结构体
	Result struct {
		Item   interface{}
		Params map[string]string
	}
)

// 创建树
func NewTree() *Tree {
	return &Tree{
		root: newNode(nil),
	}
}

// 添加一个路径
func (t *Tree) Add(route string, item interface{}) error {
	// 路径不为空，且必须以‘/’开头
	if len(route) == 0 || route[0] != slash {
		return ErrNotFromRoot
	}

	if item == nil {
		return ErrEmptyItem
	}

	return add(t.root, route[1:], item)
}

// 匹配搜索路径
func (t *Tree) Search(route string) (Result, bool) {
	if len(route) == 0 || route[0] != slash {
		return NotFound, false
	}

	var result Result
	ok := t.next(t.root, route[1:], &result)
	return result, ok
}

func (t *Tree) next(n *node, route string, result *Result) bool {
	if len(route) == 0 && n.item != nil {
		result.Item = n.item
		return true
	}

	for i := range route {
		if route[i] == slash {
			token := route[:i]
			for _, children := range n.children {
				for k, v := range children {
					if r := match(k, token); r.found {
						if t.next(v, route[i+1:], result) {
							if r.named {
								addParam(result, r.key, r.value)
							}

							return true
						}
					}
				}
			}

			return false
		}
	}

	for _, children := range n.children {
		for k, v := range children {
			if r := match(k, route); r.found && v.item != nil {
				result.Item = v.item
				if r.named {
					addParam(result, r.key, r.value)
				}

				return true
			}
		}
	}

	return false
}

// 获取子节点
func (nd *node) getChildren(route string) map[string]*node {
	if len(route) > 0 && route[0] == colon {
		return nd.children[1]
	} else {
		return nd.children[0]
	}
}

func add(nd *node, route string, item interface{}) error {
	if len(route) == 0 {	// 处理最后一个字符是'/'的情况
		if nd.item != nil {
			return ErrDupItem
		}

		nd.item = item
		return nil
	}

	if route[0] == slash {
		return ErrDupSlash
	}

	// 递归处理
	for i := range route {
		if route[i] == slash {
			token := route[:i]
			children := nd.getChildren(token)
			if child, ok := children[token]; ok {
				if child != nil {
					return add(child, route[i+1:], item)
				} else {
					return ErrInvalidState
				}
			} else {
				child := newNode(nil)
				children[token] = child
				return add(child, route[i+1:], item)
			}
		}
	}

	// 处理URI最后一个的字字符串（以'/'分割）
	children := nd.getChildren(route)
	if child, ok := children[route]; ok {
		if child.item != nil {
			return ErrDupItem
		}

		child.item = item
	} else {
		children[route] = newNode(item)
	}

	return nil
}

func addParam(result *Result, k, v string) {
	if result.Params == nil {
		result.Params = make(map[string]string)
	}

	result.Params[k] = v
}

func match(pat, token string) innerResult {
	if pat[0] == colon {
		return innerResult{
			key:   pat[1:],
			value: token,
			named: true,
			found: true,
		}
	}

	return innerResult{
		found: pat == token,
	}
}

// 创建新的节点, 初始化节点元素和孩子节点
func newNode(item interface{}) *node {
	return &node{
		item: item,
		children: [2]map[string]*node{
			make(map[string]*node),
			make(map[string]*node),
		},
	}
}
