package go_mode

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/ge-editor/gecore/define"
	"github.com/ge-editor/gecore/lang"

	"github.com/ge-editor/utils"

	"github.com/ge-editor/langs/abstract"
)

// ------------------------------------------------------------------
// GoMode implement gecore Mode interface
// ------------------------------------------------------------------

func NewGoMode() lang.Mode {
	return &GoMode{
		Worker: abstract.Worker{
			// Initialize the `exts` field of AbstractWorker
			// allowed extensions
			Exts: []string{".go" /* , ".mod", ".sum" */},
		},
	}
}

type GoMode struct {
	abstract.Worker
}

func (gm *GoMode) Name() string {
	return "Go"
}

// Format formats source code and restores the cursor position.
// The cursor position based on byte indices
func (gm *GoMode) Format(source [][]byte /* cursorRow, cursorCol int */) ([][]byte /* int, int, */, error) {
	// Combine [][]byte into a single byte slice
	sourceBytes := utils.JoinBytes(source)

	// Remove EOF Mark
	sourceBytes = sourceBytes[:len(sourceBytes)-1]

	// Format the source code
	formattedBytes, err := format.Source(sourceBytes)
	if err != nil {
		return nil /* 0, 0, */, fmt.Errorf("failed to format source: %w", err)
	}

	// Add EOF Mark
	formattedBytes = append(formattedBytes, define.EOF)

	// Split formatted source back into [][]byte
	formattedLines := bytes.SplitAfter(formattedBytes, []byte("\n"))

	// Restore the cursor position
	// newCursorRow, newCursorCol := restoreCursorPosition(source, formattedLines, cursorRow, cursorCol)

	return formattedLines /* newCursorRow, newCursorCol, */, nil
}

func (gm *GoMode) FormatBeforeSave(source [][]byte /* , cursorRow, cursorCol int */) ([][]byte /* int, int, */, error) {
	return gm.Format(source /* , cursorRow, cursorCol */)
}

// Restore the cursor position based on byte indices
/*
func restoreCursorPosition(
	original [][]byte,
	formatted [][]byte,
	cursorRow int,
	cursorCol int,
) (int, int) {
	newCursorRow := cursorRow
	newCursorCol := cursorCol

	// Ensure the cursor row is within the bounds of the formatted lines
	if cursorRow < len(formatted) {
		originalLine := original[cursorRow]
		formattedLine := formatted[cursorRow]

		// Extract the portion of the original line up to the cursor column (byte-based)
		if cursorCol <= len(originalLine) {
			originalPrefix := originalLine[:cursorCol]
			// Find the position of the original prefix in the formatted line
			newCursorCol = bytes.Index(formattedLine, originalPrefix)
			if newCursorCol == -1 {
				newCursorCol = len(formattedLine) // Fallback to the end of the line if not found
			}
		} else {
			// If the cursor column exceeds the line length, set it to the end of the formatted line
			newCursorCol = len(formattedLine)
		}
	} else {
		// If the cursor row exceeds the number of formatted lines, set it to the last line
		newCursorRow = len(formatted) - 1
		newCursorCol = len(formatted[newCursorRow])
	}

	return newCursorRow, newCursorCol
}
*/
