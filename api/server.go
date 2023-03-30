package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	addr string
}

func New(addr string) *Api {
	return &Api{
		addr: addr,
	}
}

func (a *Api) Run() error {
	r := mux.NewRouter()

	sub := r.PathPrefix("/api").Subrouter()

	sub.HandleFunc("/auth/google/login", handleGoogleLogin)
	sub.HandleFunc("/auth/google/callback", handleGoogleCallback)

	log.Printf("Api running on addr %s", a.addr)

	ser := &http.Server{
		Addr:    a.addr,
		Handler: r,
	}

	return ser.ListenAndServe()
}
