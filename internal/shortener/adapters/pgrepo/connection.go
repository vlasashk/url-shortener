package pgrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener/models/service"
	"os"
	"time"
)

type PgRepo struct {
	db *pgxpool.Pool
}

func New(cfg config.PostgresCfg) (service.Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.NameDB)
	dbPool, err := pgxpool.New(ctx, url)
	if err != nil {
		return PgRepo{}, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err = dbPool.Ping(ctx); err != nil {
		return PgRepo{}, fmt.Errorf("unable to ping connection pool: %v", err)
	}
	instance := PgRepo{dbPool}
	if err = instance.NewTable(ctx, cfg); err != nil {
		return PgRepo{}, err
	}
	return instance, nil
}

func (db PgRepo) NewTable(ctx context.Context, cfg config.PostgresCfg) error {
	err := db.db.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		query, err := os.ReadFile(cfg.InitFilePath)
		if err != nil {
			return fmt.Errorf("failed to read sql file: %v", err)
		}
		if _, err = conn.Exec(ctx, string(query)); err != nil {
			return fmt.Errorf("failed to init tables: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	return nil
}
