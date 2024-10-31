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
	Headers string
	// Body
	Body string
}

func ParseReq(r string) HTTPReq {
	r_slc := strings.Split(r, "\r\n")
	rReqLine := strings.Split(r_slc[0], " ")
	method, target, version := rReqLine[0], rReqLine[1], rReqLine[2]

	rHeaders := r_slc[1]
	var body string
	if len(r_slc) > 2 {
		body = r_slc[2]
	}
	return HTTPReq{
		Method:  method,
		Target:  target,
		Version: version,

		Headers: rHeaders,
		Body:    body,
	}
}

func (r HTTPReq) URL() string {
	return r.Target
}
