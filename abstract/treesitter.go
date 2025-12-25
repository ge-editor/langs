package abstract

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v3"
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/ge-editor/gecore/lang"

	"github.com/ge-editor/theme"
)

func (aw *Worker) ColorizeEvents(ctx context.Context, oldTree *sitter.Tree, sourceCode []byte) ([]lang.Event, *sitter.Tree, error) {
	return nil, nil, nil
}

// Returns the event index, tcell.Style, and error that match the current source row and column.
// 現在のソース行と列に一致するイベントインデックス、tcell.Style、及びエラーを返す。
func (aw *Worker) EventIndex(ctx context.Context, currentRow, currentCol int, source [][]byte, events []lang.Event, eventIndex int) (int, tcell.Style, error) {
	return -1, theme.ColorDefault, fmt.Errorf("event not found")
}

// スイープラインの処理でネストした色付けを処理
type ColorStack []tcell.Style

func (cs *ColorStack) PushColor(color tcell.Style) {
	*cs = append(*cs, color)
}

func (cs *ColorStack) PopColor() tcell.Style {
	if len(*cs) == 0 {
		return theme.ColorDefault // Reset color
	}
	color := (*cs)[len(*cs)-1] // スタックの末尾を取得
	*cs = (*cs)[:len(*cs)-1]   // スタックから末尾を削除
	return color
}
