package run

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	libconfig "github.com/ekristen/libnuke/pkg/config"
	libnuke "github.com/ekristen/libnuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	libscanner "github.com/ekristen/libnuke/pkg/scanner"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/hetzner-nuke/pkg/commands/global"
	"github.com/cgroschupp/hetzner-nuke/pkg/common"
	"github.com/cgroschupp/hetzner-nuke/pkg/config"
	"github.com/cgroschupp/hetzner-nuke/pkg/hetzner"
	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
)

type log2LogrusWriter struct {
	entry *logrus.Entry
}

func (w *log2LogrusWriter) Write(b []byte) (int, error) {
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	w.entry.Trace(string(b))
	return n, nil
}

func execute(c *cli.Context) error { //nolint:funlen,gocyclo
	_, cancel := context.WithCancel(c.Context)
	defer cancel()

	log.SetOutput(&log2LogrusWriter{
		entry: logrus.WithField("source", "standard-logger"),
	})

	logrus.Trace("preparing to run nuke")

	params := &libnuke.Parameters{
		Force:              c.Bool("force"),
		ForceSleep:         c.Int("force-sleep"),
		Quiet:              c.Bool("quiet"),
		NoDryRun:           c.Bool("no-dry-run"),
		Includes:           c.StringSlice("include"),
		Excludes:           c.StringSlice("exclude"),
		WaitOnDependencies: !c.Bool("no-wait-on-dependencies"),
	}

	parsedConfig, err := config.New(libconfig.Options{
		Path:         c.Path("config"),
		Deprecations: registry.GetDeprecatedResourceTypeMapping(),
	})
	if err != nil {
		return err
	}

	hapi := hetzner.NewClient(hcloud.WithToken(c.String("hcloud-token")))
	token, _, err := hapi.Token.Current(c.Context)
	if err != nil {
		return err
	}
	filters, err := parsedConfig.Filters(strconv.Itoa(int(token.Token.Project.ID)))
	if err != nil {
		return fmt.Errorf("parse filter: %w", err)
	}

	logrus.Tracef("project id: %d", token.Token.Project.ID)

	// Initialize the underlying nuke process
	n := libnuke.New(params, filters, parsedConfig.Settings)

	n.SetRunSleep(5 * time.Second)
	n.SetLogger(logrus.WithField("component", "nuke"))

	n.RegisterVersion(fmt.Sprintf("> %s", common.AppVersion.String()))

	p := &hetzner.Prompt{
		Parameters: params,
		Project:    &token.Token.Project,
	}
	n.RegisterPrompt(p.Prompt)

	tenantResourceTypes := types.ResolveResourceTypes(
		registry.GetNamesForScope(nuke.Account),
		[]types.Collection{
			n.Parameters.Includes,
			parsedConfig.ResourceTypes.GetIncludes(),
		},
		[]types.Collection{
			n.Parameters.Excludes,
			parsedConfig.ResourceTypes.Excludes,
		},
		nil,
		nil,
	)

	logrus.Debug("registering scanner for account")
	scanner, err := libscanner.New(&libscanner.Config{
		ResourceTypes: tenantResourceTypes,
		Owner:         "account",
		Opts:          &nuke.ListerOpts{Client: hapi},
	})
	if err != nil {
		return err
	}
	if err := n.RegisterScanner(nuke.Account, scanner); err != nil {
		return err
	}

	logrus.Debug("running ...")

	return n.Run(c.Context)
}

func init() {
	flags := []cli.Flag{
		&cli.PathFlag{
			Name:  "config",
			Usage: "path to config file",
			Value: "config.yaml",
		},
		&cli.StringSliceFlag{
			Name:  "include",
			Usage: "only include this specific resource",
		},
		&cli.StringSliceFlag{
			Name:  "exclude",
			Usage: "exclude this specific resource (this overrides everything)",
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
			Usage:   "hide filtered messages",
		},
		&cli.BoolFlag{
			Name:  "no-dry-run",
			Usage: "actually run the removal of the resources after discovery",
		},
		&cli.BoolFlag{
			Name:    "no-prompt",
			Usage:   "disable prompting for verification to run",
			Aliases: []string{"force"},
		},
		&cli.BoolFlag{
			Name:  "no-wait-on-dependencies",
			Usage: "disable waiting for dependencies",
		},
		&cli.IntFlag{
			Name:    "prompt-delay",
			Usage:   "seconds to delay after prompt before running (minimum: 3 seconds)",
			Value:   10,
			Aliases: []string{"force-sleep"},
		},
		&cli.StringFlag{
			Name:     "hcloud-token",
			Usage:    "the hcloud token",
			EnvVars:  []string{"HCLOUD_TOKEN"},
			Required: true,
		},
	}

	cmd := &cli.Command{
		Name:    "run",
		Aliases: []string{"nuke"},
		Usage:   "run nuke against an hetzner account to remove all configured resources",
		Flags:   append(flags, global.Flags()...),
		Before:  global.Before,
		Action:  execute,
	}

	common.RegisterCommand(cmd)
}
