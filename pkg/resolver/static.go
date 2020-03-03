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

	link := s.Settings.GetLink(database)
	if link.Address == "" {
		return nil, errors.New("database " + database + " is not allowed by pgproker-static")
	}

	return net.Dial("tcp", link.Address)
}

func (s *Static) RewriteParameters(parameters map[string]string) map[string]string {
	database := parameters["database"]

	link := s.Settings.GetLink(database)
	if link.Datname != "" {
		parameters["database"] = link.Datname
	}

	return parameters
}
