package buffer

import (
	"bufio"
	"errors"
	"io"
	"log"
	"regexp"
	"sort"
	"sync"
	"unicode/utf8"

	"github.com/gunererd/grease/internal/types"
)

var (
	ErrInvalidOffset = errors.New("invalid offset")
	ErrInvalidLine   = errors.New("invalid line number")
	ErrNoCursor      = errors.New("no cursor available")

	wordPattern    = regexp.MustCompile(`\w+`)
	bigWordPattern = regexp.MustCompile(`\S+`)
)

// Buffer represents the text content and provides operations to modify it
type Buffer struct {
	lines   [][]rune // each line is a slice of runes
	cursors []types.Cursor
	mu      sync.RWMutex

	nextCursorID int
}

// New creates a new empty buffer
func New() *Buffer {
	b := &Buffer{
		lines: [][]rune{{}}, // start with one empty line
	}
	// Create primary cursor at start of buffer
	b.AddCursor(NewPosition(0, 0), 100) // Primary cursor gets high priority
	return b
}

// NewFromString creates a buffer from a string
func NewFromString(content string) *Buffer {
	b := &Buffer{}
	if content == "" {
		b.lines = [][]rune{{}}
	} else {
		var lines [][]rune
		currentLine := make([]rune, 0, 128)

		for len(content) > 0 {
			r, size := utf8.DecodeRuneInString(content)
			content = content[size:]

			if r == '\n' {
				lines = append(lines, currentLine)
				currentLine = make([]rune, 0, 128)
			} else {
				currentLine = append(currentLine, r)
			}
		}
		lines = append(lines, currentLine)
		b.lines = lines
	}

	// Create primary cursor at start of buffer
	b.AddCursor(NewPosition(0, 0), 100)
	return b
}

// LoadFromReader loads content from an io.Reader into the buffer
func (b *Buffer) LoadFromReader(r io.Reader) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	scanner := bufio.NewScanner(r)
	b.lines = make([][]rune, 0)

	for scanner.Scan() {
		b.lines = append(b.lines, []rune(scanner.Text()))
	}

	// Ensure at least one line exists
	if len(b.lines) == 0 {
		b.lines = append(b.lines, []rune(""))
	}

	return scanner.Err()
}

// AddCursor adds a new cursor at the specified position
func (b *Buffer) AddCursor(pos types.Position, priority int) (types.Cursor, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err := b.validatePosition(pos); err != nil {
		return nil, err
	}

	cursor := NewCursor(pos, b.nextCursorID, priority)
	b.nextCursorID++

	// Insert cursor in priority order
	insertIdx := sort.Search(len(b.cursors), func(i int) bool {
		return b.cursors[i].GetPriority() <= priority
	})

	b.cursors = append(b.cursors, nil)
	copy(b.cursors[insertIdx+1:], b.cursors[insertIdx:])
	b.cursors[insertIdx] = cursor

	return cursor, nil
}

// RemoveCursor removes a cursor by its ID
func (b *Buffer) RemoveCursor(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, c := range b.cursors {
		if c.ID() == id {
			b.cursors = append(b.cursors[:i], b.cursors[i+1:]...)
			return
		}
	}
}

// GetPrimaryCursor returns the highest priority cursor
func (b *Buffer) GetPrimaryCursor() (types.Cursor, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.cursors) == 0 {
		return nil, ErrNoCursor
	}
	return b.cursors[0], nil
}

func (b *Buffer) MoveCursorRelative(cursorID int, lineOffset, columnOffset int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	cursor := b.findCursor(cursorID)
	if cursor == nil {
		return ErrNoCursor
	}

	newPos := cursor.GetPosition().Add(lineOffset, columnOffset)

	// Validate line bounds
	if newPos.Line() < 0 || newPos.Line() >= len(b.lines) {
		log.Println("Invalid line:", newPos.Line())
		return ErrInvalidLine
	}

	// When moving vertically, if target line is shorter than current column,
	// move cursor to end of the target line
	if lineOffset != 0 && columnOffset == 0 {
		if newPos.Column() > len(b.lines[newPos.Line()]) {
			col := len(b.lines[newPos.Line()])
			newPos = NewPosition(newPos.Line(), col)
		}
	} else {
		// For horizontal movement, validate column bounds normally
		if newPos.Column() < 0 || newPos.Column() > len(b.lines[newPos.Line()]) {
			return ErrInvalidOffset
		}
	}

	cursor.SetPosition(newPos)
	return nil
}

