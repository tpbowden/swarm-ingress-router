package router

import (
	"github.com/valyala/fasthttp"
)

// RedirectHandler stores the URL and HTTP code of a redirect
type RedirectHandler struct {
	url  string
	code int
}

// ServerHTTP writes a redirect to an HTTP response
func (r *RedirectHandler) ServeHTTP(ctx *fasthttp.RequestCtx) {
	ctx.Redirect(r.url, r.code)
}

// NewRedirectHandler returns a new redirect handler
func NewRedirectHandler(url string, code int) fasthttp.RequestHandler {
	handler := RedirectHandler{url: url, code: code}
	return fasthttp.RequestHandler(handler.ServeHTTP)
}
