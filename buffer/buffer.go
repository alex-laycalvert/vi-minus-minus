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
	return len(buffer.lines[index])
}

func (buffer *Buffer) Line(index int) string {
	return buffer.lines[index]
}

func (buffer *Buffer) Lines() []string {
	return buffer.lines
}

func (buffer *Buffer) InsertLine(line string, index int) {
	buffer.lines = append(buffer.lines, "")
	copy(buffer.lines[index:], buffer.lines[index-1:])
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
