package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener/adapters/pgrepo/converter"
	"github.com/vlasashk/url-shortener/internal/shortener/models"
	"github.com/vlasashk/url-shortener/pkg/migration"
	"github.com/vlasashk/url-shortener/pkg/pgconnect"
)

const (
	keyCollisionCode = "23505"
)

const (
	createQry       = `INSERT INTO url (alias, original, expires_at, visits) VALUES ($1, $2, $3, $4)`
	getAliasQry     = `SELECT original FROM url WHERE alias = $1`
	updateVisitsQry = `UPDATE url SET visits = visits + 1 WHERE alias = $1`
)

type PgRepo struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.PostgresCfg, logger zerolog.Logger) (*PgRepo, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.NameDB)

	dbPool, err := pgconnect.Connect(ctx, url, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err = migration.Up(dbPool, cfg.Migrations); err != nil {
		return nil, fmt.Errorf("unable to apply migrations: %w", err)
	}

	return &PgRepo{dbPool}, nil
}

func (db *PgRepo) SaveAlias(ctx context.Context, original, alias string) error {
	newURL := converter.New(alias, original)

	if _, err := db.Pool.Exec(ctx, createQry, newURL.Alias, newURL.Original, newURL.ExpiresAt, newURL.Visits); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return errorHandler(pgErr)
		}
		return fmt.Errorf("exec transaction fail: %w", err)
	}

	return nil
}
func (db *PgRepo) GetOrigURL(ctx context.Context, alias string) (string, error) {
	var originalURL string

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin transaction fail: %w", err)
	}
	defer func() {
		txFinisher(ctx, tx, err)
	}()

	err = tx.QueryRow(ctx, getAliasQry, alias).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrInvalidAlias
		}
		return "", fmt.Errorf("query execution fail: %w", err)
	}

	_, err = tx.Exec(ctx, updateVisitsQry, alias)
	if err != nil {
		return "", fmt.Errorf("query execution fail: %w", err)
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
