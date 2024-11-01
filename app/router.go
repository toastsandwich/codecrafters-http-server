package main

type Router map[string]Handler

func (r *Router) Handle(path string, h Handler) {
	if *r == nil {
		*r = make(Router)
	}
	if path == "" {
		path = "/"
	}
	(*r)[path] = h
}
