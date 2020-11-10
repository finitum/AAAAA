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
	auth *auth.Authenticator
	db   store.PackageStore
	exec executor.Executor
}

func New(cfg *config.Config, db store.PackageStore, auth *auth.Authenticator, exec executor.Executor) *Routes {
	return &Routes{cfg, auth, db, exec}
}

func (*Routes) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
