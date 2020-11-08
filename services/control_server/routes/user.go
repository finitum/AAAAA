package routes

import (
	"encoding/json"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
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

func (rs *Routes) GetUsers(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := rs.db.AllUsers()
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Errorf("failed to get users (%v)", err)
	}

	users := make([]render.Renderer, len(dbUsers))
	for i, user := range dbUsers {
		// No need to send these
		user.Password = ""
		users[i] = user
	}

	_ = render.RenderList(w, r, users)
}

func (rs *Routes) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	allUsers, err := rs.db.AllUserNames()
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Errorf("failed to get users (%v)", err)
	}
	if len(allUsers) == 0 {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		log.Errorf("invalid request: can't remove last user (%v)", err)
	}

	err = rs.db.DelUser(username)
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Errorf("failed to remove user (%v)", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (rs *Routes) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	rs.auth.Update(&user)
	w.WriteHeader(http.StatusCreated)
}
