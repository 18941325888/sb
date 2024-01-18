package gapi

import (
	"testing"
	"time"

	db "github.com/18941325888/sb/db/sqlc"
	"github.com/18941325888/sb/util"
	"github.com/18941325888/sb/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}
