package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

type testAlphaMode struct {
	Value string `json:"value"`
}
type testBetaMode struct {
	Count int `json:"count"`
}

func testModeParsers() map[string]engine.ModeParser {
	return map[string]engine.ModeParser{
		"alpha": func(raw *string) (any, error) {
			if raw == nil || *raw != `alpha:{"value":"first"}` {
				return nil, fmt.Errorf("unexpected alpha input")
			}
			return testAlphaMode{Value: "first"}, nil
		},
		"beta": func(raw *string) (any, error) {
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
	got, err := resolveModes([]string{alpha, beta}, parsers)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"alpha": testAlphaMode{Value: "first"}, "beta": testBetaMode{Count: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("modes = %#v, want %#v", got, want)
	}
	reversed, err := resolveModes([]string{beta, alpha}, parsers)
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
			_, err := resolveModes(tc.values, testModeParsers())
			if err == nil || !strings.Contains(strings.ToLower(err.Error()), tc.want) {
				t.Fatalf("got error %v, want %s", err, tc.want)
			}
		})
	}
}
