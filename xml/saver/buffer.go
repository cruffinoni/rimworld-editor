package saver

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"os"
	"strings"
)

type Buffer struct {
	buffer      []byte
	indentation int
	depth       int
}

func NewBuffer() *Buffer {
	b := &Buffer{
		buffer:      make([]byte, 0),
		indentation: 0,
		depth:       0,
	}
	b.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>\n`))
	return b
}

func (b *Buffer) Write(p []byte) {
	b.buffer = append(b.buffer, p...)
}

func (b *Buffer) OpenTag(tag string, attr xml.Attributes) {
	b.IncreaseDepth()
	b.Write([]byte(strings.Repeat("\t", b.indentation)))
	if attr != nil {
		b.Write([]byte(`<` + tag + ` ` + attr.Join(" ") + `>`))
	} else {
		b.Write([]byte(`<` + tag + `>`))
	}
}

func (b *Buffer) CloseTag(tag string) {
	b.Write([]byte(strings.Repeat("\t", b.indentation)))
	b.Write([]byte(`</` + tag + `>\n`))
	b.DecreaseDepth()
}

func (b *Buffer) IncreaseDepth() {
	b.depth++
	b.indentation++
}

func (b *Buffer) DecreaseDepth() {
	b.depth--
	b.indentation--
}

func (b *Buffer) ToFile(path string) error {
	return os.WriteFile(path, b.buffer, 0644)
}
