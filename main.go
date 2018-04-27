package main

import (
	"os"

	"github.com/ArcherWheeler/bit/bit"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "The best way to learn git (TM)"
	app.Commands = []cli.Command{
		{
			Name:    "feature",
			Aliases: []string{"f"},
			Usage:   "Create a new feature branch",
			Action:  bit.NewBranch,
		},
		{
			Name:    "commit",
			Aliases: []string{"c"},
			Usage:   "Commit all changes to the branch",
			Action:  bit.Commit,
		},
		{
			Name:    "undo",
			Aliases: []string{"u"},
			Usage:   "Undo the last commit",
			Action:  bit.Undo,
		},
		{
			Name:    "switch",
			Aliases: []string{"sw"},
			Usage:   "Move to a different branch",
			Action:  bit.SwitchTo,
		},
		{
			Name:    "publish",
			Aliases: []string{"pb"},
			Usage:   "Push your local changes to the remote repository",
			Action:  bit.Publish,
		},
		{
			Name:    "status",
			Aliases: []string{"st"},
			Usage:   "Your current status",
			Action:  bit.Status,
		},
		{
			Name:    "sync",
			Aliases: []string{"sy"},
			Usage:   "Update and merge with remote changes",
			Action:  bit.Sync,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		bit.Fail(err)
	}
}
