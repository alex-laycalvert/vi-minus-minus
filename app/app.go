package app

import (
	"fmt"

	"github.com/alex-laycalvert/vimm/buffer"
	"github.com/gdamore/tcell"
)

const (
	Normal = iota
	Insert = iota
)

type App struct {
	currentBuffer int
	buffers       []*buffer.Buffer
	screen        tcell.Screen
	cols          int
	rows          int
	col           int
	row           int
	mode          int
	startCol      int
	startRow      int
	endCol        int
	endRow        int
}

func New() (*App, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	err = screen.Init()
	if err != nil {
		return nil, err
	}
	cols, rows := screen.Size()
	return &App{
		currentBuffer: -1,
		buffers:       []*buffer.Buffer{},
		screen:        screen,
		cols:          cols,
		rows:          rows,
		startCol:      0,
		startRow:      0,
		endCol:        cols - 1,
		endRow:        rows - 1,
	}, nil
}

func (app *App) AddBuffer(buffer *buffer.Buffer) {
	app.buffers = append(app.buffers, buffer)
	if app.currentBuffer < 0 {
		app.currentBuffer = 0
	}
}

func (app *App) Show() {
	if len(app.buffers) == 0 {
		return
	}
	app.drawBackground()
	buf := app.CurrentBuffer()
	headerStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Bold(true)
	switch app.mode {
	case Normal:
		app.drawText(0, 0, headerStyle, "NORMAL")
	case Insert:
		app.drawText(0, 0, headerStyle, "INSERT")
	}

	startRow := app.startRow
	endRow := min(app.endRow, buf.Len()-1)
	sideBar := 1 + len(fmt.Sprintf("%v", buf.Len()))
	for r := startRow; r <= endRow; r++ {
		line := buf.Line(r)
		if app.startCol < len(line) {
			line = line[app.startCol:]
		}
		app.drawText(0, r+1, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue), fmt.Sprintf("%v", r+1))
		app.drawText(sideBar, r+1, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite), line)
	}
	app.screen.ShowCursor(app.col+sideBar, app.row+1)
	app.screen.Show()
}

func (app *App) ProcessEvent() bool {
	if len(app.buffers) == 0 {
		return true
	}
	ev := app.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		app.Resize()
	case *tcell.EventKey:
		return app.processKeyEvent(ev)
	}
	return false
}

func (app *App) CurrentBuffer() *buffer.Buffer {
	return app.buffers[app.currentBuffer]
}

func (app *App) Resize() {
	app.screen.Sync()
	cols, rows := app.screen.Size()
	app.cols = cols
	app.rows = rows
}

func (app *App) End() {
	app.screen.Fini()
}

func (app *App) drawBackground() {
	for r := 0; r < app.rows; r++ {
		for c := 0; c < app.cols; c++ {
			app.screen.SetContent(c, r, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}
}

func (app *App) drawText(col int, row int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		app.screen.SetContent(col, row, r, nil, style)
		col++
		if col >= app.cols {
			break
		}
	}
}

func (app *App) processKeyEvent(event *tcell.EventKey) bool {
	buffer := app.CurrentBuffer()
	if app.mode == Normal {
		if event.Key() == tcell.KeyCtrlC {
			return true
		}
		if event.Rune() == 'h' {
			app.col -= 1
		}
		if event.Rune() == 'j' {
			if buffer.Len() > app.row+1 {
				app.row += 1
				if buffer.Len() <= app.row {
					buffer.AppendLine("")
				}
				app.col = min(buffer.LineLen(app.row), app.col)
			}
		}
		if event.Rune() == 'k' {
			app.row--
			if app.row >= 0 {
				app.col = min(buffer.LineLen(app.row), app.col)
			}
		}
		if event.Rune() == 'l' {
			app.col++
			if app.col > buffer.LineLen(app.row) {
				app.col = buffer.LineLen(app.row) - 1
				app.startCol++
			}
		}
		if event.Rune() == 'g' {
			app.row = 0
			app.col = buffer.LineLen(app.row)
		}
		if event.Rune() == 'G' {
			app.row = buffer.Len() - 1
			app.col = buffer.LineLen(app.row)
		}
		if event.Rune() == 'i' {
			app.mode = Insert
		}
		if event.Rune() == 'a' {
			app.mode = Insert
			app.col++
			if app.col > buffer.LineLen(app.row) {
				buffer.AppendToLine(app.row, " ")
			}
		}
		if event.Rune() == 'I' {
			app.mode = Insert
			app.col = 0
			app.startCol = 0
		}
		if event.Rune() == 'A' {
			app.mode = Insert
			app.col = buffer.LineLen(app.row)
		}
		if event.Rune() == 'o' {
			app.mode = Insert
			app.col = 0
			app.startCol = 0
			app.row++
			buffer.InsertLine("", app.row)
		}
		if event.Rune() == 'O' {
			app.mode = Insert
			app.col = 0
			app.startCol = 0
			buffer.InsertLine("", app.row)
		}
	} else if app.mode == Insert {
		if event.Key() == tcell.KeyEscape {
			app.mode = Normal
		} else if event.Key() == tcell.KeyEnter {
			app.row++
			app.col = 0
			app.startCol = 0
			buffer.InsertLine("", app.row)
		} else if event.Key() == tcell.KeyBackspace2 {
			if buffer.LineLen(app.row) > 0 {
				buffer.ReplaceLine(app.row, buffer.Line(app.row)[:buffer.LineLen(app.row)-1])
			}
			app.col -= 1
			if app.col < 0 && app.row > 0 {
				if app.row > 0 {
					buffer.RemoveLine(app.row)
				}
				app.row--
				app.col = buffer.LineLen(app.row)
			}
		} else if event.Key() == tcell.KeyCtrlW {
			length := buffer.LineLen(app.row)
			if length == 0 {
				if app.row > 0 {
					buffer.RemoveLine(app.row)
				}
				app.row--
				if app.row >= 0 {
					app.col = buffer.LineLen(app.row)
				}
			} else {
				foundChar := false
				for i := length - 1; i >= 0; i-- {
					if buffer.Line(app.row)[i] != ' ' && !foundChar {
						foundChar = true
					}
					if buffer.Line(app.row)[i] == ' ' && foundChar {
						buffer.ReplaceLine(app.row, buffer.Line(app.row)[:i+1])
						app.col = buffer.LineLen(app.row)
						break
					}
					if i == 0 {
						buffer.ReplaceLine(app.row, "")
						app.col = 0
						app.startCol = 0
					}
				}
			}
		} else if event.Key() == tcell.KeyTab {
			buffer.AppendToLine(app.row, "    ")
			app.col += 4
		} else {
			if buffer.LineLen(app.row) == 0 {
				buffer.AppendToLine(app.row, string(event.Rune()))
			} else {
				buffer.ReplaceLine(app.row, buffer.Line(app.row)[:app.col]+string(event.Rune())+buffer.Line(app.row)[app.col:])
			}
			app.col += 1
		}
	}
	if app.col < 0 {
		app.col = 0
		app.startCol--
	}
	if app.col >= app.cols {
		app.col = app.cols - 1
		app.startCol++
	}
	if app.row < 0 {
		app.row = 0
	}
	if app.row >= app.rows {
		app.row = app.rows - 1
	}
	return false
}
