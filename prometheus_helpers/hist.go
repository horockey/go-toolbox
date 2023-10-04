package prometheus_helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/horockey/go-toolbox/options"
	"github.com/prometheus/client_golang/prometheus"
)

func NewHistOpts(name string, opts ...options.Option[prometheus.HistogramOpts]) (*prometheus.HistogramOpts, error) {
	histOpts := prometheus.HistogramOpts{
		Name: name,
		Buckets: []float64{
			float64(time.Microsecond * 100),
			float64(time.Millisecond),
			float64(time.Millisecond * 5),
			float64(time.Millisecond * 10),
			float64(time.Millisecond * 20),
			float64(time.Millisecond * 50),
			float64(time.Millisecond * 100),
		},
	}

	if err := options.ApplyOptions(&histOpts, opts...); err != nil {
		return nil, fmt.Errorf("applying opts: %w", err)
	}

	return &histOpts, nil
}

func HistOptsWithNamespace(ns string) options.Option[prometheus.HistogramOpts] {
	return func(target *prometheus.HistogramOpts) error {
		if ns == "" {
			return errors.New("got empty namespace")
		}
		target.Namespace = ns
		return nil
	}
}

func HistOptsWithSubsystem(ss string) options.Option[prometheus.HistogramOpts] {
	return func(target *prometheus.HistogramOpts) error {
		if ss == "" {
			return errors.New("got empty subsystem")
		}
		target.Subsystem = ss
		return nil
	}
}

func HistOptsWithHelp(h string) options.Option[prometheus.HistogramOpts] {
	return func(target *prometheus.HistogramOpts) error {
		if h == "" {
			return errors.New("got empty help")
		}
		target.Help = h
		return nil
	}
}
