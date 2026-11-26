package integration_tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

// The helper runs in a separate process so the test exercises the real OS lock.
func TestRenderRootLockHolder(t *testing.T) {
	root := os.Getenv("AGENTBLD_TEST_LOCK_ROOT")
	if root == "" {
		return
	}
	release, err := engine.LockRenderRoot(root, false)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if err := os.WriteFile(os.Getenv("AGENTBLD_TEST_LOCK_READY"), []byte("ready"), 0600); err != nil {
		t.Fatal(err)
	}
	releaseFile := os.Getenv("AGENTBLD_TEST_LOCK_RELEASE")
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(releaseFile); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("lock holder timed out")
}

type testLockHolder struct {
	cmd         *exec.Cmd
	done        chan error
	releaseFile string
}

func startTestLockHolder(t *testing.T, f *contractFixture) *testLockHolder {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ready := filepath.Join(f.root, "holder-ready")
	release := filepath.Join(f.root, "holder-release")
	cmd := exec.Command(exe, "-test.run=^TestRenderRootLockHolder$")
	cmd.Env = append(os.Environ(), "AGENTBLD_TEST_LOCK_ROOT="+f.render, "AGENTBLD_TEST_LOCK_READY="+ready, "AGENTBLD_TEST_LOCK_RELEASE="+release)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	h := &testLockHolder{cmd: cmd, done: make(chan error, 1), releaseFile: release}
	go func() { h.done <- cmd.Wait() }()
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			<-h.done
		}
	})
	until := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			return h
		}
		select {
		case err := <-h.done:
			t.Fatalf("lock holder exited before ready: %v", err)
		default:
		}
		if time.Now().After(until) {
			t.Fatal("lock holder did not acquire lock")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
func startLockedCommand(t *testing.T, args ...string) (*exec.Cmd, <-chan error, *strings.Builder) {
	t.Helper()
	cmd := exec.Command(contractBinary, args...)
	var output strings.Builder
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			<-done
		}
	})
	select {
	case err := <-done:
		t.Fatalf("command unexpectedly finished while root locked: %v\n%s", err, &output)
	case <-time.After(200 * time.Millisecond):
	}
	return cmd, done, &output
}
func finishLockedCommand(t *testing.T, done <-chan error, out *strings.Builder) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("command after lock release: %v\n%s", err, out)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("command did not finish after lock release")
	}
}

func TestPruneWaitsForActivationUnderRenderRootLock(t *testing.T) {
	f := modeFixture(t)
	render := func() {
		lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render)
	}
	render()
	first := f.latest()
	number := lifecycleOrdinal(t, first)
	lifecycleRun(t, f, 0, "sync", number, "--render-dir", f.render, "--codex-dir", f.codex)
	agents := filepath.Join(f.codex, "AGENTS.md")
	target, err := os.Readlink(agents)
	if err != nil {
		t.Fatal(err)
	}
	// Temporarily remove every active link to model the gap before activation.
	for _, pattern := range []string{"AGENTS.md", "skills/*", "resources/*", "rules/*"} {
		paths, err := filepath.Glob(filepath.Join(f.codex, pattern))
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range paths {
			if info, err := os.Lstat(p); err == nil && info.Mode()&os.ModeSymlink != 0 {
				if err = os.Remove(p); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	for i := 0; i < 5; i++ {
		f.put(f.layer, "AGENTS.md.d/10-base.md", strings.Repeat("x", i+1))
		render()
	}
	holder := startTestLockHolder(t, f)
	_, done, out := startLockedCommand(t, "prune", "--render-dir", f.render, "--codex-dir", f.codex, "--sections-dir", filepath.Join(f.config, "sections"))
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("prune removed generation while locked: %v", err)
	}
	if err := os.Symlink(target, agents); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(holder.releaseFile, []byte("release"), 0600); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-holder.done:
		if err != nil {
			t.Fatalf("holder: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("holder did not release lock")
	}
	// Prune may refuse to remove an active generation rather than succeed.
	select {
	case err := <-done:
		if err != nil {
			e, ok := err.(*exec.ExitError)
			if !ok || e.ExitCode() != 2 || !strings.Contains(strings.ToLower(out.String()), "active") {
				t.Fatalf("unexpected prune error: %v\n%s", err, out)
			}
		}
	case <-time.After(10 * time.Second):
		t.Fatal("prune remained blocked")
	}
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("active generation pruned: %v", err)
	}
	if _, err := os.ReadFile(agents); err != nil {
		t.Fatalf("active instructions broken: %v", err)
	}
}
func TestRenderRootLockReleasedOnHolderDeath(t *testing.T) {
	f := modeFixture(t)
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render)
	number := lifecycleOrdinal(t, f.latest())
	holder := startTestLockHolder(t, f)
	_, done, out := startLockedCommand(t, "sync", number, "--render-dir", f.render, "--codex-dir", f.codex)
	if err := holder.cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-holder.done:
	case <-time.After(10 * time.Second):
		t.Fatal("killed lock holder did not exit")
	}
	finishLockedCommand(t, done, out)
	f.symlink(filepath.Join(f.codex, "AGENTS.md"))
}
