package list

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/urfave/cli/v3"

	"github.com/cgroschupp/hetzner-nuke/pkg/commands/global"
	"github.com/cgroschupp/hetzner-nuke/pkg/common"
	"github.com/cgroschupp/hetzner-nuke/pkg/hetzner"
	_ "github.com/cgroschupp/hetzner-nuke/resources"
)

func execute(ctx context.Context, c *cli.Command) error {
	client := hetzner.NewClient(hcloud.WithToken(c.String("hcloud-token")))
	token, _, err := client.Token.Current(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Hetzner project with the ID %d and the name %q.\n", token.Token.Project.ID, token.Token.Project.Name)
	return nil
}

func init() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "hcloud-token",
			Usage:    "the hcloud token",
			Sources:  cli.EnvVars("HCLOUD_TOKEN"),
			Required: true,
		},
	}
	cmd := &cli.Command{
		Name:   "project-info",
		Usage:  "list available resources to nuke",
		Flags:  append(flags, global.Flags()...),
		Before: global.Before,
		Action: execute,
	}

	common.RegisterCommand(cmd)
}
