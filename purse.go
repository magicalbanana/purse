package purse

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
	ext string = ".sql"
)

// Purse is the interface representing a collection of SQL files whose
// contents are accessed via the basic Get method.
//
// Implementations of this interface should be safe for concurrent use by
// multiple goroutines.
type Purse interface {
	Get(string) (string, bool)
}

// MemoryPurse is an implementation of Purse that uses a map of strings to
// represent the contents of SQL files found within a specified directory.
// It is safe for concurrent use by multiple goroutines.
type MemoryPurse struct {
	mu    sync.RWMutex
	files map[string]string
}

// New loads SQL files' contents in the specified directory dir into memory.
// A file's contents can be accessed with the Get method.
//
// New returns an error if the directory does not exist or is not a directory.
func New(dir string) (*MemoryPurse, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fis, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	p := &MemoryPurse{files: make(map[string]string, 0)}

	for _, fi := range fis {
		if !fi.IsDir() && filepath.Ext(fi.Name()) == ext {
			f, err := os.Open(filepath.Join(dir, fi.Name()))
			if err != nil {
				return nil, err
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
			p.files[fi.Name()] = string(b)
			f.Close()
		}
	}
	return p, nil
}

// Get returns a SQL file's contents as a string.
// If the file does not exist or does exists but had a read error,
// then v == "" and ok == false.
func (p *MemoryPurse) Get(filename string) (v string, ok bool) {
	p.mu.RLock()
	v, ok = p.files[filename]
	p.mu.RUnlock()
	return
}

// Files returns a slice of filenames for all loaded SQL files.
func (p *MemoryPurse) Files() []string {
	fs := make([]string, len(p.files))
	i := 0
	for k, _ := range p.files {
		fs[i] = k
		i++
	}
	return fs
}
