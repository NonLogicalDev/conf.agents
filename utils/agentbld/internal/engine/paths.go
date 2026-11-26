package engine

import (
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
