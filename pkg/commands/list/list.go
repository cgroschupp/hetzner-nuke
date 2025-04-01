package list

import (
	"context"
	"sort"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/cgroschupp/hetzner-nuke/pkg/commands/global"
	"github.com/cgroschupp/hetzner-nuke/pkg/common"
	_ "github.com/cgroschupp/hetzner-nuke/resources"
)

func execute(ctx context.Context, c *cli.Command) error {
	ls := registry.GetNames()

	sort.Strings(ls)

	for _, name := range ls {
		reg := registry.GetRegistration(name)

		if reg.AlternativeResource != "" {
			_, _ = color.New(color.Bold).Printf("%-55s\n", name)
			_, _ = color.New(color.Bold, color.FgYellow).Printf("  > %-55s", reg.AlternativeResource)
			_, _ = color.New(color.FgCyan).Printf("alternative resource\n")
		} else {
			_, _ = color.New(color.Bold).Printf("%-55s", name)
			c := color.FgGreen
			_, _ = color.New(c).Printf("%s\n", string(reg.Scope))
		}
	}

	return nil
}

func init() {
	cmd := &cli.Command{
		Name:    "resource-types",
		Aliases: []string{"list-resources"},
		Usage:   "list available resources to nuke",
		Flags:   global.Flags(),
		Before:  global.Before,
		Action:  execute,
	}

	common.RegisterCommand(cmd)
}
