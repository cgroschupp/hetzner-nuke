package main

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/cgroschupp/hetzner-nuke/pkg/common"

	_ "github.com/cgroschupp/hetzner-nuke/pkg/commands/list"
	_ "github.com/cgroschupp/hetzner-nuke/pkg/commands/project-info"
	_ "github.com/cgroschupp/hetzner-nuke/pkg/commands/run"

	_ "github.com/cgroschupp/hetzner-nuke/resources"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// log panics forces exit
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "remove everything from an hetzner account"
	app.Version = common.AppVersion.Summary
	app.Authors = []*cli.Author{
		{
			Name:  "Christian Groschupp",
			Email: "christian@groschupp.org",
		},
	}

	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logrus.Fatalf("Command %s not found.", command)
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
