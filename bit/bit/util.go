package bit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SmartStash() {
	if currentChanges() {
		if onMaster() {
			Fail("Do not commit to master")
		}

		git("add", "-A")
		git("commit", "-m", "WIP-BIT-SMART-STASH")
	}
}

func SmartUnstash() {
	lastCommit, _ := gitF("log", "-1", "--pretty=%B")
	if lastCommit == "WIP-BIT-SMART-STASH" {
		Undo()
	}
}

func currentChanges() bool {
	state := git("status", "-s")
	return state != ""
}

func currentBranch() string {
	return git("symbolic-ref", "--short", "HEAD")
}

func onMaster() bool {
	return currentBranch() == "master"
}

func git(args ...string) string {
	out, _ := gitF(args...)
	return out
}

func gitF(args ...string) (string, error) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd := exec.Command("git", args...)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	err := cmd.Run()
	if err != nil {
		if errBuf.String() != "" {
			Fail(errBuf.String())
		}
		Fail(outBuf.String())
	}
	return strings.TrimSpace(outBuf.String()), err
}

func Fail(msg interface{}) {
	fmt.Fprint(os.Stderr, msg)
	os.Exit(1)
}
