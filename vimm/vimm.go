package vimm

import (
	"fmt"

	"github.com/alex-laycalvert/vimm/buffer"
	"github.com/gdamore/tcell"
)

type Vimm struct {
	currentBuffer int
	buffers       []*buffer.Buffer
	screen        tcell.Screen
	cols          int
	rows          int
}

func New() (*Vimm, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	err = screen.Init()
	if err != nil {
		return nil, err
	}
	cols, rows := screen.Size()
	return &Vimm{
		currentBuffer: -1,
		buffers:       []*buffer.Buffer{},
		screen:        screen,
		cols:          cols,
		rows:          rows,
	}, nil
}

func (v *Vimm) AddBuffer(buffer *buffer.Buffer) {
	v.buffers = append(v.buffers, buffer)
	if v.currentBuffer < 0 {
		v.currentBuffer = 0
	}
}

func (v *Vimm) Show() {
	if len(v.buffers) == 0 {
		return
	}
	for r := 0; r < v.rows; r++ {
		for c := 0; c < v.cols; c++ {
			v.screen.SetContent(c, r, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}
	buf := v.CurrentBuffer()
	headerStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Bold(true)
	switch buf.Mode() {
	case buffer.Normal:
		v.drawText(0, 0, headerStyle, "NORMAL")
	case buffer.Insert:
		v.drawText(0, 0, headerStyle, "INSERT")
	}
	iter := buf.Iter()
	sideBar := 1 + len(fmt.Sprintf("%v", buf.Len()))
	for iter.HasMore() {
		r, l := iter.Next()
		v.drawText(0, r+1, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue), fmt.Sprintf("%v", r+1))
		v.drawText(sideBar, r+1, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite), l)
	}
	col, row := buf.Position()
	v.screen.ShowCursor(col+sideBar, row+1)
	v.screen.Show()
}

func (v *Vimm) ProcessEvent() bool {
	if len(v.buffers) == 0 {
		return true
	}
	ev := v.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		v.Resize()
	case *tcell.EventKey:
		buf := v.CurrentBuffer()
		return buf.SendInput(ev, v.cols, v.rows)
	}
	return false
}

func (v *Vimm) CurrentBuffer() *buffer.Buffer {
	return v.buffers[v.currentBuffer]
}

func (v *Vimm) Resize() {
	v.screen.Sync()
	cols, rows := v.screen.Size()
	v.cols = cols
	v.rows = rows
}

func (v *Vimm) End() {
	v.screen.Fini()
}

func (v *Vimm) drawText(col int, row int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		v.screen.SetContent(col, row, r, nil, style)
		col++
	}
}
