package bit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBit(t *testing.T) {
	// These tests build and commit to a git repo. In order to stop people
	// from automatically running the tests and creating a bunch of garabe files,
	// this failsafe is used to make sure the tests are only run in the dockerfile.
	failsafe := os.Getenv("TEST_FAILSAFE")
	if failsafe != "off" {
		t.FailNow()
	}

	setUpRepos()

	t.Run("Run bit help", func(t *testing.T) {
		out, err := bit("help")
		require.NotEqual(t, "", out)
		require.Equal(t, "", err)
	})

	t.Run("Bit commit first file", func(t *testing.T) {
		shell("echo foo >> txt")
		out, err := bit("commit", "\"foo\"")

		require.Equal(t, "", out)
		require.Equal(t, "", err)
	})

	t.Run("bit status", func(t *testing.T) {
		out, err := bit("status")
		fmt.Println(out)
		require.Equal(t,
			`On branch master
Your branch is ahead of 'origin/master' by 1 commit.
  (use "git push" to publish your local commits)

nothing to commit, working tree clean`,
			out)
		require.Equal(t, "", err)
	})

	t.Run("bit feature", func(t *testing.T) {
		out, err := bit("feature", "other-branch")
		require.Equal(t, "", out)
		require.Equal(t, "", err)

		branch := shell("git symbolic-ref --short HEAD")
		require.Equal(t, "other-branch", branch)
	})

	t.Run("Bit switch", func(t *testing.T) {
		shell("echo foo >> txt")
		st := shell("git status -s")
		require.Equal(t, "M txt", st)

		bit("switch", "master")
		st = shell("git status -s")
		require.Equal(t, "", st)

		bit("switch", "other-branch")

		st = shell("git status -s")
		require.Equal(t, "M txt", st)
	})
}

func setUpRepos() {
	os.Mkdir("testRepo", os.ModeDir)
	os.Chdir("testRepo")
	shell("git init")

	shell("git config --global user.name 'John Doe'")
	shell("git config --global user.email 'john@doe.com'")

	shell("git config receive.denyCurrentBranch ignore")
	shell("echo foo >> txt")
	shell("git add .")
	shell("git commit -m \"hello world\"")

	os.Chdir("..")
	os.Mkdir("local", os.ModeDir)
	os.Chdir("local")
	shell("git clone ../testRepo")
	os.Chdir("testRepo")
}

func shell(args ...string) string {
	argz := append([]string{"-c"}, args...)
	cmd := exec.Command("sh", argz...)
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	_ = cmd.Run()

	fmt.Println(args)
	fmt.Print(outBuf.String())
	fmt.Print(errBuf.String())

	fmt.Println("")

	return strings.TrimSpace(outBuf.String())
}

func bit(args ...string) (string, string) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd := exec.Command("bit", args...)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	// We want to test the top level api. We care about stdin and stdout not
	// Go's runtime error representation
	_ = cmd.Run()

	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String())
}
