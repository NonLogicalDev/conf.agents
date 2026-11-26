package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFragmentHarnessValidation(t *testing.T) {
	for _, header := range []string{"harnesses: alpha", "harnesses: []", "harnesses: [typo]", "harnesses: [alpha, alpha]", "harnesses: [null]", "harnesses: [alpha]\nharnesses: [beta]"} {
		t.Run(header, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "fragment.md")
			if err := os.WriteFile(path, []byte("---\n"+header+"\n---\nbody\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if _, _, err := (Registry{"alpha": {}, "beta": {}}).FragmentForHarness(path, "alpha"); err == nil {
				t.Fatalf("accepted %q", header)
			}
		})
	}
}
