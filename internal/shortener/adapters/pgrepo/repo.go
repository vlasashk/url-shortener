package pgrepo

import "context"

func (db PgRepo) CrateAlias(ctx context.Context, url, alias string) error {
	return nil
}
func (db PgRepo) GetOrigURL(ctx context.Context, alias string) (string, error) {
	return "", nil
}
