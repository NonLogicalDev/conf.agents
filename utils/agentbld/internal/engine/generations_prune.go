package engine

import (
	"fmt"
	"os"
	"path/filepath"
)

func (registry Registry) Prune(root string) error {
	gens, e := registry.DiscoverGenerations(root)
	if e != nil {
		return e
	}
	split := len(gens) - 5
	if split < 0 {
		split = 0
	}
	remove, kept := gens[:split], gens[split:]
	if len(remove) > 0 {
		if !AliasMatches(root, gens[len(gens)-1].Path) {
			return fmt.Errorf("latest generation alias must target the newest generation before pruning")
		}
		if err := registry.pruneTargets(root, remove); err != nil {
			return err
		}

		for _, g := range remove {
			if e = removeGeneration(g.Path); e != nil {
				return e
			}
		}
	}
	latest := "none"
	if len(kept) > 0 {
		latest = filepath.Base(kept[len(kept)-1].Path)
	}
	fmt.Printf("pruned agentbld: kept=%d removed=%d latest=%s destination=%s\n", len(kept), len(remove), latest, root)
	return nil
}

// Restore permissions on surviving directories if pruning fails partway.
func removeGeneration(root string) (err error) {
	modes := map[string]os.FileMode{}
	defer func() {
		if err != nil {
			for path, mode := range modes {
				_ = os.Chmod(path, mode)
			}
		}
	}()
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	modes[root] = info.Mode().Perm()
	if err = os.Chmod(root, modes[root]|0200); err != nil {
		return err
	}
	if err = walk(root, func(path string, d os.DirEntry) error {
		if !d.IsDir() {
			return nil
		}
		i, e := d.Info()
		if e != nil {
			return e
		}
		modes[path] = i.Mode().Perm()
		return os.Chmod(path, modes[path]|0200)
	}); err != nil {
		return err
	}
	return os.RemoveAll(root)
}
