package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/isksss/mcgo/commands"
	"github.com/isksss/mcgo/config"
)

const (
	javaCmd = "java"
)

// initialize
func init() {
	_, err := config.GetCmdPath(javaCmd)
	if err != nil {
		panic(err)
	}

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&commands.DownloadCmd{}, "")
	subcommands.Register(&commands.InitCmd{}, "")
	subcommands.Register(&commands.RunCmd{}, "")

	flag.Parse()
}

func main() {
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
