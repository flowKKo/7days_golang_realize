package gee

import (
	"net/http"
)

type HandlerFunc func( *Context)

type Engine struct{
	router *router
}

func New() *Engine{
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc){
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func(engine *Engine)GET(pattern string, handler HandlerFunc){
	engine.router.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func(engine *Engine)POST(pattern string, handler HandlerFunc){
	engine.router.addRoute("POST", pattern, handler)
}

func(engine *Engine)Run (addr string)(err error){
	return http.ListenAndServe(addr, engine)
}

func(engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	c := newContext(w, req)
	engine.router.handle(c)
}