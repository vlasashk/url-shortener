package usecase

import (
	"context"
	"errors"
	"path"
	"time"

	"github.com/vlasashk/url-shortener/internal/shortener/models"
	"github.com/vlasashk/url-shortener/pkg/aliasgen"
)

const (
	reqTimeOut = 3 * time.Second
)

//go:generate mockery --name=AliasSaver
type AliasSaver interface {
	SaveAlias(ctx context.Context, original, alias string) error
}

//go:generate mockery --name=OriginalProvider
type OriginalProvider interface {
	GetOrigURL(ctx context.Context, alias string) (string, error)
}

type UseCase struct {
	address          string
	aliasSaver       AliasSaver
	originalProvider OriginalProvider
}

func New(address string,
	aliasSaver AliasSaver,
	originalProvider OriginalProvider) *UseCase {
	return &UseCase{
		address:          address,
		aliasSaver:       aliasSaver,
		originalProvider: originalProvider,
	}
}

func (s *UseCase) CreateAlias(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, reqTimeOut)
	defer cancel()

	alias := aliasgen.Generate()

	err := s.aliasSaver.SaveAlias(ctx, url, alias)
	for err != nil && errors.Is(err, models.ErrAliasCollision) {
		alias = aliasgen.Generate()
		err = s.aliasSaver.SaveAlias(ctx, url, alias)
	}

	if err != nil {
		return "", err
	}

	return path.Join(s.address, alias), nil
}

func (s *UseCase) GetOrigURL(ctx context.Context, alias string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, reqTimeOut)
	defer cancel()

	return s.originalProvider.GetOrigURL(ctx, alias)
}
