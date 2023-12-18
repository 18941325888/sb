package gapi

import (
	"fmt"

	db "github.com/18941325888/sb/db/sqlc"
	"github.com/18941325888/sb/pb"
	"github.com/18941325888/sb/token"
	"github.com/18941325888/sb/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}

	server := &Server{
		config: config, store: store, tokenMaker: tokenMaker,
	}

	return server, nil
}
