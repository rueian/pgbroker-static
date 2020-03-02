package proxy

import (
	"bufio"
	"fmt"
	"github.com/rueian/pgbroker/backend"
	"github.com/rueian/pgbroker/proxy"
	"io"
	"log"
	"net"
	"os"
	"regexp"
)

func NewPGBroker(resolver backend.PGResolver, logging bool) *proxy.Server {
	clientStreamCallbacks := proxy.NewStreamCallbackFactories()
	serverStreamCallbacks := proxy.NewStreamCallbackFactories()

	if logging {
		clientStreamCallbacks.SetFactory('Q', loggingFactory)
		clientStreamCallbacks.SetFactory('P', loggingFactory)
	}

	server := &proxy.Server{
		PGResolver:                    resolver,
		ConnInfoStore:                 backend.NewInMemoryConnInfoStore(),
		ServerStreamCallbackFactories: serverStreamCallbacks,
		ClientStreamCallbackFactories: clientStreamCallbacks,
		OnHandleConnError: func(err error, ctx *proxy.Ctx, conn net.Conn) {
			if err == io.EOF {
				return
			}

			client := conn.RemoteAddr().String()
			server := ""
			if ctx.ConnInfo.ServerAddress != nil {
				server = ctx.ConnInfo.ServerAddress.String()
			}
			user := ""
			database := ""
			if ctx.ConnInfo.StartupParameters != nil {
				user = ctx.ConnInfo.StartupParameters["user"]
				database = ctx.ConnInfo.StartupParameters["database"]
			}

			log.Printf("Error: client=%s server=%s user=%s db=%s err=%s", client, server, user, database, err.Error())
		},
	}

	return server
}

var spaces = regexp.MustCompile(`\s+`)

func loggingFactory(ctx *proxy.Ctx) proxy.StreamCallback {
	buf := bufio.NewWriter(os.Stdout)

	user := ctx.ConnInfo.StartupParameters["user"]
	database := ctx.ConnInfo.StartupParameters["database"]
	buf.WriteString(fmt.Sprintf("Query: db=%s user=%s query=", database, user))

	return func(slice proxy.Slice) proxy.Slice {
		if !slice.Head {
			var query string
			if slice.Data[len(slice.Data)-1] == 0 {
				query = string(slice.Data[:len(slice.Data)-1])
			} else {
				query = string(slice.Data)
			}
			query = spaces.ReplaceAllString(query, " ")
			buf.WriteString(query)
		}
		if slice.Last {
			buf.WriteString("\n")
			buf.Flush()
		}

		return slice
	}
}
