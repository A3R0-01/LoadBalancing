package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend interface {
	SetAlive(bool)
	IsAlive() bool
	GetUrl() *url.URL
	GetActiveConnections() int
	Serve(http.ResponseWriter, *http.Request)
}
type backend struct {
	url          *url.URL
	isAlive      bool
	mux          sync.RWMutex
	connections  int
	reverseProxy *httputil.ReverseProxy
}
