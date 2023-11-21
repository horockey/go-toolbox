package prometheus_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/horockey/go-toolbox/options"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	registry *prometheus.Registry
	server   *http.Server

	shutdownTimeout time.Duration
}

func New(addr string, opts ...options.Option[Server]) (*Server, error) {
	s := Server{
		registry:        prometheus.NewRegistry(),
		server:          &http.Server{Addr: addr},
		shutdownTimeout: time.Second,
	}

	if err := options.ApplyOptions(&s, opts...); err != nil {
		return nil, fmt.Errorf("applying opts: %w", err)
	}

	return &s, nil
}

func (s *Server) Register(cols ...prometheus.Collector) error {
	for idx, col := range cols {
		if err := s.registry.Register(col); err != nil {
			return fmt.Errorf("registering collector %d: %w", idx, err)
		}
	}
	return nil
}

func (s *Server) Start(ctx context.Context) error {
	http.Handle("/metrics", promhttp.HandlerFor(
		s.registry,
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		},
	))
	errs := make(chan error)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		if err != nil {
			return fmt.Errorf("running server: %w", err)
		}
	case <-ctx.Done():
		sdCtx, cancel := context.WithTimeout(context.TODO(), s.shutdownTimeout)
		defer cancel()

		if err := s.server.Shutdown(sdCtx); err != nil {
			return fmt.Errorf("shutting down server: %w", err)
		}

		return nil
	}

	return nil
}
