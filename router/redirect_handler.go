package router

import (
	"net/http"
)

type RedirectHandler struct {
	url  string
	code int
}

func (r *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, r.url, r.code)
}

func NewRedirectHandler(url string, code int) http.Handler {
	return http.Handler(&RedirectHandler{url: url, code: code})
}
