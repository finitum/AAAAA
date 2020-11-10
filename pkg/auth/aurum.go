package auth

import (
	"github.com/finitum/aurum/clients/go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

func (a *Aurum) Login(user string, pass string) (string, error) {
	tp, err := a.au.Login(user, pass)
	if err != nil {
		return "", err
	}

	return tp.LoginToken, nil
}

func (a *Aurum) Register(user FullUser) error {
	if user.Email == "" {
		user.Email = "no-email"
	}

	err := a.au.Register(user.Username, user.Password, user.Email)
	return errors.Wrap(err, "aurum signup failed")
}

func (a *Aurum) Update(user FullUser, token string) error {
	log.Error("updating user unsupported")
	return errors.New("unsupported")
}

func (a *Aurum) Verify(token string) (ret Claims, _ bool) {
	claims, err := a.au.Verify(token)
	if err != nil {
		return Claims{}, false
	}

	if err := claims.Valid(); err != nil {
		return Claims{}, false
	}


	ret.Username = claims.Username
	ret.RawToken = token
	return ret, !claims.Refresh
}
