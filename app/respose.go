package main

import (
	"fmt"
	"strings"
)

// http response struct
type HTTPResp struct {
	Version string
	Status  int
	Phrase  string

	Headers map[string]string

	Body string
}

// \r\n is crlf which marks an end
func (h HTTPResp) Format() string {
	// Use strings.Builder for more efficient string concatenation
	var builder strings.Builder

	// Write status line
	builder.WriteString(fmt.Sprintf("%s %d %s\r\n", h.Version, h.Status, h.Phrase))

	// Write headers
	if len(h.Headers) > 0 {
		for k, v := range h.Headers {
			builder.WriteString(k)
			builder.WriteString(": ")
			builder.WriteString(v)
			builder.WriteString("\r\n")
		}
	}

	// Add empty line between headers and body
	builder.WriteString("\r\n")

	// Add body if it exists
	if h.Body != "" {
		builder.Write([]byte(h.Body))
	}

	return builder.String()
}

func (h *HTTPResp) SetHeader(key string, value any) {
	if h.Headers == nil {
		h.Headers = make(map[string]string)
	}
	h.Headers[key] = fmt.Sprintf("%v", value)
}
