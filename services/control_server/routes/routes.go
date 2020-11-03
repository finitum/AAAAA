package routes

import (
	"github.com/finitum/AAAAA/pkg/auth"
	"github.com/finitum/AAAAA/pkg/executor"
	"github.com/finitum/AAAAA/pkg/store"
	"net/http"
)

type Routes struct {
	db   store.Store
	auth auth.AuthenticationService
	exec executor.Executor
}

func New(db store.Store, auth auth.AuthenticationService, exec executor.Executor) *Routes {
	return &Routes{db, auth, exec}
}

func (*Routes) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
