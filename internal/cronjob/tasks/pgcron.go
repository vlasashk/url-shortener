package tasks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/pkg/pgconnect"
)

type Conn struct {
	conn *pgxpool.Pool
}

func New(ctx context.Context, cfg config.PostgresCfg, logger zerolog.Logger) (*Conn, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.NameDB)

	pool, err := pgconnect.Connect(ctx, url, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return &Conn{pool}, nil
}

func (c *Conn) Close() {
	c.conn.Close()
}

func (c *Conn) CleanUp(ctx context.Context) error {
	_, err := c.conn.Exec(ctx, `DELETE FROM url WHERE expires_at < now()::date`)
	return err
}
