package auth

import (
	"context"
	"github.com/finitum/AAAAA/pkg/models"
	"net/http"
	"strings"
)

type Claims struct {
	Username string
	RawToken string
}

const VerifierContextKey = "AAAAA jwt verification key"

type AuthenticationService interface {
	Login(user *models.User) (string, error)
	Register(user *models.User) error
	Update(user *models.User, token string) error

	// Verify verifies a jwt token returns claims, true if success nil, false otherwise
	Verify(token string) (Claims, bool)
}

// VerificationMiddleware calls AuthenticationService.Verify, 401s on failure and puts the claims in the request context
// on success, these can be retrieved with FromContext
func VerificationMiddleware(a AuthenticationService) func (next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			auth = strings.TrimPrefix(auth, "Bearer ")

			claims, valid := a.Verify(auth)
			if !valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims.RawToken = auth

			ctx := context.WithValue(r.Context(), VerifierContextKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// FromContext retrieves the Claims from a request context, returns Claims, true on success nil, false otherwise
func FromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(VerifierContextKey).(Claims)
	if !ok {
		return Claims{}, false
	}
	return claims, true
}
