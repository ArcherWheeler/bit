package bit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

func stashOnHEAD(c *cli.Context) {
	if onMaster() {
		Fail("Do not modify master")
	}

	if currentChanges() {
		if lastCommitMessage() == "WIP-BIT-SAVE" {
			Undo(c)
		}

		git("add", "-A")
		git("commit", "-m", "WIP-BIT-SAVE")
	}
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
			Fail(errBuf.String())
		}
		Fail(outBuf.String())
	}
	return outBuf.String()
}

func Fail(msg interface{}) {
	fmt.Print(msg)
	os.Exit(1)
}