func (b *Buffer) MoveCursor(cursorID int, lineOffset, columnOffset int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	cursor := b.findCursor(cursorID)
	if cursor == nil {
		return ErrNoCursor
	}

	newPos := NewPosition(lineOffset, columnOffset)

	// Validate line bounds
	if newPos.Line() < 0 || newPos.Line() >= len(b.lines) {
		log.Println("Invalid line:", newPos.Line())
		return ErrInvalidLine
	}

	// When moving vertically, if target line is shorter than current column,
	// move cursor to end of the target line
	if lineOffset != 0 && columnOffset == 0 {
		if newPos.Column() > len(b.lines[newPos.Line()]) {
			col := len(b.lines[newPos.Line()])
			newPos = NewPosition(newPos.Line(), col)
		}
	} else {
		// For horizontal movement, validate column bounds normally
		if newPos.Column() < 0 || newPos.Column() > len(b.lines[newPos.Line()]) {
			return ErrInvalidOffset
		}
	}

	cursor.SetPosition(newPos)
	return nil
}

// Insert inserts text at all cursor positions
func (b *Buffer) Insert(text string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.cursors) == 0 {
		return ErrNoCursor
	}

	// Sort cursors in reverse order to handle insertions from bottom to top
	cursors := make([]types.Cursor, len(b.cursors))
	copy(cursors, b.cursors)
	sort.Slice(cursors, func(i, j int) bool {
		return cursors[j].GetPosition().Before(cursors[i].GetPosition())
	})

	for _, cursor := range cursors {
		if err := b.insertAt(cursor.GetPosition(), text); err != nil {
			return err
		}
		// Update cursor position
		if text == "\n" {
			cursor.SetPosition(NewPosition(cursor.GetPosition().Line()+1, 0))
		} else {
			col := cursor.GetPosition().Column() + len([]rune(text))
			cursor.SetPosition(NewPosition(cursor.GetPosition().Line(), col))
		}
	}

	return nil
}

// insertAt inserts text at a specific position (internal method)
func (b *Buffer) insertAt(pos types.Position, text string) error {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	if text == "\n" {
		// Handle newline insertion
		line := b.lines[pos.Line()]
		newLine := append([]rune{}, line[pos.Column():]...)
		b.lines[pos.Line()] = line[:pos.Column()]
		b.lines = append(b.lines[:pos.Line()+1], append([][]rune{newLine}, b.lines[pos.Line()+1:]...)...)
		return nil
	}

	// Handle regular text insertion
	line := b.lines[pos.Line()]
	b.lines[pos.Line()] = append(line[:pos.Column()], append(runes, line[pos.Column():]...)...)
	return nil
}

// Delete deletes the specified number of characters at all cursor positions
func (b *Buffer) Delete(count int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.cursors) == 0 {
		return ErrNoCursor
	}

	// Sort cursors in reverse order to handle deletions from bottom to top
	cursors := make([]types.Cursor, len(b.cursors))
	copy(cursors, b.cursors)
	sort.Slice(cursors, func(i, j int) bool {
		return cursors[j].GetPosition().Before(cursors[i].GetPosition())
	})

	for _, cursor := range cursors {
		pos := cursor.GetPosition()
		// For backspace (count < 0), move cursor left first
		if count < 0 {
			if pos.Column() > 0 {
				newPos := NewPosition(pos.Line(), pos.Column()-1)
				cursor.SetPosition(newPos)
				if err := b.deleteAt(newPos, 1); err != nil {
					return err
				}
			}
		} else {
			// For delete (count > 0), delete at current position
			if err := b.deleteAt(pos, count); err != nil {
				return err
			}
		}
	}

	return nil
}

// deleteAt deletes characters at a specific position (internal method)
func (b *Buffer) deleteAt(pos types.Position, count int) error {
	line := b.lines[pos.Line()]
	if pos.Column()+count <= len(line) {
		// Simple case: deletion within the same line
		b.lines[pos.Line()] = append(line[:pos.Column()], line[pos.Column()+count:]...)
		return nil
	}

	// Handle multi-line deletion
	return errors.New("multi-line deletion not implemented yet")
}

// Get returns the content of the buffer as a string
func (b *Buffer) Get() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var result []rune
	for i, line := range b.lines {
		result = append(result, line...)
		if i < len(b.lines)-1 {
			result = append(result, '\n')
		}
	}
	return string(result)
}

// GetLine returns the content of a specific line
func (b *Buffer) GetLine(line int) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if line < 0 || line >= len(b.lines) {
		return "", ErrInvalidLine
	}
	return string(b.lines[line]), nil
}

// LineCount returns the number of lines in the buffer
func (b *Buffer) LineCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.lines)
}

// LineLen returns the length of a specific line
func (b *Buffer) LineLen(line int) (int, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if line < 0 || line >= len(b.lines) {
		return 0, ErrInvalidLine
	}
	return len(b.lines[line]), nil
}

