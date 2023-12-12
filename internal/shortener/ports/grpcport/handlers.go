package grpcport

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener/models/service"
	"github.com/vlasashk/url-shortener/pkg/transport"
	"google.golang.org/grpc"
	"net"
)

type Handler struct {
	transport.UnimplementedAliasServiceServer
	s   service.Service
	log zerolog.Logger
}

func Run(serv service.Service, log zerolog.Logger, cfg config.AppCfg) error {
	lis, err := net.Listen("tcp", ":"+cfg.Port) // Choose your port
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	transport.RegisterAliasServiceServer(grpcServer, &Handler{s: serv, log: log})
	log.Info().Str("address", cfg.Host+":"+cfg.Port).Msg("starting listening")
	return grpcServer.Serve(lis)
}

func (h *Handler) CreateAlias(_ context.Context, orig *transport.URLRequest) (*transport.AliasResp, error) {
	alias, err := h.s.CrateAlias(orig.Original)
	if err != nil {
		h.log.Error().Err(err).Send()
		return nil, fmt.Errorf("could not create alias: %v", err)
	}
	return &transport.AliasResp{Alias: alias}, nil
}
func (h *Handler) GetOrigURL(_ context.Context, alias *transport.AliasReq) (*transport.OriginalURLResp, error) {
	orig, err := h.s.GetOrigURL(alias.Alias)
	if err != nil {
		h.log.Error().Err(err).Send()
		return nil, fmt.Errorf("alias search fail: %v", err)
	}
	return &transport.OriginalURLResp{OriginalUrl: orig}, nil
}
