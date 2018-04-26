package main

import (
	"bytes"
	"fmt"
	"log"
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
			Action:  rollBack,
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
			Usage:   "Save your current changes and push them to github",
			Action:  save,
		},
		{
			Name:    "status",
			Aliases: []string{"st"},
			Usage:   "Your current status",
			Action:  status,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func newBranch(c *cli.Context) {
	branchName := c.Args().First()

	changes := currentChanges()
	if changes {
		save(c)
	}

	git("checkout", "master")
	git("pull")
	git("checkout", "-b", branchName)
	git("push", "--set-upstream", "origin", branchName)
}

func commit(c *cli.Context) {
	message := c.Args().First()

	if onMaster() {
		log.Fatal("Do not commit to master")
	}

	git("add", "-A")
	git("commit", "-m", message)
}

func rollBack(c *cli.Context) {
	numCommits := git("rev-list", "--count", "master..HEAD")
	if numCommits == "0" {
		log.Fatal("No commits since master to undo")
	}

	git("reset", "--soft", "HEAD^")
	git("reset", "HEAD", ".")
}

func switchTo(c *cli.Context) {
	branchName := c.Args().First()

	if currentChanges() {
		save(c)
	}

	git("checkout", branchName)

	if lastCommitMessage() == "WIP-BIT-SAVE" {
		rollBack(c)
	}
}

func save(c *cli.Context) {
	if onMaster() {
		log.Fatal("Do not modify master")
	}

	if currentChanges() {
		if lastCommitMessage() == "WIP-BIT-SAVE" {
			rollBack(c)
		}

		git("add", "-A")
		git("commit", "-m", "WIP-BIT-SAVE")
	}
	git("push")
}

func status(c *cli.Context) {
	output := git("status")
	fmt.Print(output)
}

func lastCommitMessage() string {
	lastCommit := git("log", "-1", "--pretty=%B", "|", "tee")
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
	buf := new(bytes.Buffer)
	cmd := exec.Command("git", args...)
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