// NextWordPosition moves the cursor to the start of the next word
func (b *Buffer) NextWordPosition(pos types.Position, bigWord bool) types.Position {
	b.mu.RLock()
	defer b.mu.RUnlock()

	line := pos.Line()
	col := pos.Column()

	for line < len(b.lines) {
		lineStr := string(b.lines[line][col:])
		pattern := wordPattern
		if bigWord {
			pattern = bigWordPattern
		}

		// Find the next word boundary
		loc := pattern.FindStringIndex(lineStr)

		// If we found a match on this line
		if loc != nil {
			// If we're inside a word, find the next one
			if loc[0] == 0 {
				// Look for another word after this one
				nextLoc := pattern.FindStringIndex(lineStr[loc[1]:])
				if nextLoc != nil {
					return NewPosition(line, col+loc[1]+nextLoc[0])
				}
			} else {
				// Move to the start of the found word
				return NewPosition(line, col+loc[0])
			}
		}

		// If no match on this line or at end of line, try next line
		if line+1 >= len(b.lines) {
			break
		}
		line++
		col = 0
	}

	lastLine := len(b.lines) - 1
	return NewPosition(lastLine, len(b.lines[lastLine])-1)
}

// NextWordEndPosition moves the cursor to the end of the next word
func (b *Buffer) NextWordEndPosition(pos types.Position, bigWord bool) types.Position {
	b.mu.RLock()
	defer b.mu.RUnlock()

	line := pos.Line()
	col := pos.Column()

	for line < len(b.lines) {
		lineStr := string(b.lines[line][col:])
		pattern := wordPattern
		if bigWord {
			pattern = bigWordPattern
		}

		// Find the next word
		loc := pattern.FindStringIndex(lineStr)

		// If we found a match on this line
		if loc != nil {
			// If we're inside a word
			if loc[0] == 0 {
				// First check if we're not at the end of current word
				if col+loc[1]-1 > col {
					return NewPosition(line, col+loc[1]-1)
				}
				// If we are at the end, find the next word
				nextLoc := pattern.FindStringIndex(lineStr[loc[1]:])
				if nextLoc != nil {
					return NewPosition(line, col+loc[1]+nextLoc[1]-1)
				}
			} else {
				// Move to the end of the found word
				return NewPosition(line, col+loc[1]-1)
			}
		}

		// If no match on this line, try next line
		if line+1 >= len(b.lines) {
			break
		}
		line++
		col = 0
	}

	lastLine := len(b.lines) - 1
	return NewPosition(lastLine, len(b.lines[lastLine])-1)
}

// PrevWordPosition moves the cursor to the start of the previous word
func (b *Buffer) PrevWordPosition(pos types.Position, bigWord bool) types.Position {
	b.mu.RLock()
	defer b.mu.RUnlock()

	line := pos.Line()
	col := pos.Column()

	for line >= 0 {
		// Get the text before the cursor on this line
		lineStr := string(b.lines[line][:col])
		pattern := wordPattern
		if bigWord {
			pattern = bigWordPattern
		}

		// Find all matches in the line up to cursor
		matches := pattern.FindAllStringIndex(lineStr, -1)

		if matches != nil {
			// Find the last match that starts before our cursor
			for i := len(matches) - 1; i >= 0; i-- {
				match := matches[i]
				// If cursor is after start of a word, move to its start
				if match[0] < col && col <= match[1] {
					return NewPosition(line, match[0])
				}
				// If cursor is at start of a word and it's not the first word,
				// move to start of previous word
				if match[0] == col && i > 0 {
					return NewPosition(line, matches[i-1][0])
				}
				// If cursor is before this word's start, use this word
				if match[0] < col {
					return NewPosition(line, match[0])
				}
			}
		}

		if line == 0 {
			return NewPosition(0, 0)
		}
		line--
		col = len(b.lines[line])
	}

	return NewPosition(0, 0)
}

// ReplaceLine replaces the content of a specific line
func (b *Buffer) ReplaceLine(line int, content string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if line < 0 || line >= len(b.lines) {
		return ErrInvalidLine
	}

	b.lines[line] = []rune(content)
	return nil
}

// validatePosition checks if a position is valid within the buffer
func (b *Buffer) validatePosition(pos types.Position) error {
	if pos.Line() < 0 || pos.Line() >= len(b.lines) {
		return ErrInvalidLine
	}
	if pos.Column() < 0 || pos.Column() > len(b.lines[pos.Line()]) {
		return ErrInvalidOffset
	}
	return nil
}

// findCursor finds a cursor by its ID (internal method)
func (b *Buffer) findCursor(id int) types.Cursor {
	for _, c := range b.cursors {
		if c.ID() == id {
			return c
		}
	}
	return nil
}
