package bit

import (
	"github.com/urfave/cli"
)

func (t *Tutor) NewBranch(c *cli.Context) {
	branchName := c.Args().First()

	if currentChanges() {
		SmartStash()
	}

	t.git("checkout", "master")
	t.git("pull")
	t.git("checkout", "-b", branchName)
	t.git("push", "--set-upstream", "origin", branchName)
}

func (t *Tutor) Sync(c *cli.Context) {
	if currentChanges() {
		Fail("Must commit changes first")
	}

	t.git("fetch")

	if onMaster() {
		return
	}

	stdout := t.git("merge", "master")
	t.finalOutput(stdout)
}

func (t *Tutor) Commit(c *cli.Context) {
	message := c.Args().First()

	if onMaster() {
		Fail("Do not commit to master")
	}

	t.git("add", "-A")
	t.git("commit", "-m", message)
}

func (t *Tutor) Undo() {
	numCommits := t.git("rev-list", "--count", "master..HEAD")
	if numCommits == "0" {
		Fail("No commits since master to undo")
	}

	t.git("reset", "--soft", "HEAD^")
	t.git("reset", "HEAD", ".")
}

func (t *Tutor) SwitchTo(c *cli.Context) {
	branchName := c.Args().First()

	if currentChanges() {
		SmartStash()
	}

	t.git("checkout", branchName)

	if branchName == "master" {
		t.git("pull")
	}

	t.SmartUnstash()
}

func (t *Tutor) Publish(c *cli.Context) {
	if currentChanges() {
		Fail("Must commit changes first")
	}

	if onMaster() {
		Fail("Do not manually publish master")
	}

	t.git("push")
}

func (t *Tutor) Status(c *cli.Context) {
	stdout := t.git("status")
	t.finalOutput(stdout)
}
