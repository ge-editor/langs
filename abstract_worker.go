package langs

import (
	"path/filepath"
	"strings"
)

// ------------------------------------------------------------------
// AbstractWorker implement gecore Mode interface
// ------------------------------------------------------------------

type AbstractWorker struct {
	exts []string
}

/*
func (aw *AbstractWorker) Name() string {
	return ""
}
*/

// Matches checks if the given file path matches any of the extensions in GoMode.exts
func (aw *AbstractWorker) HasMatchingExtension(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath)) // Get the file extension in lowercase
	for _, validExt := range aw.exts {
		if ext == strings.ToLower(validExt) { // Case-insensitive comparison
			return true
		}
	}
	return false
}

// Format formats source code and restores the cursor position.
// The cursor position based on byte indices
func (aw *AbstractWorker) Format(source [][]byte /* , cursorRow, cursorCol int */) ([][]byte /*  int, int, */, error) {
	return source /* cursorRow, cursorCol, */, nil
}

func (aw *AbstractWorker) FormatBeforeSave(source [][]byte /* , cursorRow, cursorCol int */) ([][]byte /* int, int, */, error) {
	return aw.Format(source /* , cursorRow, cursorCol */)
}

func (aw *AbstractWorker) IndentWidth() int {
	return 4
}

func (aw *AbstractWorker) IsSoftTAB() bool {
	return false // Hard TAB
}
