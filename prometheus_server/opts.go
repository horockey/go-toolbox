package prometheus_server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/horockey/go-toolbox/options"
)

func WithShutdownTimeout(t time.Duration) options.Option[Server] {
	return func(target *Server) error {
		if t <= 0 {
			return fmt.Errorf("timeout must be positive, got: %d", t)
		}
		target.shutdownTimeout = t
		return nil
	}
}

func WithServer(serv *http.Server) options.Option[Server] {
	return func(target *Server) error {
		if serv == nil {
			return errors.New("givern server is nil")
		}
		target.server = serv
		target.needToStartServer = false
		return nil
	}
}
