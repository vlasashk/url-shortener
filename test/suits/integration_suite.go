//go:build integration

package suits

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener"
	"github.com/vlasashk/url-shortener/pkg/pgconnect"
)

type IntegrationSuite struct {
	suite.Suite
	client         http.Client
	serviceAddress string
	pool           *pgxpool.Pool
	cancel         context.CancelFunc
}

func (s *IntegrationSuite) healthCheck(attempts int) error {
	var err error

	healthURL := s.serviceAddress + "/healthz"

	for attempts > 0 {
		if _, err = s.client.Get(healthURL); err != nil {
			log.Debug().Int("attempts left", attempts).Str("URL", healthURL).Msg("Service is not available for integration rests")
			time.Sleep(time.Second)
			attempts--
			continue
		}
		return nil
	}

	return err
}

func (s *IntegrationSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	cfg, err := config.NewShortener()
	if err != nil {
		s.T().Fatal(err)
	}

	s.serviceAddress = "http://" + net.JoinHostPort(cfg.App.Host, cfg.App.Port)
	s.client = http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Do not follow redirects
		},
		Timeout: time.Second,
	}

	go func() {
		if err = shortener.Run(ctx, cfg); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	if err = s.healthCheck(10); err != nil {
		s.T().Fatal(err)
	}

	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.NameDB)

	s.pool, err = pgconnect.Connect(ctx, url, zerolog.New(os.Stderr))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *IntegrationSuite) TearDownSuite() {
	s.cancel()
	s.pool.Close()
}

func (s *IntegrationSuite) TestCreateAlias() {
	tests := []struct {
		name       string
		expectCode int
		body       io.Reader
		expectResp string
	}{
		{
			name:       "CreateAliasSuccess",
			expectCode: http.StatusCreated,
			body:       bytes.NewBuffer([]byte(`{"original": "https://test.com"}`)),
		},
		{
			name:       "CreateAliasBadJSON",
			expectCode: http.StatusBadRequest,
			body:       bytes.NewBuffer([]byte(`{"BruH"}`)),
			expectResp: `{"error": "bad JSON"}`,
		},
		{
			name:       "CreateAliasSuccess",
			expectCode: http.StatusUnprocessableEntity,
			body:       bytes.NewBuffer([]byte(`{"orig": "https://test.com"}`)),
			expectResp: `{"error": "invalid JSON"}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			r, err := http.NewRequest("POST", fmt.Sprintf("%s/alias", s.serviceAddress), tt.body)
			s.Require().NoError(err)

			resp, err := s.client.Do(r)
			s.Require().NoError(err)
			s.Equal(tt.expectCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)
			_ = resp.Body.Close()
			if tt.expectCode != http.StatusCreated {
				s.JSONEq(tt.expectResp, string(body))
			}
		})
	}
}

func (s *IntegrationSuite) TestGetOriginal() {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO url (alias, original, expires_at, visits) VALUES ($1, $2, $3, $4)`,
		"ABOBUS1234",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		time.Now().AddDate(0, 1, 0),
		0,
	)
	defer func() {
		_, _ = s.pool.Exec(context.Background(), `DELETE FROM url WHERE alias = 'ABOBUS1234'`)
	}()
	s.Require().NoError(err)

	tests := []struct {
		name           string
		expectCode     int
		alias          string
		expectResp     string
		expectLocation string
	}{
		{
			name:           "GetOriginalSuccess",
			expectCode:     http.StatusFound,
			alias:          "ABOBUS1234",
			expectLocation: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		},
		{
			name:       "CreateAliasBadJSON",
			expectCode: http.StatusNotFound,
			alias:      "lolkek",
			expectResp: `{"error": "alias search fail"}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			r, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", s.serviceAddress, tt.alias), nil)
			s.Require().NoError(err)

			resp, err := s.client.Do(r)
			s.Require().NoError(err)
			s.Equal(tt.expectCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)
			_ = resp.Body.Close()
			if tt.expectCode == http.StatusFound {
				location := resp.Header.Get("Location")
				s.Equal(tt.expectLocation, location)
			} else {
				s.JSONEq(tt.expectResp, string(body))
			}
		})
	}
}
