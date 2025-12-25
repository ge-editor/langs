// Colorize Go source by Tree-sitter
// $ go get github.com/smacker/go-tree-sitter
// $ go get github.com/smacker/go-tree-sitter-go

package go_mode

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/gdamore/tcell/v3"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"

	"github.com/ge-editor/gecore/lang"

	"github.com/ge-editor/langs/abstract"

	"github.com/ge-editor/theme"
)

// URL 検出用の正規表現
var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

var parser *sitter.Parser

func init() {
	// Setup Tree-sitter parser
	parser = sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage()) // for Go language
}

func (gm *GoMode) ColorizeEvents(ctx context.Context, oldTree *sitter.Tree, sourceCode []byte) ([]lang.Event, *sitter.Tree, error) {
	// var oldTree *sitter.Tree
	// oldTree := &sitter.Tree{}

	fmt.Printf("%v,%v\n", ctx, oldTree)
	// 構文解析
	tree, err := parser.ParseCtx(ctx, oldTree, sourceCode)
	if err != nil {
		return nil, tree, err
	}
	rootNode := tree.RootNode()
	// fmt.Printf("rootNode: %v\n--------\n", rootNode)

	// パース結果を位置データに変換
	events := parseToEvents(sourceCode, rootNode)
	// fmt.Printf("events: %v\n--------\n", events)

	// 行、列でソート
	// 行、列で比較し、"end" イベントを優先する
	sort.Slice(events, func(i, j int) bool {
		if events[i].Row == events[j].Row {
			if events[i].Column == events[j].Column {
				return events[i].EventType == "end" // "end" を優先
			}
			return events[i].Column < events[j].Column
		}
		return events[i].Row < events[j].Row
	})
	fmt.Printf("sorted events: %v\n--------\n", events)

	return events, tree, nil
}

// Returns the event index, tcell.Style, and error that match the current source row and column.
// 現在のソース行と列に一致するイベントインデックス、tcell.Style、及びエラーを返す。
// eventIndex := 0                    // 現在処理中のイベントインデックス
func (gm *GoMode) EventIndex(ctx context.Context, currentRow, currentCol int, source [][]byte, events []lang.Event, eventIndex int) (int, tcell.Style, error) {
	cs := abstract.ColorStack{}
	currentColor := theme.ColorDefault //  "\033[0m" // 初期はリセット色

	for row, line := range source {
		for col := 0; col < len(line); col++ {
			// 現在の行と列に該当するイベントを処理
			for eventIndex < len(events) && events[eventIndex].Row == uint32(row) && events[eventIndex].Column == uint32(col) {
				if events[eventIndex].EventType == "start" {
					cs.PushColor(currentColor)
					currentColor = events[eventIndex].Color
				} else if events[eventIndex].EventType == "end" {
					currentColor = cs.PopColor()
				}
				eventIndex++
				// if row >= currentRow || (row == currentRow && col >= currentCol) {
				if row == currentRow && col == currentCol {
					return eventIndex, currentColor, nil
				}
			}
		}
	}
	return -1, currentColor, fmt.Errorf("event not found")
}

// Returns the event index, tcell.Style, and error that match the current source row and column.
// 現在のソース行と列に一致するイベントインデックス、tcell.Style、及びエラーを返す。
func (gm *GoMode) EventIndex_0(ctx context.Context, currentRow, currentCol int, events []lang.Event) (int, tcell.Style, error) {
	cs := abstract.ColorStack{}
	currentColor := theme.ColorDefault // Reset color
	eventIndex := 0                    // 現在処理中のイベントインデックス

	for row := 0; row <= currentRow; row++ {
		for col := 0; col <= currentCol; col++ {
			// contextの中断処理をチェック
			select {
			case <-ctx.Done():
				return -1, currentColor, ctx.Err() // abort
			default:
				// 現在の行と列に該当するイベントを処理
				for eventIndex < len(events) && events[eventIndex].Row == uint32(row) && events[eventIndex].Column == uint32(col) {
					if events[eventIndex].EventType == "start" {
						cs.PushColor(currentColor)
						currentColor = events[eventIndex].Color
						return eventIndex, currentColor, nil
					} else if events[eventIndex].EventType == "end" {
						currentColor = cs.PopColor()
						return eventIndex, currentColor, nil
					}
					eventIndex++
				}
			}
		}
	}
	return -1, currentColor, fmt.Errorf("event not found")
}

