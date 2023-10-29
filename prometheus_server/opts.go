package prometheus_server

import (
	"fmt"
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
