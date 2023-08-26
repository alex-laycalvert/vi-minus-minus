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
	clipboard     string
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
	sideBar := len(fmt.Sprintf("%v ", buf.Len()))
	for r := app.startRow; r <= max(app.rows, buf.Len()); r++ {
		if r >= buf.Len() {
			break
		}
		line := buf.Line(r)
		if app.startCol >= len(line) && app.startCol != 0 {
			line = ""
		} else if app.startCol < len(line) {
			line = line[app.startCol:]
		}
		app.drawText(0, r+1-app.startRow, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite), fmt.Sprintf("%v", r+1))
		app.drawText(sideBar, r+1-app.startRow, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite), line)
	}
	app.screen.ShowCursor(app.col+sideBar-app.startCol, app.row+1-app.startRow)
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
		if col > app.cols {
			break
		}
	}
}

func (app *App) processKeyEvent(event *tcell.EventKey) bool {
	buf := app.CurrentBuffer()
	switch app.mode {
	case Normal:
		if event.Key() == tcell.KeyCtrlC {
			return true
		} else if event.Rune() == 'h' {
			app.col--
			app.startCol--
		} else if event.Rune() == 'j' {
			app.row++
			app.col = min(app.col, buf.LineLen(app.row))
		} else if event.Rune() == 'k' {
			app.row--
		} else if event.Rune() == 'l' {
			app.col++
		} else if event.Rune() == 'g' {
			app.row = 0
			app.startRow = 0
			app.col = min(app.col, buf.LineLen(app.row))
		} else if event.Rune() == 'G' {
			app.row = buf.Len() - 1
			app.col = min(app.col, buf.LineLen(app.row))
		} else if event.Rune() == 'd' {
			app.clipboard = buf.RemoveLine(app.row)
			if buf.Len() == 0 {
				buf.AppendLine("")
			}
		} else if event.Rune() == 'p' {
			buf.InsertInLine(app.row, app.col-1, app.clipboard)
		} else if event.Rune() == 'i' {
			app.mode = Insert
		} else if event.Rune() == 'I' {
			app.mode = Insert
			app.col = 0
			app.startCol = 0
		} else if event.Rune() == 'a' {
			app.mode = Insert
			app.col++
		} else if event.Rune() == 'A' {
			app.mode = Insert
			app.col = buf.LineLen(app.row)
		} else if event.Rune() == 'o' {
			app.mode = Insert
			app.row++
			buf.InsertLine(app.row, "")
			app.col = 0
			app.startCol = 0
		} else if event.Rune() == 'O' {
			app.mode = Insert
			buf.InsertLine(app.row, "")
			app.col = 0
			app.startCol = 0
		}
	case Insert:
		if event.Key() == tcell.KeyEscape {
			app.mode = Normal
		} else if event.Key() == tcell.KeyCtrlW {
			// TODO: kill last word
		} else if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			if app.col == 0 && buf.Len() > 1 {
				line := buf.RemoveLine(app.row)
				app.row--
				app.col = buf.LineLen(app.row)
				if app.startRow > 0 {
					app.startRow--
				}
				buf.AppendToLine(app.row, line)
			} else {
				removed := buf.RemoveFromLine(app.row, app.col-1, 1)
				app.col -= len(removed)
				app.startCol -= len(removed)
			}
		} else if event.Key() == tcell.KeyEnter {
			app.row++
			buf.InsertLine(app.row, "")
			app.col = 0
			app.startCol = 0
		} else if event.Key() == tcell.KeyTab {
			buf.InsertInLine(app.row, app.col, "    ")
			app.col += 4
		} else {
			buf.InsertInLine(app.row, app.col, string(event.Rune()))
			app.col++
		}
	}
	if app.row < app.startRow {
		app.startRow = app.row
	}
	if app.row >= buf.Len() {
		app.row = buf.Len() - 1
	}
	if app.row >= app.rows+app.startRow-1 {
		app.startRow += app.row - (app.rows + app.startRow - 2)
	}
	if app.row < 0 {
		app.row = 0
		app.startRow = 0
		app.startCol = 0
		app.col = 0
	}
	sideBar := 2 + len(fmt.Sprintf("%v", buf.Len()))
	if app.col > buf.LineLen(app.row) {
		app.col = buf.LineLen(app.row)
	}
	if app.col > app.cols-sideBar+app.startCol {
		// TODO: maybe add a -1 to this?
		app.startCol += (app.col - app.startCol) - (app.cols - sideBar)
	}
	if app.col < 0 {
		app.col = 0
		app.startCol = 0
	}
	if app.col < app.startCol {
		app.startCol = app.col
	}
	if app.startCol < 0 {
		app.startCol = 0
	}
	return false
}
