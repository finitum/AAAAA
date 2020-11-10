package routes

import (
	"encoding/json"
	"github.com/finitum/AAAAA/pkg/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (rs *Routes) Login(w http.ResponseWriter, r *http.Request) {
	var user auth.FullUser

	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	token, err := rs.auth.Login(user.Username, user.Password)
	if err != nil {
		_ = render.Render(w, r, ErrUnauthorized())
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	tokenJson := map[string]string{
		"token": token,
	}

	res, err := json.Marshal(tokenJson)
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	_, _ = w.Write(res)
}

func (rs *Routes) AddUser(w http.ResponseWriter, r *http.Request) {
	var user auth.FullUser

	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.auth.Register(user); err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rs *Routes) GetUsers(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := rs.auth.GetUsers()
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Errorf("failed to get users (%v)", err)
		return
	}

	users := make([]render.Renderer, len(dbUsers))
	for i, user := range dbUsers {
		users[i] = user
	}

	_ = render.RenderList(w, r, users)
}

func (rs *Routes) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	allUsers, err := rs.auth.GetUserNames()
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Errorf("failed to get users (%v)", err)
		return
	}
	if len(allUsers) == 0 {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		log.Errorf("invalid request: can't remove last user (%v)", err)
		return
	}

	err = rs.auth.DeleteUser(username)
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Errorf("failed to remove user (%v)", err)
	}

}

func (rs *Routes) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user auth.FullUser

	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	claims, _ := auth.FromContext(r.Context())
	err := rs.auth.Update(user, claims.RawToken)
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
