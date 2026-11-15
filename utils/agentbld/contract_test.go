package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var contractBinary string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "agentbld-test-bin-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	contractBinary = filepath.Join(dir, "agentbld")
	cmd := exec.Command("go", "build", "-o", contractBinary, ".")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

type contractFixture struct {
	t                                        *testing.T
	root, config, render, codex, home, layer string
}

func newContractFixture(t *testing.T) *contractFixture {
	t.Helper()
	root := t.TempDir()
	f := &contractFixture{t: t, root: root, config: filepath.Join(root, "config"), render: filepath.Join(root, "render"), codex: filepath.Join(root, "codex"), home: filepath.Join(root, "home")}
	f.layer = filepath.Join(f.config, "sections", "00-public")
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Base {{%_resources_%}} {{%_skills_%}}\n")
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example skill\n---\nSee {{%_resources_%}}/guide/first.md\n")
	f.put(f.layer, "resources/guide/first.md", "[next](second.md)\n")
	f.put(f.layer, "resources/guide/second.md", "linked resource\n")
	f.put(f.layer, "rules/default.rules", "prefix_rule(pattern=[\"echo\"], decision=\"allow\")\n")
	t.Cleanup(func() {
		_ = filepath.WalkDir(f.root, func(path string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				_ = os.Chmod(path, 0700)
			}
			return nil
		})
	})
	return f
}
func (f *contractFixture) put(root, name, value string) string {
	f.t.Helper()
	p := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		f.t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(value), 0644); err != nil {
		f.t.Fatal(err)
	}
	return p
}
func (f *contractFixture) run(op string, want int, extra ...string) string {
	f.t.Helper()
	args := []string{op, "--config-dir", f.config, "--render-dir", f.render, "--codex-dir", f.codex, "--home", f.home}
	args = append(args, extra...)
	cmd := exec.Command(contractBinary, args...)
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			code = e.ExitCode()
		} else {
			f.t.Fatal(err)
		}
	}
	if code != want {
		f.t.Fatalf("%v: exit %d, want %d\n%s", args, code, want, out)
	}
	return string(out)
}
func (f *contractFixture) latest() string {
	f.t.Helper()
	p, err := filepath.EvalSymlinks(filepath.Join(f.render, "latest"))
	if err != nil {
		f.t.Fatal(err)
	}
	return p
}
func (f *contractFixture) read(p string) string {
	f.t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		f.t.Fatal(err)
	}
	return string(b)
}
func (f *contractFixture) resolved(p string) string {
	f.t.Helper()
	s, err := filepath.EvalSymlinks(p)
	if err != nil {
		f.t.Fatal(err)
	}
	return s
}
func (f *contractFixture) absent(p string) {
	f.t.Helper()
	if _, err := os.Lstat(p); !os.IsNotExist(err) {
		f.t.Fatalf("expected absent: %s (err %v)", p, err)
	}
}
func (f *contractFixture) symlink(p string) {
	f.t.Helper()
	i, err := os.Lstat(p)
	if err != nil || i.Mode()&os.ModeSymlink == 0 {
		f.t.Fatalf("expected symlink: %s (%v)", p, err)
	}
}
func (f *contractFixture) json(p string) map[string]any {
	f.t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(f.read(p)), &m); err != nil {
		f.t.Fatal(err)
	}
	return m
}

