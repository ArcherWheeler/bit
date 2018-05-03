package bit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func SmartStash() {
	stashSave := fmt.Sprintf("WIP-BIT-STASH-%s", currentBranch())
	git("stash", "save", "--include-untracked", stashSave)
}

func SmartUnstash() {
	stashList := git("stash", "list")
	stashes := strings.Split(stashList, "\n")

	branchName := currentBranch()
	regex := fmt.Sprintf(`^stash@{[0-9]+}: On %s: WIP-BIT-STASH-%s$`, branchName, branchName)
	r := regexp.MustCompile(regex)
	var stashNum string
	for _, line := range stashes {
		if r.MatchString(line) {
			stashPart := strings.Split(line, " ")[0]
			stashNum = regexp.MustCompile("[0-9]+").FindString(stashPart)
			break
		}
	}

	if stashNum != "" {
		git("stash", "pop", fmt.Sprintf("stash@{%s}", stashNum))
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
