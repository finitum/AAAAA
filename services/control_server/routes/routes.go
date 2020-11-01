package routes

import (
	"github.com/finitum/AAAAA/pkg/auth"
	"github.com/finitum/AAAAA/pkg/store"
	"net/http"
)

type Routes struct {
	db   store.Store
	auth auth.AuthenticationService
}

func New(db store.Store, auth auth.AuthenticationService) *Routes {
	return &Routes{db, auth}
}

func (*Routes) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
