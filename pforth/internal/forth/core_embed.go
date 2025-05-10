package forth

import (
	"os"
	"path/filepath"
	"strings"
)

var searchDirs = []string{"libs", "demos"}

func findLibPath(name string) string {
	lower := strings.ToLower(name) + ".fs"
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		for _, sd := range searchDirs {
			candidate := filepath.Join(dir, sd, lower)
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func (f *Forth) requireModule(name string) {
	upper := strings.ToUpper(name)
	if f.Loaded[upper] {
		return
	}
	f.Loaded[upper] = true

	path := findLibPath(name)
	if path == "" {
		forthError("MODULE NOT FOUND: %s", name)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		forthError("MODULE NOT FOUND: %s", name)
	}
	savedTIB := f.TIB
	savedIN := f.IN
	f.LoadSystem(string(data))
	f.TIB = savedTIB
	f.IN = savedIN
}

func (f *Forth) LoadFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		forthError("CANNOT OPEN: %s", path)
	}
	savedTIB := f.TIB
	savedIN := f.IN
	f.LoadSystem(string(data))
	f.TIB = savedTIB
	f.IN = savedIN
}

func (f *Forth) LoadSystem(data string) {
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "\\") {
			continue
		}
		if strings.HasPrefix(line, "#REQUIRE") {
			rest := strings.TrimSpace(line[8:])
			if rest != "" {
				f.requireModule(rest)
			}
			continue
		}
		f.InterpretLine(line)
	}
}
