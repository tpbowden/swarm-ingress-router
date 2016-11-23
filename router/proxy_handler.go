package router

import (
	"github.com/valyala/fasthttp"
)

// RedirectHandler stores the URL and HTTP code of a redirect
type ProxyHandler struct {
	url string
}

// ServerHTTP writes a redirect to an HTTP response
func (r *ProxyHandler) ServeHTTP(ctx *fasthttp.RequestCtx) {
	proxyClient := &fasthttp.HostClient{
		Addr: r.url,
		// set other options here if required - most notably timeouts.
	}

	req := &ctx.Request
	resp := &ctx.Response
	if err := proxyClient.Do(req, resp); err != nil {
		ctx.Logger().Printf("error when proxying the request: %s", err)
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
	}
}

// NewRedirectHandler returns a new redirect handler
func NewProxyHandler(url string) fasthttp.RequestHandler {
	handler := ProxyHandler{url: url}
	return fasthttp.RequestHandler(handler.ServeHTTP)
}
