package main

import (
	"fmt"
	"os"

	"github.com/horockey/go-toolbox/cli_helpers"
	"github.com/urfave/cli/v2"
)

var sFlag = &cli_helpers.CustomTimestampFlag{
	Layouts: cli_helpers.DefaultLayouts(),
	TimestampFlag: cli.TimestampFlag{
		Name:    "since",
		Aliases: []string{"s"},
	},
}

func main() {
	app := cli.App{
		Action: func(ctx *cli.Context) error {
			fmt.Println(sFlag.Value)
			fmt.Println(sFlag.TimestampFlag.Value)
			return nil
		},
		Flags: []cli.Flag{sFlag},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
