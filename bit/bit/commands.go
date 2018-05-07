package bit

import (
	"github.com/urfave/cli"
)

func (t *Tutor) NewBranch(c *cli.Context) {
	branchName := c.Args().First()

	if currentChanges() {
		SmartStash()
	}

	t.explain(
		`We want to build our new feature by building on current state of the master branch.
Let's switch to the master branch.`,
	).git("checkout", "master")

	t.explain(
		`We want to make sure our current local version of master is the same as the shared remote one`,
	).git("pull")

	t.explain(
		`Let's create the new branch. In Git, making a new branch doesn't change the code, but rather creates a new "branch" that starts identical to the current branch that you're on`,
	).git("checkout", "-b", branchName)

	t.explain(
		`We now want to tell the remote version of the repository about our new (empty) branch.
You could leave this out, but then Git would complain at you later once you push your local changes to the remote repository.
Bit automatically does this now for simplicity and to avoid confusion.`,
	).git("push", "--set-upstream", "origin", branchName)
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

	t.explain(
		"Let's learn about staging!\n"+
			"In Git every change you make can be staged or unstaged. Commits then are only formed from the staged lines."+
			"The unstaged lines aren't lost! They just carry over to the next commit. This can be useful if you want to"+
			"pick and choose what to commit now and later\n"+
			"However, generally you want to commit everything. The flag -A is short for -all and stages every change to be commited",
	).git("add", "-A")
	t.explain(
		"We now commit the changes.\n"+
			"The flag -m stands for -message and lets you pass the message in line to Git. If you don't use -m,"+
			"you still have to have to write a commit message, it's just the default interface is to open a text editor"+
			"to write your changes in. The default editor is usualy vim, which is _very_ confusing and archaic."+
			"We here at Bit find this to be hostle to new users",
	).git("commit", "-m", message)
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
