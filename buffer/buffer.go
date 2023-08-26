package buffer

import (
	"os"
	"strings"

	"github.com/gdamore/tcell"
)

const (
	Normal = iota
	Insert = iota
)

type Buffer struct {
	Path        string
	currentLine int
	lines       []string
	mode        int
	row         int
	col         int
}

type BufferIterator struct {
	buffer      *Buffer
	currentLine int
}

func New() *Buffer {
	return &Buffer{
		Path:  "",
		lines: []string{},
		mode:  Normal,
		row:   0,
		col:   0,
	}
}

func FromString(str string) *Buffer {
	lines := strings.Split(str, "\n")
	return &Buffer{
		Path:  "",
		lines: lines,
		mode:  Normal,
		row:   0,
		col:   0,
	}
}

func From(path string) (*Buffer, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	return &Buffer{
		Path:  path,
		lines: lines,
		mode:  Normal,
		row:   0,
		col:   0,
	}, nil
}

func (buffer *Buffer) Len() int {
	return len(buffer.lines)
}

func (buffer *Buffer) LineLen(index int) int {
	return len(buffer.lines[index])
}

func (buffer *Buffer) Mode() int {
	return buffer.mode
}

func (buffer *Buffer) Line(index int) string {
	return buffer.lines[index]
}

func (buffer *Buffer) Lines() []string {
	return buffer.lines
}

func (buffer *Buffer) InsertLine(line string, index int) {
	buffer.lines = append(buffer.lines, "")
	if index == 0 {
		copy(buffer.lines[index+1:], buffer.lines[index:])
	} else {
		copy(buffer.lines[index:], buffer.lines[index-1:])
	}
	buffer.lines[index] = line
}

func (buffer *Buffer) ReplaceLine(index int, line string) {
	buffer.lines[index] = line
}

func (buffer *Buffer) AppendLine(line string) {
	buffer.lines = append(buffer.lines, line)
}

func (buffer *Buffer) AppendToLine(index int, data string) {
	buffer.lines[index] += data
}

func (buffer *Buffer) RemoveLine(index int) {
	buffer.lines = append(buffer.lines[:index], buffer.lines[index+1:]...)
}

func (buffer *Buffer) SendInput(event *tcell.EventKey, cols int, rows int) bool {
	if buffer.mode == Normal {
		if event.Key() == tcell.KeyCtrlC {
			return true
		}
		if event.Rune() == 'h' {
			buffer.col -= 1
		}
		if event.Rune() == 'j' {
			if buffer.Len() > buffer.row+1 {
				buffer.row += 1
				if buffer.Len() <= buffer.row {
					buffer.AppendLine("")
				}
				buffer.col = min(buffer.LineLen(buffer.row), buffer.col)
			}
		}
		if event.Rune() == 'k' {
			buffer.row--
			if buffer.row >= 0 {
				buffer.col = min(buffer.LineLen(buffer.row), buffer.col)
			}
		}
		if event.Rune() == 'l' {
			buffer.col++
			if buffer.col > buffer.LineLen(buffer.row) {
				buffer.col = buffer.LineLen(buffer.row)
			}
		}
		if event.Rune() == 'g' {
			buffer.row = 0
			buffer.col = buffer.LineLen(buffer.row)
		}
		if event.Rune() == 'G' {
			buffer.row = buffer.Len() - 1
			buffer.col = buffer.LineLen(buffer.row)
		}
		if event.Rune() == 'i' {
			buffer.mode = Insert
		}
		if event.Rune() == 'a' {
			buffer.mode = Insert
			buffer.col++
			if buffer.col > buffer.LineLen(buffer.row) {
				buffer.AppendToLine(buffer.row, " ")
			}
		}
		if event.Rune() == 'I' {
			buffer.mode = Insert
			buffer.col = 0
		}
		if event.Rune() == 'A' {
			buffer.mode = Insert
			buffer.col = buffer.LineLen(buffer.row)
		}
		if event.Rune() == 'o' {
			buffer.mode = Insert
			buffer.col = 0
			buffer.row++
			buffer.InsertLine("", buffer.row)
		}
		if event.Rune() == 'O' {
			buffer.mode = Insert
			buffer.col = 0
			buffer.InsertLine("", buffer.row)
		}
	} else if buffer.mode == Insert {
		if event.Key() == tcell.KeyEscape {
			buffer.mode = Normal
		} else if event.Key() == tcell.KeyEnter {
			buffer.row++
			buffer.col = 0
			buffer.InsertLine("", buffer.row)
		} else if event.Key() == tcell.KeyBackspace2 {
			if buffer.LineLen(buffer.row) > 0 {
				buffer.ReplaceLine(buffer.row, buffer.Line(buffer.row)[:buffer.LineLen(buffer.row)-1])
			}
			buffer.col -= 1
			if buffer.col < 0 && buffer.row > 0 {
				if buffer.row > 0 {
					buffer.RemoveLine(buffer.row)
				}
				buffer.row--
				buffer.col = buffer.LineLen(buffer.row)
			}
		} else if event.Key() == tcell.KeyCtrlW {
			length := buffer.LineLen(buffer.row)
			if length == 0 {
				if buffer.row > 0 {
					buffer.RemoveLine(buffer.row)
				}
				buffer.row--
				if buffer.row >= 0 {
					buffer.col = buffer.LineLen(buffer.row)
				}
			} else {
				foundChar := false
				for i := length - 1; i >= 0; i-- {
					if buffer.Line(buffer.row)[i] != ' ' && !foundChar {
						foundChar = true
					}
					if buffer.Line(buffer.row)[i] == ' ' && foundChar {
						buffer.ReplaceLine(buffer.row, buffer.Line(buffer.row)[:i+1])
						buffer.col = buffer.LineLen(buffer.row)
						break
					}
					if i == 0 {
						buffer.ReplaceLine(buffer.row, "")
						buffer.col = 0
					}
				}
			}
		} else if event.Key() == tcell.KeyTab {
			buffer.AppendToLine(buffer.row, "    ")
			buffer.col += 4
		} else {
			if buffer.LineLen(buffer.row) == 0 {
				buffer.AppendToLine(buffer.row, string(event.Rune()))
			} else {
				buffer.ReplaceLine(buffer.row, buffer.Line(buffer.row)[:buffer.col]+string(event.Rune())+buffer.Line(buffer.row)[buffer.col:])
			}
			buffer.col += 1
		}
	}
	if buffer.col < 0 {
		buffer.col = 0
	}
	if buffer.col >= cols {
		buffer.col = cols - 1
	}
	if buffer.row < 0 {
		buffer.row = 0
	}
	if buffer.row >= rows {
		buffer.row = rows - 1
	}
	return false
}

func (buffer *Buffer) Position() (int, int) {
	return buffer.col, buffer.row
}

func (buffer *Buffer) Iter() BufferIterator {
	return BufferIterator{
		currentLine: 0,
		buffer:      buffer,
	}
}

func (bufferIter *BufferIterator) HasMore() bool {
	return bufferIter.currentLine < bufferIter.buffer.Len()
}

func (bufferIter *BufferIterator) Next() (int, string) {
	defer func() { bufferIter.currentLine++ }()
	return bufferIter.currentLine, bufferIter.buffer.Line(bufferIter.currentLine)
}
