package main

import (
	"strings"
)

type HTTPReq struct {
	// Request line
	Method  string
	Target  string
	Version string
	// Headers
	Headers map[string]string
	// Body
	Body string
}

func ParseReq(r string) HTTPReq {
	r_slc := strings.Split(r, "\r\n")
	rReqLine := strings.Split(r_slc[0], " ")
	method, target, version := rReqLine[0], rReqLine[1], rReqLine[2]

	rmap := make(map[string]string)
	for _, h := range r_slc[1 : len(r_slc)-1] {
		pair := strings.Split(h, ": ")
		if len(pair) == 2 {
			k := strings.TrimSpace(pair[0])
			v := strings.TrimSpace(pair[1])
			rmap[k] = v
		}
	}
	body := r_slc[len(r_slc)-1]
	return HTTPReq{
		Method:  method,
		Target:  target,
		Version: version,
		Body:    body,
		Headers: rmap,
		// Body:    body,
	}
}

func (r HTTPReq) URL() string {
	return r.Target
}
