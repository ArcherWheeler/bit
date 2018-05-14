package main

import (
	"os"

	"github.com/ArcherWheeler/bit/bit/bit"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "The best way to learn git (TM)"
	t, err := bit.NewTutor()
	if err != nil {
		bit.Fail(err)
	}

	app.Commands = []cli.Command{
		{
			Name:    "feature",
			Aliases: []string{"f"},
			Usage:   "Create a new feature branch",
			Action:  t.NewBranch,
		},
		{
			Name:    "commit",
			Aliases: []string{"c"},
			Usage:   "Commit all changes to the branch",
			Action:  t.CommitCmd,
		},
		{
			Name:    "undo",
			Aliases: []string{"u"},
			Usage:   "Undo the last commit",
			Action:  func(c *cli.Context) { t.Undo() },
		},
		{
			Name:    "switch",
			Aliases: []string{"sw"},
			Usage:   "Move to a different branch",
			Action:  t.SwitchTo,
		},
		{
			Name:    "publish",
			Aliases: []string{"pb"},
			Usage:   "Push your local changes to the remote repository",
			Action:  t.Publish,
		},
		{
			Name:    "status",
			Aliases: []string{"st"},
			Usage:   "Your current status",
			Action:  t.Status,
		},
		{
			Name:    "sync",
			Aliases: []string{"sy"},
			Usage:   "Update and merge with remote changes",
			Action:  t.Sync,
		},
		{
			Name:    "mode",
			Aliases: []string{"md"},
			Usage:   "Set the show mode to silent, explain or hint",
			Action:  t.SetShowMode,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		bit.Fail(err)
	}
}
