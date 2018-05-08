package bit

import (
	"github.com/urfave/cli"
)

func (t *Tutor) NewBranch(c *cli.Context) {
	branchName := c.Args().First()

	if t.currentChanges() {
		t.SmartStash()
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
	if t.currentChanges() {
		Fail("For sync to work, you must commit your changes first")
	}
	branch := t.hide().currentBranch()

	if t.onMaster() {
		return
	}
	t.explain(
		"We want to update our current feature branch with any new changes commited to master.\n\n"+
			"First we switch to the master branch.",
	).git("checkout", "master")

	t.explain(
		"We now update our local copy of master with any changes commited to the remote version",
	).git("pull")

	t.explain(
		"Now that we've updated master, we switch back to our feature branch.",
	).git("checkout", branch)

	stdout :=
		t.explain(
			"Now we merge the new changes to master into our current feature branch.\n\n"+
				"Now our feature branch will have all of the changes done to master, as well as the edits you are currently"+
				"working on.\n\n"+
				"We don't want to merge "+branch+" into master, until we are entirely done with the new feature.",
		).git("merge", "master")

	t.finalOutput(stdout)
}

func (t *Tutor) CommitCmd(c *cli.Context) {
	t.commit(c.Args().First())
}
func (t *Tutor) commit(message string) {
	if t.onMaster() {
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
	numCommits := t.hide().git("rev-list", "--count", "master..HEAD")
	if numCommits == "0" {
		Fail("No commits since this branch was made to undo")
	}

	t.explain("Git doesn't have a clean way to undo the last commit. Don't worry too much about how this works.")
	t.git("reset", "--soft", "HEAD^")
	t.git("reset", "HEAD", ".")
}

func (t *Tutor) SwitchTo(c *cli.Context) {
	branchName := c.Args().First()

	if t.currentChanges() {
		t.SmartStash()
	} else {
		t.explain("Bit can tell you have no changes since your last commit, so you don't need to save anything extra")
	}

	t.explain(
		"We need to tell git the branch we want to switch to",
	).git("checkout", branchName)

	if branchName == "master" {
		t.explain(
			"Bit won't let you modify the master branch, so if you want to look at it, it should match the most up to date" +
				"version from the remote repository",
		).git("pull")
	}

	t.explain(
		"Bit now checks wether any unfished changes were saved on this branch in a \"Work In Progress\" (WIP) commit.",
	).SmartUnstash()
}

func (t *Tutor) Publish(c *cli.Context) {
	if t.currentChanges() {
		Fail("You must commit your changes first")
	}

	if t.onMaster() {
		Fail("Do not manually edit and publish the master branch")
	}

	t.explain(
		"In Git, there is both a local copy of your branch and a remote copy stored in a place like Github.com.",
	).git("push")
}

func (t *Tutor) Status(c *cli.Context) {
	stdout :=
		t.explain(
			"This one is no different than the Git command.",
		).git("status")
	t.finalOutput(stdout)
}

func (t *Tutor) ToggleShowMode(c *cli.Context) {
	err := saveConfig(BitConfig{ShowMode: !t.ShowMode})
	if err != nil {
		Fail(err)
	}
}
