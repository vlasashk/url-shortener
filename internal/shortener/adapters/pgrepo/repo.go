package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"github.com/vlasashk/url-shortener/internal/shortener/adapters/pgrepo/converter"
	"github.com/vlasashk/url-shortener/internal/shortener/models"
)

const (
	keyCollisionCode = "23505"
)

const (
	createQry       = `INSERT INTO url (alias, original, expires_at, visits) VALUES ($1, $2, $3, $4)`
	getAliasQry     = `SELECT original FROM url WHERE alias = $1`
	updateVisitsQry = `UPDATE url SET visits = visits + 1 WHERE alias = $1`
)

func (db *PgRepo) SaveAlias(ctx context.Context, original, alias string) error {
	newURL := converter.New(alias, original)

	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction fail: %v", err)
	}
	defer func() {
		txFinisher(ctx, tx, err)
	}()

	if _, err = tx.Exec(ctx, createQry, newURL.Alias, newURL.Original, newURL.ExpiresAt, newURL.Visits); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return errorHandler(pgErr)
		}
		return fmt.Errorf("exec transaction fail: %v", err)
	}

	return nil
}
func (db *PgRepo) GetOrigURL(ctx context.Context, alias string) (string, error) {
	var originalURL string

	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin transaction fail: %v", err)
	}
	defer func() {
		txFinisher(ctx, tx, err)
	}()

	err = tx.QueryRow(ctx, getAliasQry, alias).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrInvalidAlias
		}
		return "", fmt.Errorf("query execution fail: %v", err)
	}
	_, err = tx.Exec(ctx, updateVisitsQry, alias)
	if err != nil {
		return "", fmt.Errorf("query execution fail: %v", err)
	}
	return originalURL, nil
}

func errorHandler(pgErr *pgconn.PgError) error {
	switch pgErr.Code {
	case keyCollisionCode:
		return models.ErrAliasCollision
	default:
		return pgErr
	}
}

func txFinisher(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Error().Err(err).Msg("transaction rollback for register fail")
		}
	} else {
		err = tx.Commit(ctx)
		if err != nil {
			log.Error().Err(err).Msg("transaction commit for register fail")
		}
	}
}
