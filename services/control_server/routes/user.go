package routes

import (
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/go-chi/render"
	"net/http"
)

func (rs *Routes) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	token, err := rs.auth.Login(&user)
	if err != nil {
		_ = render.Render(w, r, ErrUnauthorized())
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
}

func (rs *Routes) AddUser(w http.ResponseWriter, r *http.Request)  {
	var user models.User

	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.auth.Register(&user); err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}