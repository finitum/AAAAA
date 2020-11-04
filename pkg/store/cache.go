package store

import "github.com/finitum/AAAAA/pkg/aur"

type Cache interface {
	SetEntry(searchterm string, result aur.Results) error
	GetEntry(searchterm string) (aur.Results, bool, error)
}
