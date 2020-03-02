package resolver

import (
	"context"
	"errors"
	"github.com/rueian/pgbroker-static/pkg/config"
	"net"
)

type Static struct {
	Settings *config.Settings
}

func (s *Static) GetPGConn(ctx context.Context, clientAddr net.Addr, parameters map[string]string) (net.Conn, error) {
	database := parameters["database"]

	address := s.Settings.GetAddress(database)
	if address == "" {
		return nil, errors.New("database " + database + " is not allowed by pgproker-static")
	}

	return net.Dial("tcp", address)
}
