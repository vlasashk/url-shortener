package service

import (
	"context"
	"github.com/vlasashk/url-shortener/pkg/aliasgen"
	"time"
)

const (
	reqTimeOut = 3 * time.Second
)
const (
	AliasCollisionErr = "alias collision"
)

type URLService struct {
	db Repo
}

func New(db Repo) *URLService {
	return &URLService{db: db}
}

func (s *URLService) CrateAlias(url string) (string, error) {
	alias := aliasgen.Generate()
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeOut)
	defer cancel()
	err := s.db.CrateAlias(ctx, url, alias)
	for err != nil && err.Error() == AliasCollisionErr {
		alias = aliasgen.Generate()
		err = s.db.CrateAlias(ctx, url, alias)
	}
	if err != nil {
		return "", err
	}
	return alias, nil
}
func (s *URLService) GetOrigURL(alias string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeOut)
	defer cancel()
	return s.db.GetOrigURL(ctx, alias)
}
