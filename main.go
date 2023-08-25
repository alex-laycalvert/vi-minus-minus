package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

type Position struct {
	row int
	col int
}

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
	lines := []string{""}
	sideBarLen := 1 + len(fmt.Sprintf("%v", len(lines)-1))
	col := 0
	row := 0
	cols, rows := scr.Size()
	for {
		sideBarLen = 1 + len(fmt.Sprintf("%v", len(lines)))
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
		for r, l := range lines {
			drawText(scr, 0, r, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue), fmt.Sprintf("%v", r))
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
					if len(lines) > row+1 {
						row += 1
						if len(lines) <= row {
							lines = append(lines, "")
						}
						col = min(len(lines[row]), col)
					}
				}
				if ev.Rune() == 'k' {
					row--
					if row >= 0 {
						col = min(len(lines[row]), col)
					}
				}
				if ev.Rune() == 'l' {
					col++
					if col > len(lines[row]) {
						col = len(lines[row])
					}
				}
				if ev.Rune() == 'i' {
					mode = Insert
				}
				if ev.Rune() == 'a' {
					mode = Insert
					col++
					if col > len(lines[row]) {
						lines[row] += " "
					}
				}
				if ev.Rune() == 'I' {
					mode = Insert
					col = 0
				}
				if ev.Rune() == 'A' {
					mode = Insert
					col = len(lines[row])
				}
				if ev.Rune() == 'o' {
					mode = Insert
					col = 0
					lines = append(lines, "")
					row++
					copy(lines[row+1:], lines[row:])
					lines[row] = ""
				}
				if ev.Rune() == 'O' {
					mode = Insert
					col = 0
					lines = append(lines, "")
					copy(lines[row+1:], lines[row:])
					lines[row] = ""
				}
			} else if mode == Insert {
				if ev.Key() == tcell.KeyEscape {
					mode = Normal
				} else if ev.Key() == tcell.KeyEnter {
					row++
					col = 0
					lines = append(lines, "")
					copy(lines[row:], lines[row-1:])
					lines[row] = ""
				} else if ev.Key() == tcell.KeyBackspace2 {
					if len(lines[row]) > 0 {
						lines[row] = lines[row][0 : len(lines[row])-1]
					}
					col -= 1
					if col < 0 && row > 0 {
						row--
						col = len(lines[row])
					}
				} else if ev.Key() == tcell.KeyCtrlW {
					length := len(lines[row])
					if length == 0 {
						row--
						if row >= 0 {
							col = len(lines[row])
						}
					} else {
						foundChar := false
						for i := length - 1; i >= 0; i-- {
							if lines[row][i] != ' ' && !foundChar {
								foundChar = true
							}
							if lines[row][i] == ' ' && foundChar {
								lines[row] = lines[row][0 : i+1]
								col = len(lines[row])
								break
							}
							if i == 0 {
								lines[row] = ""
								col = 0
							}
						}
					}
				} else if ev.Key() == tcell.KeyTab {
					lines[row] += "    "
					col += 4
				} else {
					if len(lines[row]) == 0 {
						lines[row] += string(ev.Rune())
					} else {
						lines[row] = lines[row][:col] + string(ev.Rune()) + lines[row][col:]
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
