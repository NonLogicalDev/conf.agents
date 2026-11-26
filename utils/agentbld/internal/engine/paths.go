package engine

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultHome resolves the user's home alias so default paths stay consistent
// on systems where /home is a symlink. Explicit paths retain their spelling.
func DefaultHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(home)
	if err != nil {
		if os.IsNotExist(err) {
			return home, nil
		}
		return "", err
	}
	return resolved, nil
}

func Absolute(p string) (string, error) {
	if len(p) > 0 && p[0] == '~' && (len(p) == 1 || p[1] == '/') {
		h, e := DefaultHome()
		if e != nil {
			return "", e
		}
		p = filepath.Join(h, p[1:])
	}
	return filepath.Abs(p)
}

// RejectSymlinkAncestors checks the path and its parents without following links.
func RejectSymlinkAncestors(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for {
		info, err := os.Lstat(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing a symlinked destination: %s", path)
		}
		parent := filepath.Dir(path)
		if parent == path {
			return nil
		}
		path = parent
	}
}
