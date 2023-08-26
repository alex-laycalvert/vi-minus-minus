package app

import (
	"fmt"

	"github.com/alex-laycalvert/vimm/buffer"
	"github.com/gdamore/tcell"
)

type App struct {
	currentBuffer int
	buffers       []*buffer.Buffer
	screen        tcell.Screen
	cols          int
	rows          int
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
	for r := 0; r < app.rows; r++ {
		for c := 0; c < app.cols; c++ {
			app.screen.SetContent(c, r, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}
	buf := app.CurrentBuffer()
	headerStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Bold(true)
	switch buf.Mode() {
	case buffer.Normal:
		app.drawText(0, 0, headerStyle, "NORMAL")
	case buffer.Insert:
		app.drawText(0, 0, headerStyle, "INSERT")
	}
	iter := buf.Iter()
	sideBar := 1 + len(fmt.Sprintf("%v", buf.Len()))
	for iter.HasMore() {
		r, l := iter.Next()
		app.drawText(0, r+1, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue), fmt.Sprintf("%v", r+1))
		app.drawText(sideBar, r+1, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite), l)
	}
	col, row := buf.Position()
	app.screen.ShowCursor(col+sideBar, row+1)
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
		buf := app.CurrentBuffer()
		return buf.SendInput(ev, app.cols, app.rows)
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

func (app *App) drawText(col int, row int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		app.screen.SetContent(col, row, r, nil, style)
		col++
	}
}