// 構文木を走査してコールバックを実行
func walk(node *sitter.Node, source []byte, callback func(node *sitter.Node)) {
	if node == nil {
		return
	}

	// 現在のノードでコールバックを実行
	callback(node)

	// 子ノードを再帰的に処理
	for i := 0; i < int(node.ChildCount()); i++ {
		walk(node.Child(i), source, callback)
	}
}

// パース結果をスイープライン イベントに変換
func parseToEvents(source []byte, node *sitter.Node) []lang.Event {
	var events []lang.Event

	walk(node, source, func(currentNode *sitter.Node) {
		nodeType := currentNode.Type()
		color, err := getColor(nodeType)
		if err != nil {
			return
		}

		startPoint := currentNode.StartPoint()
		endPoint := currentNode.EndPoint()

		// コメントや文字列リテラル内にURLがある場合
		// Tree-sitterの構文木をさらに詳細に解析する方法もあるけれどやめ
		if urlColor, err := getColor("url"); err == nil && (nodeType == "comment" || nodeType == "interpreted_string_literal") {
			text := currentNode.Content(source)
			matches := urlRegex.FindAllStringIndex(text, -1)
			for _, match := range matches {
				urlStart := match[0]
				urlEnd := match[1]
				startRow, startCol := calculatePosition(text, currentNode, urlStart)
				endRow, endCol := calculatePosition(text, currentNode, urlEnd)

				events = append(events, lang.Event{
					Row:       startRow,
					Column:    startCol,
					EventType: "start",
					Color:     urlColor,
					NodeType:  "url",
				}, lang.Event{
					Row:       endRow,
					Column:    endCol,
					EventType: "end",
					Color:     urlColor,
					NodeType:  "url",
				})
			}
		}

		events = append(events, lang.Event{
			Row:       startPoint.Row,
			Column:    startPoint.Column,
			EventType: "start",
			Color:     color,
			NodeType:  nodeType,
		})
		events = append(events, lang.Event{
			Row:       endPoint.Row,
			Column:    endPoint.Column,
			EventType: "end",
			Color:     color,
			NodeType:  nodeType,
		})
	})

	return events
}

// ノード内の相対位置を計算して絶対位置に変換
func calculatePosition(text string, node *sitter.Node, byteIndex int) (uint32, uint32) {
	lines := bytes.Split([]byte(text[:byteIndex]), []byte("\n"))
	row := node.StartPoint().Row + uint32(len(lines)-1)
	col := uint32(len(lines[len(lines)-1]))
	if len(lines) == 1 {
		col += node.StartPoint().Column
	}
	return row, col
}

// ノードの種類に応じた色を取得
func getColor(nodeType string) (tcell.Style, error) {
	if color, ok := theme.CodeColors[nodeType]; ok {
		return color, nil
	}
	return theme.CodeColors["default"], fmt.Errorf("color not found for node type: %s", nodeType)
}

// Sweep Line スイープラインでソースコードに色付け
/*
func colorizeSource(source []byte, events []Event) {
	currentColor := theme.ColorDefault // "\033[0m" // 初期はリセット色
	eventIndex := 0                    // 現在処理中のイベントインデックス

	for row, line := range splitLines(source) {
		coloredLine := ""

		for col := 0; col < len(line); col++ {
			// 現在の行と列に該当するイベントを処理
			for eventIndex < len(events) && events[eventIndex].Row == uint32(row) && events[eventIndex].Column == uint32(col) {
				if events[eventIndex].EventType == "start" {
					pushColor(currentColor)
					currentColor = events[eventIndex].Color
				} else if events[eventIndex].EventType == "end" {
					// currentColor = "\033[0m" // リセット色
					currentColor = popColor()
				}
				eventIndex++
			}

			// 現在の文字に色を付ける
			coloredLine += currentColor + string(line[col])
		}

		// 行の終了後にリセット
		fmt.Printf("%s\033[0m", coloredLine)
		// fmt.Printf("%s", coloredLine)
	}
}
*/
