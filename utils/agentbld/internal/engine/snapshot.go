package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const SnapshotFile = "config.json"

// Snapshot preserves every selected harness under its top-level name.
type Snapshot map[string]any

func SnapshotBytes(s Snapshot) ([]byte, error) {
	b, e := json.Marshal(s)
	if e != nil {
		return nil, e
	}
	return append(b, '\n'), nil
}

// ReadSnapshot accepts only the canonical format written by this version. This
// rejects duplicate fields, unknown fields, omitted fields, and trailing data.
func (registry Registry) ReadSnapshot(gen string) (Snapshot, error) {
	var s Snapshot
	path := filepath.Join(gen, SnapshotFile)
	data, e := os.ReadFile(path)
	if os.IsNotExist(e) {
		return nil, fmt.Errorf("generation has no config.json; rerender it before syncing")
	}
	if e != nil {
		return s, e
	}
	s, e = DecodeSnapshot(data, registry)
	if e != nil {
		return nil, fmt.Errorf("invalid generation config %s: %w", path, e)
	}
	return s, nil
}

// DecodeSnapshot preserves every registered harness and rejects unknown fields,
// duplicate keys, missing fields, and noncanonical JSON.
func DecodeSnapshot(data []byte, registry map[string]Handler) (Snapshot, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil || len(raw) == 0 {
		return nil, fmt.Errorf("invalid generation config")
	}
	s := Snapshot{}
	for _, name := range SortedKeys(raw) {
		handler, ok := registry[name]
		if !ok {
			return nil, fmt.Errorf("unsupported harness in generation config: %s", name)
		}
		value, err := handler.Decode(raw[name])
		if err != nil {
			return nil, fmt.Errorf("invalid %s generation config: %w", name, err)
		}
		s[name] = value
	}
	canonical, err := SnapshotBytes(s)
	if err != nil || !bytes.Equal(data, canonical) {
		return nil, fmt.Errorf("invalid or unsupported generation config")
	}
	return s, nil
}

// RestorePlan reconstructs a harness plan solely from its immutable payload.
func RestorePlan(gen string, fileModes map[string]uint32, layers, fragments int, instructionsName string) (*RenderPlan, error) {
	p := &RenderPlan{InstructionsName: instructionsName, Skills: map[string]string{}, Resources: map[string]string{}, Dirs: map[string]bool{}, Files: map[string]string{}, Contents: map[string][]byte{}, fileModes: map[string]os.FileMode{}}
	p.Layers = make([]string, layers)
	p.Fragments = fragments
	b, e := os.ReadFile(filepath.Join(gen, instructionsName))
	if e != nil {
		return nil, e
	}
	p.Agents = string(b)
	e = walk(gen, func(path string, d os.DirEntry) error {
		rel, e := filepath.Rel(gen, path)
		if e != nil {
			return e
		}
		if rel == instructionsName {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("generation contains symlink: %s", path)
		}
		if d.IsDir() {
			p.Dirs[rel] = true
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("generation contains non-regular file: %s", path)
		}
		mode, ok := fileModes[filepath.ToSlash(rel)]
		if !ok {
			return fmt.Errorf("generation file is absent from config: %s", path)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if info.Mode().Perm() != os.FileMode(mode)&^0222 {
			return fmt.Errorf("generation file mode does not match config: %s", path)
		}
		p.Files[rel] = path
		p.fileModes[rel] = os.FileMode(mode)
		return nil
	})
	if e != nil {
		return nil, e
	}
	if len(p.Files) != len(fileModes) {
		return nil, fmt.Errorf("generation files do not match config: %s", gen)
	}
	skills := filepath.Join(gen, "skills")
	es, e := Entries(skills)
	if e != nil {
		return nil, e
	}
	for _, d := range es {
		path := filepath.Join(skills, d.Name())
		if !d.IsDir() {
			return nil, fmt.Errorf("invalid generated skill: %s", path)
		}
		name, e := ReadSkillName(filepath.Join(path, "SKILL.md"))
		if e != nil {
			return nil, e
		}
		p.Skills[name] = path
	}
	for rel, source := range p.Files {
		if strings.HasPrefix(rel, "resources"+string(filepath.Separator)) {
			p.Resources[strings.TrimPrefix(rel, "resources"+string(filepath.Separator))] = source
		}
	}
	return p, nil
}
