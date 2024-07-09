package cli_helpers

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	_ cli.Flag   = &CustomTimestampFlag{}
	_ flag.Value = &CustomTSFlagValue{}
)

type CustomTimestampFlag struct {
	cli.TimestampFlag
	Layouts []string
	Value   *CustomTSFlagValue
}

type CustomTSFlagValue struct {
	ts       time.Time
	location *time.Location
	layouts  []string
	hooks    []func(newVal time.Time)
}

func (v *CustomTSFlagValue) Set(value string) (resErr error) {
	defer func() {
		if resErr != nil {
			return
		}

		ts := v.ts
		for _, hook := range v.hooks {
			hook(ts)
		}
	}()

	var err error
	for _, layout := range v.layouts {
		if v.location != nil {
			v.ts, err = time.ParseInLocation(layout, value, v.location)
			if err != nil {
				resErr = errors.Join(
					resErr,
					fmt.Errorf(
						"parsing time in loc %s: %w",
						v.location.String(),
						err,
					),
				)
				continue
			}

			return nil
		}

		v.ts, err = time.Parse(layout, value)
		if err != nil {
			resErr = errors.Join(
				resErr,
				fmt.Errorf("parsing time: %w", err),
			)
			continue
		}

		return nil
	}

	return resErr
}

func (v *CustomTSFlagValue) String() string {
	return fmt.Sprintf("%#v", v.ts)
}

func (f *CustomTimestampFlag) Apply(set *flag.FlagSet) error {
	if len(f.Layouts) == 0 {
		f.Layouts = DefaultLayouts()
	}
	var resErr error

	if f.Value == nil {
		f.Value = &CustomTSFlagValue{
			hooks: []func(newVal time.Time){
				func(newVal time.Time) {
					f.TimestampFlag.Value = cli.NewTimestamp(newVal)
				},
			},
		}
	}
	f.Value.layouts = f.Layouts
	f.Value.location = f.TimestampFlag.Timezone

	for _, name := range f.Names() {
		if f.Destination != nil {
			set.Var(f.Destination, name, f.Usage)
			continue
		}

		set.Var(f.Value, name, f.Usage)
	}

	return resErr
}

func (f *CustomTimestampFlag) String() string {
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
