package store

import "github.com/finitum/AAAAA/pkg/aur"

type AurCache interface {
	SetEntry(searchterm string, result aur.Results) (error)
	GetEntry(searchterm string) (aur.Results, error)
}

type Cache interface {
	AurCache
}
