package service

import "context"

type Repo interface {
	CrateAlias(ctx context.Context, original, alias string) error
	GetOrigURL(ctx context.Context, alias string) (string, error)
}

type Service interface {
	CrateAlias(url string) (string, error)
	GetOrigURL(alias string) (string, error)
}
