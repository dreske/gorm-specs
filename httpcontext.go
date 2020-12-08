package go_specs

import "net/http"

type httpContext struct {
	request *http.Request
}

func (h *httpContext) GetParameter(name string) ([]string, bool) {
	values, ok := h.request.URL.Query()[name]
	return values, ok
}

func HttpContext(r *http.Request) Context {
	return &httpContext{
		request: r,
	}
}
