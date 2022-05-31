package gee

import (
	"log"
	"net/http"
	"strings"
)

// roots key eg, roots['GET'], roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
// parsePattern is used to parse the url input into single variant
// the content of pattern is like "/v2/hello/:name", so we will get the result as {v2, hello, :name}
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)

	parts := parsePattern(pattern)
	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}


// getRouter is used to get the matched trie tree node and get the mapping of params->pattern
func (r *router) getRoute(method string, path string) (*node, map[string]string) {

	// parse the input url
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {


		parts := parsePattern(n.pattern)

		// matching the dynamic router
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	// firstly, get the trie tree node and params mapping
	n, params := r.getRoute(c.Method, c.Path)
	// if the pattern exists
	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		// add the pattern process function into handler queue
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// if we have added middleware to this pattern, calling c.Next() will firstly use
	// middleware function to process and then use the pattern method
	c.Next()
}


