package network

import (
	"net/http"

	owlsystem "github.com/xssed/owlcache/system"
)

type ServerEntity struct {
	handler http.Handler
}

func (se *ServerEntity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "owlcache "+owlsystem.VERSION)
	se.handler.ServeHTTP(w, r)
}
