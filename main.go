package main

import (
	"fmt"
	"os"

	"github.com/alex-laycalvert/vi-minus-minus/buffer"
	"github.com/gdamore/tcell"
)

const (
	Normal = iota
	Insert = iota
)

func main() {
	scr, err := tcell.NewScreen()
	checkError(err)
	defer quit(scr)
	err = scr.Init()
	checkError(err)
	normalStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	insertStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Blink(true)
	mode := Normal
	style := normalStyle

	buf := buffer.New()
	buf.AppendLine("")
	sideBarLen := 1 + len(fmt.Sprintf("%v", buf.Len()))
	col := 0
	row := 0
	cols, rows := scr.Size()
	for {
		sideBarLen = 1 + len(fmt.Sprintf("%v", buf.Len()))
		if mode == Normal {
			style = normalStyle
		} else if mode == Insert {
			style = insertStyle
		}
		// scr.Clear()
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				scr.SetContent(c, r, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
			}
		}
		cols, rows = scr.Size()
		iter := buf.Iter()
		for iter.HasMore() {
			r, l := iter.Next()
			drawText(scr, 0, r, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue), fmt.Sprintf("%v", r + 1))
			drawText(scr, sideBarLen, r, style, l)
		}
		scr.ShowCursor(col+sideBarLen, row)
		scr.Show()
		ev := scr.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			scr.Sync()
		case *tcell.EventKey:
			if mode == Normal {
				if ev.Key() == tcell.KeyCtrlC {
					quit(scr)
				}
				if ev.Rune() == 'h' {
					col -= 1
				}
				if ev.Rune() == 'j' {
					if buf.Len() > row+1 {
						row += 1
						if buf.Len() <= row {
							buf.AppendLine("")
						}
						col = min(buf.LineLen(row), col)
					}
				}
				if ev.Rune() == 'k' {
					row--
					if row >= 0 {
						col = min(buf.LineLen(row), col)
					}
				}
				if ev.Rune() == 'l' {
					col++
					if col > buf.LineLen(row) {
						col = buf.LineLen(row)
					}
				}
				if ev.Rune() == 'g' {
					row = 0
					col = buf.LineLen(row)
				}
				if ev.Rune() == 'G' {
					row = buf.Len() - 1
					col = buf.LineLen(row)
				}
				if ev.Rune() == 'i' {
					mode = Insert
				}
				if ev.Rune() == 'a' {
					mode = Insert
					col++
					if col > buf.LineLen(row) {
						buf.AppendToLine(row, " ")
					}
				}
				if ev.Rune() == 'I' {
					mode = Insert
					col = 0
				}
				if ev.Rune() == 'A' {
					mode = Insert
					col = buf.LineLen(row)
				}
				if ev.Rune() == 'o' {
					mode = Insert
					col = 0
					row++
					buf.InsertLine("", row)
				}
				if ev.Rune() == 'O' {
					mode = Insert
					col = 0
					buf.InsertLine("", row)
				}
			} else if mode == Insert {
				if ev.Key() == tcell.KeyEscape {
					mode = Normal
				} else if ev.Key() == tcell.KeyEnter {
					row++
					col = 0
					buf.InsertLine("", row)
				} else if ev.Key() == tcell.KeyBackspace2 {
					if buf.LineLen(row) > 0 {
						buf.ReplaceLine(row, buf.Line(row)[:buf.LineLen(row)-1])
					}
					col -= 1
					if col < 0 && row > 0 {
						if row > 0 {
							buf.RemoveLine(row)
						}
						row--
						col = buf.LineLen(row)
					}
				} else if ev.Key() == tcell.KeyCtrlW {
					length := buf.LineLen(row)
					if length == 0 {
						if row > 0 {
							buf.RemoveLine(row)
						}
						row--
						if row >= 0 {
							col = buf.LineLen(row)
						}
					} else {
						foundChar := false
						for i := length - 1; i >= 0; i-- {
							if buf.Line(row)[i] != ' ' && !foundChar {
								foundChar = true
							}
							if buf.Line(row)[i] == ' ' && foundChar {
								buf.ReplaceLine(row, buf.Line(row)[:i+1])
								col = buf.LineLen(row)
								break
							}
							if i == 0 {
								buf.ReplaceLine(row, "")
								col = 0
							}
						}
					}
				} else if ev.Key() == tcell.KeyTab {
					buf.AppendToLine(row, "    ")
					col += 4
				} else {
					if buf.LineLen(row) == 0 {
						buf.AppendToLine(row, string(ev.Rune()))
					} else {
						buf.ReplaceLine(row, buf.Line(row)[:col]+string(ev.Rune())+buf.Line(row)[col:])
					}
					col += 1
				}
			}
			if col < 0 {
				col = 0
			}
			if col >= cols {
				col = cols - 1
			}
			if row < 0 {
				row = 0
			}
			if row >= rows {
				row = rows - 1
			}
		}
	}
}

func drawText(scr tcell.Screen, col int, row int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		scr.SetContent(col, row, r, nil, style)
		col++
	}
}

func drawTextWrapping(scr tcell.Screen, startCol int, startRow int, endCol int, endRow int, style tcell.Style, text string) {
	col := startCol
	row := startRow
	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			col = startCol
		} else {
			scr.SetContent(col, row, r, nil, style)
			col++
		}
		if col >= endCol {
			row++
			col = startCol
		}
		if row > endRow {
			break
		}
	}
}

func quit(scr tcell.Screen) {
	scr.Fini()
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
