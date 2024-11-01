package main

type Router map[string]HandlerFunc

func (r *Router) Handle(path string, h HandlerFunc) {
	if *r == nil {
		*r = make(Router)
	}
	(*r)[path] = h
}

