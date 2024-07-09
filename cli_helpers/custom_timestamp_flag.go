package cli_helpers

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var _ cli.Flag = &customTimestampFlag{}

type customTimestampFlag struct {
	cli.TimestampFlag
	Layouts []string
}

func (f *customTimestampFlag) Apply(set *flag.FlagSet) error {
	if len(f.Layouts) == 0 {
		f.Layouts = DefaultLayouts()
	}
	var resErr error
	for _, layout := range f.Layouts {
		f.Layout = layout
		if err := f.TimestampFlag.Apply(set); err != nil {
			resErr = errors.Join(resErr, fmt.Errorf("running super's Apply: %w", err))
			continue
		}
		return nil
	}
	return resErr
}

func (f *customTimestampFlag) String() string {
	return f.TimestampFlag.String() +
		" Available time formats: " +
		strings.Join(f.Layouts, ", ")
}

func DefaultLayouts() []string {
	return []string{
		time.RFC3339,
		time.DateTime,
		time.TimeOnly,
		time.Kitchen,
		"15:04",
		"3:04:05PM",
	}
}
