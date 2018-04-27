package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
			Action:  newBranch,
		},
		{
			Name:    "commit",
			Aliases: []string{"c"},
			Usage:   "Commit all changes to the branch",
			Action:  commit,
		},
		{
			Name:    "undo",
			Aliases: []string{"u"},
			Usage:   "Undo the last commit",
			Action:  undo,
		},
		{
			Name:    "switch",
			Aliases: []string{"s"},
			Usage:   "Move to a different branch",
			Action:  switchTo,
		},
		{
			Name:    "save",
			Aliases: []string{"sv"},
			Usage:   "Save your current changes",
			Action:  save,
		},
		{
			Name:    "status",
			Aliases: []string{"st"},
			Usage:   "Your current status",
			Action:  status,
		},
		{
			Name:    "sync",
			Aliases: []string{"sy"},
			Usage:   "Update and merge with remote changes",
			Action:  sync,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fail(err)
	}
}

func newBranch(c *cli.Context) {
	branchName := c.Args().First()

	changes := currentChanges()
	if changes {
		stashOnHEAD(c)
	}

	git("checkout", "master")
	git("pull")
	git("checkout", "-b", branchName)
	git("push", "--set-upstream", "origin", branchName)
}

func sync(c *cli.Context) {
	if currentChanges() {
		fail("Must commit changes first")
	}

	git("fetch")

	if onMaster() {
		return
	}

	stdout := git("merge", "master")
	fmt.Print(stdout)
}

func commit(c *cli.Context) {
	message := c.Args().First()

	if onMaster() {
		fail("Do not commit to master")
	}

	git("add", "-A")
	git("commit", "-m", message)
}

func undo(c *cli.Context) {
	numCommits := git("rev-list", "--count", "master..HEAD")
	if numCommits == "0" {
		fail("No commits since master to undo")
	}

	git("reset", "--soft", "HEAD^")
	git("reset", "HEAD", ".")
}

func switchTo(c *cli.Context) {
	branchName := c.Args().First()

	if currentChanges() {
		stashOnHEAD(c)
	}

	git("checkout", branchName)

	if lastCommitMessage() == "WIP-BIT-SAVE" {
		undo(c)
	}
}

func save(c *cli.Context) {
	stashOnHEAD(c)
	git("push")
}

func stashOnHEAD(c *cli.Context) {
	if onMaster() {
		fail("Do not modify master")
	}

	if currentChanges() {
		if lastCommitMessage() == "WIP-BIT-SAVE" {
			undo(c)
		}

		git("add", "-A")
		git("commit", "-m", "WIP-BIT-SAVE")
	}
}

func status(c *cli.Context) {
	output := git("status")
	fmt.Print(output)
}

func lastCommitMessage() string {
	lastCommit := git("log", "-1", "--pretty=%B")
	return strings.TrimSpace(lastCommit)
}

func currentChanges() bool {
	state := git("status", "-s")
	return state != ""
}

func onMaster() bool {
	branch := git("rev-parse", "--abbrev-ref", "HEAD")
	return branch == "master"
}

func git(args ...string) string {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd := exec.Command("git", args...)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	err := cmd.Run()
	if err != nil {
		if errBuf.String() != "" {
			fail(errBuf.String())
		}
		fail(outBuf.String())
	}
	return outBuf.String()
}

func fail(msg interface{}) {
	fmt.Print(msg)
	os.Exit(1)
}
