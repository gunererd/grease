package buffer

import (
	"os"
	"path/filepath"
	"sort"
)

type Entry struct {
	Name    string
	IsDir   bool
	Path    string
	Size    int64
	ModTime int64
}

type Buffer struct {
	Lines       []string
	Entries     []Entry
	LineToEntry map[int]int
	CurrentDir  string
}

func NewBuffer() *Buffer {
	return &Buffer{
		Lines:       make([]string, 0),
		Entries:     make([]Entry, 0),
		LineToEntry: make(map[int]int),
		CurrentDir:  ".",
	}
}

// ReadDirectory reads the content of the specified directory into the buffer
func (b *Buffer) ReadDirectory(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	b.Lines = make([]string, 0)
	b.Entries = make([]Entry, 0)
	b.LineToEntry = make(map[int]int)
	b.CurrentDir = absPath

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		b.Entries = append(b.Entries, Entry{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Path:    filepath.Join(absPath, entry.Name()),
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		})
	}

	// Sort entries: directories first, then files, both alphabetically
	sort.Slice(b.Entries, func(i, j int) bool {
		if b.Entries[i].IsDir != b.Entries[j].IsDir {
			return b.Entries[i].IsDir
		}
		return b.Entries[i].Name < b.Entries[j].Name
	})

	for i, entry := range b.Entries {
		suffix := ""
		if entry.IsDir {
			suffix = "/"
		}
		b.Lines = append(b.Lines, entry.Name+suffix)
		b.LineToEntry[i] = i
	}

	return nil
}

func (b *Buffer) GetLine(idx int) string {
	if idx < 0 || idx >= len(b.Lines) {
		return ""
	}
	return b.Lines[idx]
}

func (b *Buffer) NumLines() int {
	return len(b.Lines)
}

func (b *Buffer) GetEntry(lineNum int) (Entry, bool) {
	if entryIdx, ok := b.LineToEntry[lineNum]; ok {
		return b.Entries[entryIdx], true
	}
	return Entry{}, false
}

func (b *Buffer) GetCurrentDir() string {
	return b.CurrentDir
}

func (b *Buffer) GetParentDir() (string, error) {
	parentDir := filepath.Dir(b.CurrentDir)

	if _, err := os.Stat(parentDir); err != nil {
		return "", err
	}

	return parentDir, nil
}

// GetLineLength returns the length of the line at the given index
// If the index is out of bounds, returns 0
func (b *Buffer) GetLineLength(idx int) int {
	if idx < 0 || idx >= len(b.Lines) {
		return 0
	}
	return len(b.Lines[idx])
}
