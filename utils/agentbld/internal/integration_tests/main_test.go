package integration_tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var contractBinary string

func TestMain(m *testing.M) {
	// The lock holder runs only a test helper and may be killed by its parent.
	if os.Getenv("AGENTBLD_TEST_LOCK_ROOT") != "" {
		os.Exit(m.Run())
	}
	dir, err := os.MkdirTemp("", "agentbld-test-bin-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	contractBinary = filepath.Join(dir, "agentbld")
	cmd := exec.Command("go", "build", "-o", contractBinary, "../..")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}
