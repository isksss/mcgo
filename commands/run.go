package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/isksss/mcgo/config"
)

type RunCmd struct {
}

func (*RunCmd) Name() string     { return "run" }
func (*RunCmd) Synopsis() string { return "run server." }
func (*RunCmd) Usage() string {
	return ""
}

func (p *RunCmd) SetFlags(f *flag.FlagSet) {
}

func (p *RunCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// configを取得
	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	if err := config.Check(); err != nil {
		panic(err)
	}

	// サーバーを起動
	if err := config.RunServer(); err != nil {
		panic(err)
	}
	return subcommands.ExitSuccess
}
