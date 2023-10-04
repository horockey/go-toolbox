package prometheus_helpers

import (
	"time"

	"github.com/horockey/go-toolbox/options"
	"github.com/prometheus/client_golang/prometheus"
)

func NewHistOpts(name string, opts ...options.Option[prometheus.HistogramOpts]) *prometheus.HistogramOpts {
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

	options.ApplyOptions(&histOpts, opts...)

	return &histOpts
}

func HistOptsWithNamespace(ns string) options.Option[prometheus.HistogramOpts] {
	return func(target *prometheus.HistogramOpts) error {
		target.Namespace = ns
		return nil
	}
}

func HistOptsWithSubsystem(ss string) options.Option[prometheus.HistogramOpts] {
	return func(target *prometheus.HistogramOpts) error {
		target.Subsystem = ss
		return nil
	}
}

func HistOptsWithHelp(h string) options.Option[prometheus.HistogramOpts] {
	return func(target *prometheus.HistogramOpts) error {
		target.Help = h
		return nil
	}
}
