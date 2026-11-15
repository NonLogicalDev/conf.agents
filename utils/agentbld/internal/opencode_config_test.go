package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCodePathPrecedence(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("OPENCODE_CONFIG_DIR", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	cases := []struct{ explicit, env, xdg, want string }{
		{"", "", "", filepath.Join(home, ".config", "opencode")},
		{"", "", "/tmp/xdg", "/tmp/xdg/opencode"},
		{"", "/tmp/opencode", "/tmp/xdg", "/tmp/opencode"},
		{"/tmp/explicit", "/tmp/opencode", "/tmp/xdg", "/tmp/explicit"},
	}
	for _, tc := range cases {
		t.Setenv("OPENCODE_CONFIG_DIR", tc.env)
		t.Setenv("XDG_CONFIG_HOME", tc.xdg)
		got, err := openCodePath(tc.explicit)
		if err != nil || got != tc.want {
			t.Fatalf("got %q, %v; want %q", got, err, tc.want)
		}
	}
}

func TestOpenCodeNativePluginPrecedenceAndGuards(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "00-first")
	second := filepath.Join(root, "01-second")
	for _, layer := range []string{first, second} {
		if err := os.MkdirAll(filepath.Join(layer, "opencode/plugins"), 0755); err != nil {
			t.Fatal(err)
		}
	}
	for _, layer := range []string{first, second} {
		if err := os.WriteFile(filepath.Join(layer, "opencode/plugins/example.ts"), []byte(layer), 0644); err != nil {
			t.Fatal(err)
		}
	}
	plugins, err := discoverOpenCodePlugins([]string{first, second})
	if err != nil || plugins["example.ts"] != filepath.Join(second, "opencode/plugins/example.ts") {
		t.Fatalf("plugins %v: %v", plugins, err)
	}
	if err := os.Symlink(filepath.Join(first, "opencode/plugins/example.ts"), filepath.Join(second, "opencode/plugins/unsafe.js")); err != nil {
		t.Fatal(err)
	}
	if _, err := discoverOpenCodePlugins([]string{first, second}); err == nil {
		t.Fatal("accepted a symlinked plugin")
	}
	if err := os.Remove(filepath.Join(second, "opencode/plugins/unsafe.js")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(second, "opencode/plugins/unsafe.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := discoverOpenCodePlugins([]string{first, second}); err == nil {
		t.Fatal("accepted a non-plugin file")
	}
}

func TestOpenCodeRejectsSymlinkedPluginRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "opencode")); err != nil {
		t.Fatal(err)
	}
	if _, err := discoverOpenCodePlugins([]string{root}); err == nil {
		t.Fatal("accepted a symlinked OpenCode root")
	}
}

func TestOpenCodeModeValidation(t *testing.T) {
	for _, tc := range []struct {
		raw   string
		valid bool
	}{
		{"opencode", true}, {`opencode:{"agents":false,"skills":false}`, true},
		{`opencode:{"agents":false,"skills":false,"plugins":false}`, false},
		{`opencode:{"plugins":false,"plugins":true}`, false}, {`opencode:{"plugins":"yes"}`, false}, {`opencode:{"unknown":true}`, false},
	} {
		_, err := parseOpenCodeMode(&tc.raw, syncFlags{})
		if (err == nil) != tc.valid {
			t.Fatalf("mode %s: %v", tc.raw, err)
		}
	}
}
