package models

import "errors"

var (
	ErrAliasCollision = errors.New("alias collision")
	ErrInvalidAlias   = errors.New("invalid alias")
)
