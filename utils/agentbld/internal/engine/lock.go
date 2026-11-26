package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// LockRenderRoot serializes generation and target mutations across processes.
// The kernel releases the advisory lock even if a process exits unexpectedly.
func LockRenderRoot(root string, create bool) (func(), error) {
	if err := ValidDir(root, "render root"); err != nil {
		return nil, err
	}
	if create {
		if err := os.MkdirAll(root, 0755); err != nil {
			return nil, err
		}
	} else if _, err := os.Lstat(root); err != nil {
		return nil, err
	}
	if err := ValidDir(root, "render root"); err != nil {
		return nil, err
	}
	path := filepath.Join(root, ".agentbld.lock")
	fd, err := syscall.Open(path, syscall.O_CREAT|syscall.O_RDWR|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, 0600)
	if err != nil {
		return nil, fmt.Errorf("open render lock %s: %w", path, err)
	}
	file := os.NewFile(uintptr(fd), path)
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	if !info.Mode().IsRegular() {
		file.Close()
		return nil, fmt.Errorf("render lock is not a regular file: %s", path)
	}
	for {
		err = syscall.Flock(fd, syscall.LOCK_EX)
		if err != syscall.EINTR {
			break
		}
	}
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("lock render root: %w", err)
	}
	return func() { _ = syscall.Flock(fd, syscall.LOCK_UN); _ = file.Close() }, nil
}
