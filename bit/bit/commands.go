package bit

import (
	"fmt"

	"github.com/urfave/cli"
)

func NewBranch(c *cli.Context) {
	branchName := c.Args().First()

	if currentChanges() {
		smartStash()
	}

	git("checkout", "master")
	git("pull")
	git("checkout", "-b", branchName)
	git("push", "--set-upstream", "origin", branchName)
}

func Sync(c *cli.Context) {
	if currentChanges() {
		Fail("Must commit changes first")
	}

	git("fetch")

	if onMaster() {
		return
	}

	stdout := git("merge", "master")
	fmt.Print(stdout)
}

func Commit(c *cli.Context) {
	message := c.Args().First()

	if onMaster() {
		Fail("Do not commit to master")
	}

	git("add", "-A")
	git("commit", "-m", message)
}

func Undo(c *cli.Context) {
	numCommits := git("rev-list", "--count", "master..HEAD")
	if numCommits == "0" {
		Fail("No commits since master to undo")
	}

	git("reset", "--soft", "HEAD^")
	git("reset", "HEAD", ".")
}

func SwitchTo(c *cli.Context) {
	branchName := c.Args().First()

	if currentChanges() {
		smartStash()
	}

	git("checkout", branchName)

	if branchName == "master" {
		git("pull")
	}

	smartUnstash()
}

func Publish(c *cli.Context) {
	if currentChanges() {
		Fail("Must commit changes first")
	}

	if onMaster() {
		Fail("Do not manually publish master")
	}

	git("push")
}

func Status(c *cli.Context) {
	stdout := git("status")
	fmt.Print(stdout)
}
