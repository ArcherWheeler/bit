package bit

import (
	"os"

	"github.com/dollarshaveclub/line"
)

func (t *Tutor) SmartStash() {
	if t.currentChanges() {
		if t.onMaster() {
			Fail("Do not commit to master")
		}

		t.explain(
			p(`You currently have changes since your last commit. We need to save them before switching branches.`),
			p(`Annoyingly Git treats uncommited changes seperatly from everything else. If you checkout a branch
			 with uncommited changes the changes move with you. This can cause unexpected conflicts between you changes.
			 and the new branch. While this can be useful, we at Bit think this shouldn't be the default behavior.`),
			p(`To get around this Bit saves any current changes into a "Work In Progress" (WIP) commit. Then when you
			 come back to the branch it checks and undoes that commit so you're back where you left off.`),
		).commit("WIP-BIT-SMART-STASH")
	}
}

func (t *Tutor) SmartUnstash() {
	t.explain(
		p(`Bit now checks wether any unfished changes were saved on this branch in a "Work In Progress" (WIP) commit.`),
	)
	lastCommit, _ := t.hide().gitF("log", "-1", "--pretty=%B")
	if lastCommit == "WIP-BIT-SMART-STASH" {
		t.explain(
			p(`Bit has noticed the last commit was a "Work In Progess" (WIP) commit created by switching away from a branch
			with which had unfinished changes.`),
			p(`We now run "Bit undo" to undo this temporary commit`),
		).Undo()
	} else {
		t.explain(
			p("The last commit message was " + lastCommit + " which was not auto generated by bit. No need to make changes."),
		)
	}
}

func (t *Tutor) currentChanges() bool {
	state := t.hide().git("status", "-s")
	return state != ""
}

func (t *Tutor) currentBranch() string {
	return t.hide().git("symbolic-ref", "--short", "HEAD")
}

func (t *Tutor) onMaster() bool {
	return t.currentBranch() == "master"
}

func Fail(msg interface{}) {
	output := line.New(os.Stderr, "", "", line.RedColor)
	output.Println(msg)
	os.Exit(1)
}
