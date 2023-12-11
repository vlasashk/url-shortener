package memrepo

import (
	"context"
	"errors"
	"github.com/vlasashk/url-shortener/internal/shortener/models/service"
	"sync"
)

type MemRepo struct {
	DB map[string]string
	mu *sync.RWMutex
}

func New() *MemRepo {
	return &MemRepo{
		DB: make(map[string]string),
		mu: &sync.RWMutex{},
	}
}

func (db *MemRepo) CrateAlias(ctx context.Context, original, alias string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.DB[alias]; ok {
		return errors.New(service.AliasCollisionErr)
	}
	db.DB[alias] = original
	return nil
}
func (db *MemRepo) GetOrigURL(ctx context.Context, alias string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	original, ok := db.DB[alias]
	if !ok {
		return "", errors.New(service.InvalidAliasErr)
	}
	return original, nil
}
