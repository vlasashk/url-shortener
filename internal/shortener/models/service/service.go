package service

import (
	"context"
	"github.com/vlasashk/url-shortener/pkg/aliasgen"
	"path"
	"time"
)

const (
	reqTimeOut = 3 * time.Second
)
const (
	AliasCollisionErr = "alias collision"
	InvalidAliasErr   = "alias doesn't exist"
)

type URLService struct {
	DB      Repo
	Address string
}

func New(db Repo, address string) *URLService {
	return &URLService{DB: db, Address: address}
}

func (s *URLService) CrateAlias(url string) (string, error) {
	alias := aliasgen.Generate()
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeOut)
	defer cancel()
	err := s.DB.CrateAlias(ctx, url, alias)
	for err != nil && err.Error() == AliasCollisionErr {
		alias = aliasgen.Generate()
		err = s.DB.CrateAlias(ctx, url, alias)
	}
	if err != nil {
		return "", err
	}

	return path.Join(s.Address, alias), nil
}
func (s *URLService) GetOrigURL(alias string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeOut)
	defer cancel()
	return s.DB.GetOrigURL(ctx, alias)
}
