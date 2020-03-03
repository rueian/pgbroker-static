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
	"time"
)

func NewPGBroker(resolver backend.PGResolver, rewriter backend.PGStartupMessageRewriter, logging bool) *proxy.Server {
	clientStreamCallbacks := proxy.NewStreamCallbackFactories()
	serverStreamCallbacks := proxy.NewStreamCallbackFactories()

	if logging {
		// endless channel and goroutine
		ch := make(chan string, 1000)
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			buf := bufio.NewWriter(os.Stdout)
			for {
				select {
				case <-ticker.C:
					buf.Flush()
				case s := <-ch:
					buf.WriteString(s)
				}
			}
		}()
		clientStreamCallbacks.SetFactory('Q', loggingFactory(ch))
		clientStreamCallbacks.SetFactory('P', loggingFactory(ch))
	}

	server := &proxy.Server{
		PGResolver:                    resolver,
		PGStartupMessageRewriter:      rewriter,
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

func loggingFactory(ch chan<- string) func(ctx *proxy.Ctx) proxy.StreamCallback {
	return func(ctx *proxy.Ctx) proxy.StreamCallback {
		user := ctx.ConnInfo.StartupParameters["user"]
		database := ctx.ConnInfo.StartupParameters["database"]
		ch <- fmt.Sprintf("Query: db=%s user=%s query=", database, user)

		return func(slice proxy.Slice) proxy.Slice {
			var query string
			if !slice.Head {
				if slice.Data[len(slice.Data)-1] == 0 {
					query = string(slice.Data[:len(slice.Data)-1])
				} else {
					query = string(slice.Data)
				}
				query = spaces.ReplaceAllString(query, " ")
			}
			if slice.Last {
				query = query + "\n"
			}
			ch <- query

			return slice
		}
	}
}
