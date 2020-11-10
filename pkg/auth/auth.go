package auth

import (
	"context"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const VerifierContextKey = "AAAAA jwt verification key"

type FullUser struct {
	models.User
	Password string
	Email string		`json:"-,omitempty"`
}

func (u *FullUser) Bind(*http.Request) error {
	if u.Username == "" || u.Password == "" {
		return errors.New("invalid user")
	}

	return nil
}

type Claims struct {
	models.User
	RawToken string
}

type AuthenticationService interface {
	// Login should login a user and return a token
	Login(user string, pass string) (string, error)
	// Register should create a user in the service
	Register(user FullUser) error
	// Update should update the email and password on the service
	Update(user FullUser, token string) error

	// Verify verifies a jwt token returns claims, true if success nil, false otherwise
	Verify(token string) (Claims, bool)
}

// FromContext retrieves the Claims from a request context, returns Claims, true on success nil, false otherwise
func FromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(VerifierContextKey).(Claims)
	if !ok {
		return Claims{}, false
	}
	return claims, true
}

/* Authenticator */

type Authenticator struct {
	as AuthenticationService
	us store.UserStore
}

func NewAuthenticator(as AuthenticationService, us store.UserStore) *Authenticator {
	return &Authenticator{as, us}
}

func (a Authenticator) Login(username, password string) (string, error) {
	if _, err := a.us.GetUser(username); err != nil {
		return "", err
	}

	return a.as.Login(username, password)
}

func (a Authenticator) Register(user FullUser) error {
	if err := a.as.Register(user); err != nil {
		return err
	}

	if err := a.us.AddUser(&user.User); err != nil {
		return err
	}

	return nil
}

func (a Authenticator) Update(user FullUser, token string) error {
	if _, err := a.us.GetUser(user.Username); err != nil {
		return err
	}

	return a.as.Update(user, token)
}

func (a Authenticator) GetUsers() ([]*models.User, error) {
	return a.us.AllUsers()
}

func (a Authenticator) GetUserNames() ([]string, error) {
	return a.us.AllUserNames()
}

func (a Authenticator) DeleteUser(username string) error {
	return a.us.DelUser(username)
}

// VerificationMiddleware calls AuthenticationService.Verify, 401s on failure and puts the claims in the request context
// on success, these can be retrieved with FromContext
func (a Authenticator)VerificationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		auth = strings.TrimPrefix(auth, "Bearer ")

		claims, valid := a.as.Verify(auth)
		if !valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims.RawToken = auth

		ctx := context.WithValue(r.Context(), VerifierContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
