package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/isksss/mcgo/config"
)

type InitCmd struct {
	project string
}

func (*InitCmd) Name() string     { return "init" }
func (*InitCmd) Synopsis() string { return "Create default mcgo.toml." }
func (*InitCmd) Usage() string {
	return ""
}

func (p *InitCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.project, "project", config.ProjectPaper, "project name")
}

func (p *InitCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if p.project == config.ProjectPaper {
		config.DefaultConfig(config.ProjectPaper)
	} else if p.project == config.ProjectVelocity {
		config.DefaultConfig(config.ProjectVelocity)
	} else {
		panic("project name is invalid")
	}
	return subcommands.ExitSuccess
}
