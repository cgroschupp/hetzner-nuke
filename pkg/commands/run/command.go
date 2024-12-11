package run

import (
	"context"
	"fmt"
	"log"
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

	logrus.Tracef("tenant id: %s", c.String("tenant-id"))

	logrus.Trace("preparing to run nuke")

	params := &libnuke.Parameters{
		Force:      c.Bool("force"),
		ForceSleep: c.Int("force-sleep"),
		Quiet:      c.Bool("quiet"),
		NoDryRun:   c.Bool("no-dry-run"),
		Includes:   c.StringSlice("include"),
		Excludes:   c.StringSlice("exclude"),
	}

	parsedConfig, err := config.New(libconfig.Options{
		Path:         c.Path("config"),
		Deprecations: registry.GetDeprecatedResourceTypeMapping(),
	})
	if err != nil {
		return err
	}

	// Initialize the underlying nuke process
	n := libnuke.New(params, nil, parsedConfig.Settings)

	n.SetRunSleep(5 * time.Second)
	n.SetLogger(logrus.WithField("component", "nuke"))

	n.RegisterVersion(fmt.Sprintf("> %s", common.AppVersion.String()))

	// p := &hetzner.Prompt{Parameters: params}
	// n.RegisterPrompt(p.Prompt)

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
	if err := n.RegisterScanner(nuke.Account,
		libscanner.New("account", tenantResourceTypes, &nuke.ListerOpts{Client: hcloud.NewClient(hcloud.WithToken(c.String("hcloud-token")))})); err != nil {
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
