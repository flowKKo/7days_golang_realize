package gee

import "strings"

// use trie tree to realize dynamic router

type node struct {
	pattern  string  // the router wait for match
	part     string  // part of router
	children []*node // current node's children node
	isWild   bool    // whether it is an accurate match, if current part is :filename or *filename, then isWild = true
}

// get the first matching node
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// get all the matching nodes
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// use iteration to register a new router
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height{
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)

	// if the first char of the part is ':', that means it is a dynamic router
	if child == nil{
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height + 1)
}

// use iteration to get the router result
func (n* node) search(parts []string, height int) *node{
	if len(parts) == height || strings.HasPrefix(n.part, "*"){
		if n.pattern == ""{
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children{
		result := child.search(parts, height + 1)
		if result != nil{
			return result
		}
	}
	return nil
}

