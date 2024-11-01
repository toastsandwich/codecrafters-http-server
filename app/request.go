package main

import (
	"strings"
)

type URL []string

func url(target string) URL {
	return URL(strings.Split(target, "/")[1:])
}

func (u *URL) Value() string {
	if len(*u) > 1 {
		return (*u)[1]
	}
	return ""
}

type Header map[string]any

func (h *Header) Get(key string) any {
	var val any
	if v, ok := (*h)[key]; ok {
		val = v
	}
	return val
}

func (h *Header) Set(key string, val any) {
	if *h == nil {
		*h = make(Header)
	}
	(*h)[key] = val
}

type HTTPReq struct {
	// Request line
	Method  string
	Target  string
	Version string

	// custom url
	URL
	// Header
	Header
	// Body
	Body string
}

func ParseReq(r string) HTTPReq {
	r_slc := strings.Split(r, "\r\n")
	rReqLine := strings.Split(r_slc[0], " ")
	method, target, version := rReqLine[0], rReqLine[1], rReqLine[2]

	var header Header
	for _, h := range r_slc[1 : len(r_slc)-1] {
		pair := strings.Split(h, ": ")
		if len(pair) == 2 {
			k := strings.TrimSpace(pair[0])
			v := strings.TrimSpace(pair[1])
			header.Set(k, v)
		}
	}
	body := r_slc[len(r_slc)-1]
	return HTTPReq{
		Method:  method,
		Target:  target,
		Version: version,
		URL:     url(target),
		Header:  header,
		Body:    body,
	}
}
