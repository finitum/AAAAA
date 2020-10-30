package main

import (
	"github.com/finitum/AAAAA/pkg/store"
	"net/http"
)

type Routes struct {
	db store.Store
}

func NewRoutes(db store.Store) Routes {
	return Routes{db}
}

func (rs Routes) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
