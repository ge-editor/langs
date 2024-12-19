package abstract

import (
	"path/filepath"
	"strings"
)

// ------------------------------------------------------------------
// AbstractWorker implement gecore Mode interface
// ------------------------------------------------------------------

type Worker struct {
	Exts []string
}

/*
func (aw *AbstractWorker) Name() string {
	return ""
}
*/

// Matches checks if the given file path matches any of the extensions in GoMode.exts
func (aw *Worker) HasMatchingExtension(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath)) // Get the file extension in lowercase
	for _, validExt := range aw.Exts {
		if ext == strings.ToLower(validExt) { // Case-insensitive comparison
			return true
		}
	}
	return false
}

// Format formats source code and restores the cursor position.
// The cursor position based on byte indices
func (aw *Worker) Format(source [][]byte /* , cursorRow, cursorCol int */) ([][]byte /*  int, int, */, error) {
	return source /* cursorRow, cursorCol, */, nil
}

func (aw *Worker) FormatBeforeSave(source [][]byte /* , cursorRow, cursorCol int */) ([][]byte /* int, int, */, error) {
	return aw.Format(source /* , cursorRow, cursorCol */)
}

func (aw *Worker) IndentWidth() int {
	return 4
}

func (aw *Worker) IsSoftTAB() bool {
	return false // Hard TAB
}
