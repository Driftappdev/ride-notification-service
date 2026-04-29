package httpadapter

import "net/http"

type Router interface {
	Handle(pattern string, handler http.Handler)
}

type Adapter struct {
	router Router
}

func New(router Router) *Adapter { return &Adapter{router: router} }
