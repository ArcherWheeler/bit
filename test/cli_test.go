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
	os.Mkdir("testRepo", os.ModeDir)
	os.Chdir("testRepo")

	t.Run("Run bit help", func(t *testing.T) {
		out, err := bit("help")
		require.NotEqual(t, "", out)
		require.Equal(t, "", err)
	})

	t.Run("Set up git repo", func(t *testing.T) {
		cmd := exec.Command("git", "init")
		err := cmd.Run()
		require.NoError(t, err)

		shell("git config --global user.name 'John Doe'")
		shell("git config --global user.email 'john@doe.com'")
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
nothing to commit, working tree clean`,
			out)
		require.Equal(t, "", err)
	})
}

func shell(args ...string) {
	argz := append([]string{"-c"}, args...)
	cmd := exec.Command("sh", argz...)
	_ = cmd.Run()
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
