package router

import (
	"net/http"
)

// RedirectHandler stores the URL and HTTP code of a redirect
type RedirectHandler struct {
	url  string
	code int
}

// ServerHTTP writes a redirect to an HTTP response
func (r *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, r.url, r.code)
}

// NewRedirectHandler returns a new redirect handler
func NewRedirectHandler(url string, code int) http.Handler {
	return http.Handler(&RedirectHandler{url: url, code: code})
}