func TestContractRenderCheckAndImmutableGenerations(t *testing.T) {
	f := newContractFixture(t)
	f.run("check", 1)
	f.run("render", 0)
	first := f.latest()
	agents := f.read(filepath.Join(first, "AGENTS.md"))
	for _, p := range []string{filepath.Join(f.codex, "resources"), filepath.Join(f.codex, "skills")} {
		if !strings.Contains(agents, p) {
			t.Fatalf("missing %s in %s", p, agents)
		}
	}
	if got := f.read(filepath.Join(first, "skills/Alpha/resources/guide/second.md")); got != "linked resource\n" {
		t.Fatal(got)
	}
	if !strings.Contains(f.read(filepath.Join(first, "skills/Alpha/SKILL.md")), "resources/guide/first.md") {
		t.Fatal("resource placeholder not expanded")
	}
	err := filepath.WalkDir(first, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		i, e := d.Info()
		if e != nil {
			return e
		}
		if i.Mode().Perm()&0222 != 0 {
			return fmt.Errorf("writable generation path %s", p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	f.run("check", 0)
	f.run("render", 0)
	if first != f.latest() {
		t.Fatal("unchanged generation replaced")
	}
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Changed\n")
	f.run("check", 1)
	f.run("render", 0)
	if first == f.latest() {
		t.Fatal("generation not changed")
	}
	if !strings.Contains(f.read(filepath.Join(first, "AGENTS.md")), "Base ") {
		t.Fatal("prior generation modified")
	}
	f.run("check", 0)
}
func TestContractPrecedenceAndGuidance(t *testing.T) {
	f := newContractFixture(t)
	later := filepath.Join(f.config, "sections/99-local")
	f.put(later, "AGENTS.md.d/10-base.md", "Winning fragment\n")
	f.put(later, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Replacement\n---\nWinner\n")
	f.put(later, "resources/guide/first.md", "winning resource\n")
	f.put(f.layer, "skills/Agent_Guidance/SKILL.md", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n{{%_agents_%}}\n")
	f.run("render", 0)
	gen := f.latest()
	agents := f.read(filepath.Join(gen, "AGENTS.md"))
	if !strings.Contains(agents, "Winning fragment") || strings.Contains(agents, "Base ") {
		t.Fatal(agents)
	}
	if !strings.Contains(f.read(filepath.Join(gen, "skills/Alpha/SKILL.md")), "Winner") {
		t.Fatal("skill precedence")
	}
	if f.read(filepath.Join(gen, "resources/guide/first.md")) != "winning resource\n" {
		t.Fatal("resource precedence")
	}
	if !strings.Contains(f.read(filepath.Join(gen, "skills/Agent_Guidance/SKILL.md")), "Winning fragment") {
		t.Fatal("guidance not expanded")
	}
}
func TestContractSyncAndInstructionsOnly(t *testing.T) {
	f := newContractFixture(t)
	f.run("sync", 0)
	first := f.latest()
	for dest, src := range map[string]string{"AGENTS.md": "AGENTS.md", "skills/Alpha": "skills/Alpha", "rules/default.rules": "rules/default.rules"} {
		if f.resolved(filepath.Join(f.codex, dest)) != filepath.Join(first, src) {
			t.Fatal(dest)
		}
	}
	local := filepath.Join(f.codex, "skills/Local")
	if err := os.Mkdir(local, 0755); err != nil {
		t.Fatal(err)
	}
	f.put(f.layer, "AGENTS.md.d/10-base.md", "New instruction\n")
	f.run("sync", 0)
	if first == f.latest() || f.resolved(filepath.Join(f.codex, "skills/Alpha")) != filepath.Join(f.latest(), "skills/Alpha") {
		t.Fatal("links not updated")
	}
	f.run("sync", 0, "--instructions-only")
	f.absent(filepath.Join(f.codex, "skills/Alpha"))
	if _, err := os.Stat(local); err != nil {
		t.Fatal(err)
	}
	f.symlink(filepath.Join(f.codex, "AGENTS.md"))
}
func TestContractSyncCollisionIsPreflighted(t *testing.T) {
	f := newContractFixture(t)
	p := filepath.Join(f.codex, "skills/Alpha")
	if err := os.MkdirAll(p, 0755); err != nil {
		t.Fatal(err)
	}
	f.run("sync", 2)
	f.absent(filepath.Join(f.codex, "AGENTS.md"))
	f.absent(filepath.Join(f.codex, "rules"))
	if _, err := os.Stat(p); err != nil {
		t.Fatal(err)
	}
}
func TestContractRulesOnly(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.codex, "AGENTS.md", "unmanaged")
	f.run("sync", 0, "--rules-only")
	if f.read(filepath.Join(f.codex, "AGENTS.md")) != "unmanaged" {
		t.Fatal("instructions changed")
	}
	f.symlink(filepath.Join(f.codex, "rules/default.rules"))
	f.absent(filepath.Join(f.codex, "skills"))
}
func TestContractPruneActiveGenerationAndRetention(t *testing.T) {
	f := newContractFixture(t)
	f.run("sync", 0)
	first := f.latest()
	for n := 1; n <= 5; n++ {
		f.put(f.layer, "AGENTS.md.d/10-base.md", fmt.Sprintf("Generation %d\n", n))
		f.run("render", 0)
	}
	f.run("prune", 2)
	if _, err := os.Stat(first); err != nil {
		t.Fatal(err)
	}
	f.run("sync", 0)
	f.run("prune", 0)
	f.absent(first)
	entries, err := os.ReadDir(f.render)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			count++
		}
	}
	if count != 5 {
		t.Fatalf("retained %d", count)
	}
	f.run("check", 0)
}
func TestContractResetBackups(t *testing.T) {
	f := newContractFixture(t)
	f.run("sync", 0)
	f.put(f.codex, "skills/Local/keep", "data")
	f.run("reset", 0)
	f.absent(filepath.Join(f.codex, "AGENTS.md"))
	f.absent(filepath.Join(f.codex, "skills/Alpha"))
	if f.read(filepath.Join(f.codex, "skills/Local/keep")) != "data" {
		t.Fatal("local file lost")
	}
	a, _ := filepath.Glob(filepath.Join(f.codex, "AGENTS.*.bak.md"))
	s, _ := filepath.Glob(filepath.Join(f.codex, "skills.*.bak"))
	if len(a) != 1 || len(s) != 1 {
		t.Fatalf("backups: %v %v", a, s)
	}
	f.symlink(a[0])
	f.symlink(filepath.Join(s[0], "Alpha"))
	if f.read(filepath.Join(s[0], "Local/keep")) != "data" {
		t.Fatal("backup missing data")
	}
}
func TestContractRejectUnsafeSources(t *testing.T) {
	f := newContractFixture(t)
	p := filepath.Join(f.layer, "resources/guide/evil")
	if err := os.Symlink("/etc/passwd", p); err != nil {
		t.Fatal(err)
	}
	f.run("render", 2)
	f.absent(f.render)
	if err := os.Remove(p); err != nil {
		t.Fatal(err)
	}
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Test\n---\n{{%_resources_%}}/missing.md\n")
	f.run("render", 2)
	f.absent(f.render)
}
func TestContractRejectIncompatibleModes(t *testing.T) {
	f := newContractFixture(t)
	f.run("render", 2, "--instructions-only")
	f.run("sync", 2, "--rules-only", "--instructions-only")
	f.absent(f.render)
}

func TestContractMarketplaceBuildPreservesCatalogAndRejectsEdits(t *testing.T) {
	f := newContractFixture(t)
	catalogPath := f.put(f.home, ".agents/plugins/marketplace.json", `{"name":"custom","plugins":[{"name":"other","source":{"source":"npm","package":"example"},"category":"Other"}]}`)
	f.run("sync", 0, "--skill-marketplace")
	plugin := filepath.Join(f.home, "plugins/skillset-uoleg")
	manifest := filepath.Join(plugin, ".codex-plugin/plugin.json")
	if f.json(manifest)["version"] != "0.1.0" {
		t.Fatal("initial version")
	}
	if !strings.Contains(f.read(filepath.Join(plugin, "skills/Alpha/SKILL.md")), "Example skill") {
		t.Fatal("skill missing")
	}
	f.absent(filepath.Join(f.codex, "skills/Alpha"))
	catalog := f.json(catalogPath)
	if catalog["name"] != "custom" {
		t.Fatal(catalog)
	}
	plugins := catalog["plugins"].([]any)
	if len(plugins) != 2 || plugins[0].(map[string]any)["name"] != "other" {
		t.Fatal(plugins)
	}
	var entry map[string]any
	for _, p := range plugins {
		m := p.(map[string]any)
		if m["name"] == "skillset-uoleg" {
			entry = m
		}
	}
	if entry == nil {
		t.Fatal("missing entry")
	}
	entry["policy"].(map[string]any)["installation"] = "INSTALLED_BY_DEFAULT"
	b, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(catalogPath, b, 0644); err != nil {
		t.Fatal(err)
	}
	f.run("sync", 0, "--skill-marketplace")
	if f.json(manifest)["version"] != "0.1.1" {
		t.Fatal("version did not increment")
	}
	catalog = f.json(catalogPath)
	plugins = catalog["plugins"].([]any)
	if len(plugins) != 2 || plugins[0].(map[string]any)["name"] != "other" {
		t.Fatal(plugins)
	}
	for _, p := range plugins {
		m := p.(map[string]any)
		if m["name"] == "skillset-uoleg" && m["policy"].(map[string]any)["installation"] != "INSTALLED_BY_DEFAULT" {
			t.Fatal("policy changed")
		}
	}
	before := f.read(catalogPath)
	f.put(plugin, "skills/Alpha/SKILL.md", "local edits")
	f.run("sync", 2, "--skill-marketplace")
	if f.read(filepath.Join(plugin, "skills/Alpha/SKILL.md")) != "local edits" || f.json(manifest)["version"] != "0.1.1" || f.read(catalogPath) != before {
		t.Fatal("edited plugin or catalog overwritten")
	}
}
func TestContractMarketplaceRejectsSymlinkDestination(t *testing.T) {
	f := newContractFixture(t)
	outside := filepath.Join(f.root, "outside")
	if err := os.MkdirAll(outside, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(f.home, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(f.home, "plugins")); err != nil {
		t.Fatal(err)
	}
	f.run("sync", 2, "--skill-marketplace")
	items, err := os.ReadDir(outside)
	if err != nil || len(items) != 0 {
		t.Fatalf("destination modified %v %v", items, err)
	}
	f.absent(filepath.Join(f.home, ".agents/plugins/marketplace.json"))
}
func TestContractMarketplaceRejectsInvalidYAML(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: [unterminated\n---\ntext\n")
	f.run("sync", 2, "--skill-marketplace")
	f.absent(filepath.Join(f.home, "plugins/skillset-uoleg"))
	f.absent(filepath.Join(f.home, ".agents/plugins/marketplace.json"))
	f.absent(filepath.Join(f.codex, "AGENTS.md"))
}

func TestContractCLIRejectsInvalidArguments(t *testing.T) {
	for _, args := range [][]string{
		{"invalid", "--config-dir", "/tmp/config", "--render-dir", "/tmp/render"},
		{"render", "--unknown", "--config-dir", "/tmp/config", "--render-dir", "/tmp/render"},
		{"render"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			cmd := exec.Command(contractBinary, args...)
			out, err := cmd.CombinedOutput()
			exit, ok := err.(*exec.ExitError)
			if !ok || exit.ExitCode() != 2 {
				t.Fatalf("%q: got %v, want exit 2\n%s", args, err, out)
			}
		})
	}
}

func TestContractResetBacksUpReadOnlyNestedDirectories(t *testing.T) {
	f := newContractFixture(t)
	original := f.put(f.codex, "skills/Local/locked/deeper/content.txt", "preserved\n")
	locked := filepath.Join(f.codex, "skills/Local/locked")
	deeper := filepath.Dir(original)
	if err := os.Chmod(original, 0755); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{deeper, locked} {
		if err := os.Chmod(p, 0555); err != nil {
			t.Fatal(err)
		}
	}
	args := []string{"-c", `umask 077; exec "$@"`, "sh", contractBinary, "reset", "--config-dir", f.config, "--render-dir", f.render, "--codex-dir", f.codex}
	out, err := exec.Command("sh", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("reset: %v\n%s", err, out)
	}
	backups, err := filepath.Glob(filepath.Join(f.codex, "skills.*.bak"))
	if err != nil || len(backups) != 1 {
		t.Fatalf("backup: %v, %v", backups, err)
	}
	copy := filepath.Join(backups[0], "Local/locked/deeper/content.txt")
	if f.read(copy) != "preserved\n" {
		t.Fatal("backup contents changed")
	}
	fileInfo, err := os.Stat(copy)
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.Mode().Perm() != 0755 {
		t.Fatalf("backup file mode %o, want 755", fileInfo.Mode().Perm())
	}
	for _, p := range []string{filepath.Join(backups[0], "Local/locked"), filepath.Dir(copy)} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0555 {
			t.Fatalf("mode of %s is %o, want 555", p, info.Mode().Perm())
		}
	}
}

func TestContractRejectMismatchedSkillNameQuotes(t *testing.T) {
	for _, name := range []string{`"Alpha'`, `'Alpha"`, `"Alpha`, `Alpha"`} {
		t.Run(name, func(t *testing.T) {
			f := newContractFixture(t)
			f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: "+name+"\ndescription: Example\n---\n")
			f.run("render", 2)
			f.absent(f.render)
		})
	}
}

func TestContractMarkdownEscapedBackslashResource(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "resources/guide/first.md", "[literal](literal\\\\*.md)\n")
	f.put(f.layer, `resources/guide/literal\*.md`, "correct target\n")
	f.put(f.layer, "resources/guide/literal*.md", "wrong target\n")
	f.run("render", 0)
	bundled := filepath.Join(f.latest(), "skills/Alpha/resources/guide")
	if f.read(filepath.Join(bundled, `literal\*.md`)) != "correct target\n" {
		t.Fatal("escaped target missing")
	}
	f.absent(filepath.Join(bundled, "literal*.md"))
}

func TestContractRenderRespectsRestrictiveUmask(t *testing.T) {
	f := newContractFixture(t)
	args := []string{"-c", `umask 077; exec "$@"`, "sh", contractBinary, "render", "--config-dir", f.config, "--render-dir", f.render, "--codex-dir", f.codex}
	out, err := exec.Command("sh", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}
	info, err := os.Stat(f.latest())
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0500 {
		t.Fatalf("generation mode %o, want 500", got)
	}
}

func TestContractResourceReferenceDecodesValidPercentBeforeInvalid(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\n{{%_resources_%}}/guide/a%20b.md %\n")
	f.put(f.layer, "resources/guide/a b.md", "decoded\n")
	f.run("render", 0)
	if f.read(filepath.Join(f.latest(), "skills/Alpha/resources/guide/a b.md")) != "decoded\n" {
		t.Fatal("percent-encoded resource missing")
	}
}

func TestContractPruneFailurePreservesGenerationMode(t *testing.T) {
	f := newContractFixture(t)
	for i := 0; i < 6; i++ {
		f.put(f.layer, "AGENTS.md.d/10-base.md", fmt.Sprintf("Generation %d\n", i))
		f.run("render", 0)
	}
	entries, err := os.ReadDir(f.render)
	if err != nil {
		t.Fatal(err)
	}
	var first string
	for _, e := range entries {
		if e.IsDir() {
			first = filepath.Join(f.render, e.Name())
			break
		}
	}
	if first == "" {
		t.Fatal("no generation")
	}
	info, err := os.Stat(first)
	if err != nil {
		t.Fatal(err)
	}
	before := info.Mode().Perm()
	if err := os.Chmod(f.render, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(f.render, 0755)
	probe := filepath.Join(f.render, ".permission-probe")
	if file, err := os.Create(probe); err == nil {
		file.Close()
		os.Remove(probe)
		t.Skip("current user can write a read-only directory")
	}
	f.run("prune", 2)
	info, err = os.Stat(first)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != before {
		t.Fatalf("generation mode changed from %o to %o after failure", before, info.Mode().Perm())
	}
}
