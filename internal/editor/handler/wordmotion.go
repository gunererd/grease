package handler

import (
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

var (
	// Matches any Unicode whitespace
	whitespaceRegex = regexp.MustCompile(`^[\s]$`)
	// Matches word characters (letters, digits, underscore)
	wordCharRegex = regexp.MustCompile(`^[\p{L}\p{N}_]$`)
)

type wordMotionType int

const (
	nextWordStart wordMotionType = iota
	nextWordEnd
	prevWordStart
)

type wordMotionCommand struct {
	motionType wordMotionType
	bigWord    bool
	changeMode bool
}

func isWordChar(r rune) bool {
	return wordCharRegex.MatchString(string(r))
}

func isWhitespace(r rune) bool {
	return whitespaceRegex.MatchString(string(r))
}

func getCharType(r rune) int {
	if isWhitespace(r) {
		return 0 // whitespace
	}
	if isWordChar(r) {
		return 1 // word character
	}
	return 2 // punctuation/symbol
}

// getTargetPosition calculates where the motion should end
func (wmc *wordMotionCommand) getTargetPosition(buf types.Buffer, curPos types.Position) types.Position {
	switch wmc.motionType {
	case nextWordStart:
		return wmc.getNextWordStartTarget(buf, curPos)
	case nextWordEnd:
		return buf.NextWordEndPosition(curPos, wmc.bigWord)
	case prevWordStart:
		return buf.PrevWordPosition(curPos, wmc.bigWord)
	}
	return curPos
}

// getNextWordStartTarget handles the complex logic for 'w' and 'W' motions
func (wmc *wordMotionCommand) getNextWordStartTarget(buf types.Buffer, curPos types.Position) types.Position {
	line, _ := buf.GetLine(curPos.Line())
	runes := []rune(line)
	col := curPos.Column()

	if col >= len(runes) {
		return curPos
	}

	if wmc.bigWord {
		// For 'W', move to next non-whitespace after whitespace
		if wmc.changeMode {
			// For 'cW', find first whitespace
			for i := col; i < len(runes); i++ {
				if isWhitespace(runes[i]) {
					return buffer.NewPosition(curPos.Line(), i)
				}
			}
			return buffer.NewPosition(curPos.Line(), len(runes))
		}

		// Skip non-whitespace
		for col < len(runes) && !isWhitespace(runes[col]) {
			col++
		}
		// Skip whitespace
		for col < len(runes) && isWhitespace(runes[col]) {
			col++
		}
	} else {
		// For 'w', handle word characters and punctuation separately
		startType := getCharType(runes[col])
		col++

		// Skip characters of the same type
		for col < len(runes) {
			currentType := getCharType(runes[col])
			if currentType != startType || isWhitespace(runes[col]) {
				break
			}
			col++
		}

		// Skip whitespace
		for col < len(runes) && isWhitespace(runes[col]) {
			col++
		}
	}

	if col >= len(runes) {
		col = len(runes)
	}

	return buffer.NewPosition(curPos.Line(), col)
}

// deleteText handles the text deletion logic
func (wmc *wordMotionCommand) deleteText(buf types.Buffer, curPos, targetPos types.Position) {
	if targetPos.Line() == curPos.Line() {
		var charsToDelete int
		if wmc.motionType == prevWordStart {
			charsToDelete = curPos.Column() - targetPos.Column()
		} else {
			charsToDelete = targetPos.Column() - curPos.Column()
			if wmc.motionType == nextWordEnd {
				charsToDelete++ // include the last character for 'e' motions
			}
		}
		buf.Delete(charsToDelete)
	}
}

func (wmc *wordMotionCommand) execute(e types.Editor) (tea.Model, tea.Cmd) {
	cursor, _ := e.Buffer().GetPrimaryCursor()
	curPos := cursor.GetPosition()

	targetPos := wmc.getTargetPosition(e.Buffer(), curPos)

	if wmc.motionType == prevWordStart {
		e.Buffer().MoveCursor(cursor.ID(), curPos.Line(), targetPos.Column())
		e.HandleCursorMovement()
	}

	wmc.deleteText(e.Buffer(), curPos, targetPos)

	if wmc.changeMode {
		e.SetMode(state.InsertMode)
	}

	return e, nil
}

// createWordMotionCommand is a factory function for word motion commands
func createWordMotionCommand(mType wordMotionType, bigWord bool, changeMode bool) func(e types.Editor) (tea.Model, tea.Cmd) {
	cmd := &wordMotionCommand{
		motionType: mType,
		bigWord:    bigWord,
		changeMode: changeMode,
	}
	return cmd.execute
}
