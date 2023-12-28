package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/isksss/mcgo/config"
)

type DownloadCmd struct {
}

func (*DownloadCmd) Name() string     { return "download" }
func (*DownloadCmd) Synopsis() string { return "download server and plugin." }
func (*DownloadCmd) Usage() string {
	return ""
}

func (p *DownloadCmd) SetFlags(f *flag.FlagSet) {
}

func (p *DownloadCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	// get config
	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	if err := config.Check(); err != nil {
		panic(err)
	}

	// download server
	fmt.Println("Download server...")
	if err := config.DownloadServer(); err != nil {
		panic(err)
	}

	// download plugins
	fmt.Println("Download plugins...")
	if err := config.DownloadPlugins(); err != nil {
		panic(err)
	}
	return subcommands.ExitSuccess
}
