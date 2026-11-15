package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type testAlphaMode struct {
	Value string `json:"value"`
}
type testBetaMode struct {
	Count int `json:"count"`
}

func testModeParsers() map[string]modeParser {
	return map[string]modeParser{
		"alpha": func(raw *string, _ syncFlags) (any, error) {
			if raw == nil || *raw != `alpha:{"value":"first"}` {
				return nil, fmt.Errorf("unexpected alpha input")
			}
			return testAlphaMode{Value: "first"}, nil
		},
		"beta": func(raw *string, _ syncFlags) (any, error) {
			if raw == nil || *raw != `beta:{"count":2}` {
				return nil, fmt.Errorf("unexpected beta input")
			}
			return testBetaMode{Count: 2}, nil
		},
	}
}
func TestResolveModesRetainsDistinctHarnesses(t *testing.T) {
	parsers := testModeParsers()
	alpha := `alpha:{"value":"first"}`
	beta := `beta:{"count":2}`
	got, err := resolveModes([]string{alpha, beta}, syncFlags{}, parsers)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"alpha": testAlphaMode{Value: "first"}, "beta": testBetaMode{Count: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("modes = %#v, want %#v", got, want)
	}
	reversed, err := resolveModes([]string{beta, alpha}, syncFlags{}, parsers)
	if err != nil {
		t.Fatal(err)
	}
	firstJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	secondJSON, err := json.Marshal(reversed)
	if err != nil {
		t.Fatal(err)
	}
	if string(firstJSON) != string(secondJSON) {
		t.Fatalf("harness order changed resolved JSON: %s vs %s", firstJSON, secondJSON)
	}
}
func TestResolveModesRejectsDuplicatesAndUnsupportedHarnesses(t *testing.T) {
	for _, tc := range []struct {
		values []string
		want   string
	}{
		{[]string{`alpha:{"value":"first"}`, `alpha:{"value":"first"}`}, "duplicate"},
		{[]string{`alpha:{"value":"first"}`, `gamma:{}`}, "unsupported"},
	} {
		t.Run(strings.Join(tc.values, "_"), func(t *testing.T) {
			_, err := resolveModes(tc.values, syncFlags{}, testModeParsers())
			if err == nil || !strings.Contains(strings.ToLower(err.Error()), tc.want) {
				t.Fatalf("got error %v, want %s", err, tc.want)
			}
		})
	}
}
func TestUnsupportedHarnessModeFailsBeforeWrites(t *testing.T) {
	f := modeFixture(t)
	out, err := exec.Command(contractBinary, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "test-harness:{}").CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 2 {
		t.Fatalf("expected exit 2: %v\n%s", err, out)
	}
	f.absent(f.render)
	f.absent(f.codex)
	f.absent(f.home)
}

func testSnapshotRegistry() map[string]harnessHandler {
	return map[string]harnessHandler{
		"alpha": {decode: func(data json.RawMessage) (any, error) {
			var v testAlphaMode
			err := json.Unmarshal(data, &v)
			return v, err
		}},
		"beta": {decode: func(data json.RawMessage) (any, error) {
			var v testBetaMode
			err := json.Unmarshal(data, &v)
			return v, err
		}},
	}
}
func TestSnapshotRetainsEachRegisteredHarnessAcrossGenerationAndCheck(t *testing.T) {
	want := snapshot{"alpha": testAlphaMode{Value: "first"}, "beta": testBetaMode{Count: 2}}
	metadata, err := snapshotBytes(want)
	if err != nil {
		t.Fatal(err)
	}
	plan := &renderPlan{agents: marker + "\n", dirs: map[string]bool{}, files: map[string]string{}, contents: map[string][]byte{}, skills: map[string]string{}, resources: map[string]string{}}
	root := t.TempDir()
	gen, err := performRender(root, plan, metadata)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(gen, 0700) })
	saved, err := os.ReadFile(filepath.Join(gen, snapshotFile))
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := decodeSnapshot(saved, testSnapshotRegistry())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded snapshot = %#v, want %#v", decoded, want)
	}
	if !renderedMatches(gen, plan, metadata) {
		t.Fatal("matching harness configuration did not match generation")
	}
	withoutBeta, err := snapshotBytes(snapshot{"alpha": testAlphaMode{Value: "first"}})
	if err != nil {
		t.Fatal(err)
	}
	if renderedMatches(gen, plan, withoutBeta) {
		t.Fatal("check ignored removed harness metadata")
	}
}
func TestSnapshotRejectsUnknownAndDuplicateHarnesses(t *testing.T) {
	for _, data := range []string{`{"alpha":{"value":"first"},"gamma":{}}` + "\n", `{"alpha":{"value":"first"},"alpha":{"value":"first"}}` + "\n"} {
		if _, err := decodeSnapshot([]byte(data), testSnapshotRegistry()); err == nil {
			t.Fatalf("accepted invalid snapshot: %s", data)
		}
	}
}
