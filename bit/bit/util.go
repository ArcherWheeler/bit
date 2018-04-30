package bit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func smartStash() {
	stashSave := fmt.Sprintf("WIP-BIT-STASH-%s", currentBranch())
	git("stash", "save", "--include-untracked", stashSave)
}

func smartUnstash() {
	stashSave := fmt.Sprintf("stash^{/WIP-BIT-STASH-%s}", currentBranch())
	git("stash", "pop", stashSave)
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
