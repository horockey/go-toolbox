package prometheus_helpers

import (
	"errors"
	"fmt"

	"github.com/horockey/go-toolbox/options"
	"github.com/prometheus/client_golang/prometheus"
)

func NewSummaryOpts(name string, opts ...options.Option[prometheus.SummaryOpts]) (*prometheus.SummaryOpts, error) {
	sumOpts := prometheus.SummaryOpts{
		Name: name,
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.75: 0.01,
			0.95: 0.001,
			0.99: 0.001,
		},
	}

	if err := options.ApplyOptions(&sumOpts, opts...); err != nil {
		return nil, fmt.Errorf("applying opts: %w", err)
	}

	return &sumOpts, nil
}

func SummaryOptsWithNamespace(ns string) options.Option[prometheus.SummaryOpts] {
	return func(target *prometheus.SummaryOpts) error {
		if ns == "" {
			return errors.New("got empty namespace")
		}
		target.Namespace = ns
		return nil
	}
}

func SummaryOptsWithSubsystem(ss string) options.Option[prometheus.SummaryOpts] {
	return func(target *prometheus.SummaryOpts) error {
		if ss == "" {
			return errors.New("got empty subsystem")
		}
		target.Subsystem = ss
		return nil
	}
}

func SummaryOptsWithHelp(h string) options.Option[prometheus.SummaryOpts] {
	return func(target *prometheus.SummaryOpts) error {
		if h == "" {
			return errors.New("got empty help")
		}
		target.Help = h
		return nil
	}
}
