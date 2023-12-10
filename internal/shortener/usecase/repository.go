package usecase

import "context"

type Repo interface {
	CrateAlias(ctx context.Context, url, alias string) error
	GetOrigURL(ctx context.Context, alias string) (string, error)
}
