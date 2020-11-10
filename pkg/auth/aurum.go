package auth

import (
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/aurum/clients/go"
	"github.com/pkg/errors"
)

type Aurum struct {
	au *aurum.Aurum
}

func NewAurum(url string) (*Aurum, error) {
	au, err := aurum.Connect(url)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to aurum failed")
	}

	return &Aurum{ au }, nil
}

func (a *Aurum) Login(user *models.User) (string, error) {
	tp, err := a.au.Login(user.Username, user.Password)
	if err != nil {
		return "", err
	}

	return tp.LoginToken, nil
}

func (a *Aurum) Register(user *models.User) error {
	return errors.Wrap(a.au.Register(user.Username, user.Password, "nomail@AAAAA"), "aurum signup failed")
}

func (a *Aurum) Update(user *models.User, token string) error {
	return errors.New("unsupported")
}

func (a *Aurum) Verify(token string) (Claims, bool) {
	claims, err := a.au.Verify(token)
	if err != nil {
		return Claims{}, false
	}

	if err := claims.Valid(); err != nil {
		return Claims{}, false
	}

	return Claims{
		Username: claims.Username,
		RawToken: token,
	}, !claims.Refresh
}
