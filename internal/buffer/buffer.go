package buffer

import (
	"log"
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
	Lines         []string
	Entries       []Entry
	LineToEntry   map[int]int
	CurrentDir    string
	input         string
	ModifiedLines map[int]bool
	isDirty       bool
}

func NewBuffer() *Buffer {
	return &Buffer{
		Lines:         make([]string, 0),
		Entries:       make([]Entry, 0),
		LineToEntry:   make(map[int]int),
		CurrentDir:    ".",
		input:         "",
		ModifiedLines: make(map[int]bool),
		isDirty:       false,
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
	b.ModifiedLines = make(map[int]bool)
	b.isDirty = false

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return err
	}

	// First collect all entries
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Printf("Error getting info for %s: %v", entry.Name(), err)
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

	// Sort entries
	sort.Slice(b.Entries, func(i, j int) bool {
		// Directories come first
		if b.Entries[i].IsDir != b.Entries[j].IsDir {
			return b.Entries[i].IsDir
		}
		return b.Entries[i].Name < b.Entries[j].Name
	})

	// Initialize Lines from Entries and build LineToEntry mapping
	for i, entry := range b.Entries {
		b.Lines = append(b.Lines, entry.Name)
		b.LineToEntry[len(b.Lines)-1] = i
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

func (b *Buffer) NumEntries() int {
	return len(b.Entries)
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

// UpdateLine updates the content of a line at the given index
func (b *Buffer) UpdateLine(idx int, content string) {
	if idx >= 0 && idx < len(b.Lines) {
		b.Lines[idx] = content
	}
}

// IsLineModified returns true if the line has been modified
func (b *Buffer) IsLineModified(idx int) bool {
	return b.ModifiedLines[idx]
}

// InsertCharAtCursor inserts a character at the specified position
func (b *Buffer) InsertCharAtCursor(c string, row, col int) {
	if row >= 0 && row < len(b.Lines) {
		line := b.Lines[row]
		if col >= 0 && col <= len(line) {
			b.Lines[row] = line[:col] + c + line[col:]
			b.ModifiedLines[row] = true
			b.isDirty = true
		}
	}
}

// Input handling methods
func (b *Buffer) GetInput() string {
	return b.input
}

func (b *Buffer) AppendInputChar(c string) {
	log.Println("Appending char:", c)
	b.input += c
	log.Println("Input:", b.input)
}

func (b *Buffer) DeleteInputChar() {
	if len(b.input) > 0 {
		b.input = b.input[:len(b.input)-1]
	}
}

func (b *Buffer) ClearInput() {
	b.input = ""
}
