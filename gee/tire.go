package gee

import (
	"fmt"
	"strings"
)

// 定义 trie 树的节点结构体
type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

// 实现 fmt.Stringer 接口的函数，返回该节点的字符串描述
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 向 trie 树中插入一个路由模式与处理函数之间的映射关系
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果当前位置已经到了路由模式的末尾，将该模式保存在当前节点上
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 获取当前位置对应的部分
	part := parts[height]

	// 查找匹配的子节点，如果不存在则新建子节点
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// 递归向下插入
	child.insert(pattern, parts, height+1)
}

// 在 trie 树中搜索与指定 URL 匹配的路由模式，并返回其对应的节点
func (n *node) search(parts []string, height int) *node {
	// 如果已经到达 URL 的末尾，或者当前节点是以 * 开头的通配符节点，则返回当前节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 获取当前位置对应的部分
	part := parts[height]

	// 遍历当前节点的子节点列表，寻找匹配的子节点，递归向下搜索
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 遍历 trie 树，将所有包含路由模式的节点添加到指定的列表中
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}

	for _, child := range n.children {
		child.travel(list)
	}
}

// 在当前节点的子节点列表中，寻找第一个匹配 part 参数的子节点，或第一个通配符子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
