package auth

import "github.com/finitum/AAAAA/pkg/models"

type AuthenticationService interface {
	Login(user *models.User) (string, error)
	Register(user *models.User) error
}
