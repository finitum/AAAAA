package routes

import (
	"github.com/finitum/AAAAA/pkg/auth"
	"github.com/finitum/AAAAA/pkg/executor"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/finitum/AAAAA/services/control_server/config"
	"net/http"
)

type Routes struct {
	cfg  *config.Config
	db   store.Store
	auth auth.AuthenticationService
	exec executor.Executor
}

func New(cfg *config.Config, db store.Store, auth auth.AuthenticationService, exec executor.Executor) *Routes {
	return &Routes{cfg, db, auth, exec}
}

func (*Routes) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
