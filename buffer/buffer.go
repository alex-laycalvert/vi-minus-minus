package buffer

import (
	"os"
	"strings"
)

type Buffer struct {
	Path        string
	currentLine int
	lines       []string
}

type BufferIterator struct {
	buffer      *Buffer
	currentLine int
}

func New() *Buffer {
	return &Buffer{
		Path:  "",
		lines: []string{},
	}
}

func FromString(str string) *Buffer {
	lines := strings.Split(str, "\n")
	return &Buffer{
		Path:  "",
		lines: lines,
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
	}, nil
}

func (buffer *Buffer) Len() int {
	return len(buffer.lines)
}

func (buffer *Buffer) LineLen(index int) int {
	if index >= buffer.Len() {
		return 0
	}
	return len(buffer.lines[index])
}

func (buffer *Buffer) Line(index int) string {
	if index > buffer.Len() {
		return ""
	}
	return buffer.lines[index]
}

func (buffer *Buffer) Lines() []string {
	return buffer.lines
}

func (buffer *Buffer) InsertLine(index int, line string) {
	buffer.lines = append(buffer.lines, "")
	if index == 0 {
		copy(buffer.lines[index+1:], buffer.lines[index:])
	} else {
		copy(buffer.lines[index:], buffer.lines[index-1:])
	}
	buffer.lines[index] = line
}

func (buffer *Buffer) InsertInLine(lineIndex int, index int, data string) {
	if len(buffer.lines[lineIndex]) == 0 {
		buffer.lines[lineIndex] = data
		return
	}
	buffer.lines[lineIndex] = buffer.lines[lineIndex][:index] + data + buffer.lines[lineIndex][index:]
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

func (buffer *Buffer) RemoveLine(index int) string {
	if index >= buffer.Len() {
		return ""
	}
	line := buffer.lines[index]
	buffer.lines = append(buffer.lines[:index], buffer.lines[index+1:]...)
	return line
}

func (buffer *Buffer) RemoveFromLine(lineIndex int, index int, count int) int {
	if index >= len(buffer.lines[lineIndex]) || index < 0 {
		return 0
	}
	buffer.lines[lineIndex] = buffer.lines[lineIndex][:index] + buffer.lines[lineIndex][index+count:]
	return count
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
