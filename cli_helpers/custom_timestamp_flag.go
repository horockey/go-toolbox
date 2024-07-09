package cli_helpers

import (
	"errors"
	"flag"
	"fmt"

	"github.com/urfave/cli/v2"
)

var _ cli.Flag = &customTimestampFlag{}

type customTimestampFlag struct {
	cli.TimestampFlag
	Layouts []string
}

func (f *customTimestampFlag) Apply(set *flag.FlagSet) error {
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
