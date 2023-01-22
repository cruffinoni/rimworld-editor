package saver

import (
	"os"
	"regexp"
	"strings"

	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type Flag uint

type Buffer struct {
	buffer    []byte
	depth     int
	bufferLen int
}

func NewBuffer() *Buffer {
	b := &Buffer{
		buffer: make([]byte, 0),
		// depth starts at -1 because the first tag is not indented.
		depth: -1,
	}
	b.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))
	return b
}

func (b *Buffer) Write(p []byte) {
	b.bufferLen += len(p)
	b.buffer = append(b.buffer, p...)
}

func (b *Buffer) Len() int {
	return b.bufferLen
}

func (b *Buffer) WriteString(s string) {
	b.Write([]byte(s))
}

func (b *Buffer) WriteStringWithIndent(s string) {
	b.WriteWithIndent([]byte(s))
}

func (b *Buffer) WriteWithIndent(p []byte) {
	b.Write([]byte(strings.Repeat("\t", b.depth)))
	b.Write(p)
}

func (b *Buffer) OpenTag(tag string, attr attributes.Attributes) {
	b.writeTag(tag, attr, true)
}

func (b *Buffer) WriteEmptyTag(tag string, attr attributes.Attributes) {
	b.writeTag(tag, attr, false)
}

func (b *Buffer) writeTag(tag string, attr attributes.Attributes, open bool) {
	if tag == "" {
		return
	}
	if b.buffer[b.bufferLen-1] != '\n' {
		b.Write([]byte("\n"))
	}
	b.IncreaseDepth()
	if attr != nil && !attr.Empty() {
		if open {
			b.WriteWithIndent([]byte("<" + tag + " " + attr.Join(" ") + ">"))
		} else {
			b.WriteWithIndent([]byte("<" + tag + ` ` + attr.Join(" ") + " />"))
		}
	} else {
		if open {
			b.WriteWithIndent([]byte("<" + tag + ">"))
		} else {
			b.WriteWithIndent([]byte("<" + tag + " />"))
		}
	}
	if !open {
		b.DecreaseDepth()
	}
}

func (b *Buffer) CloseTag(tag string) {
	if tag == "" {
		return
	}
	b.Write([]byte(`</` + tag + ">"))
	b.DecreaseDepth()
}

func (b *Buffer) CloseTagWithIndent(tag string) {
	if tag == "" {
		return
	}
	b.Write([]byte(strings.Repeat("\t", b.depth)))
	b.CloseTag(tag)
}

func (b *Buffer) IncreaseDepth() {
	b.depth++
}

func (b *Buffer) DecreaseDepth() {
	b.depth--
}

func (b *Buffer) GetDepth() int {
	return b.depth
}

func (b *Buffer) ToFile(path string) error {
	return os.WriteFile(path, b.buffer, 0644)
}

func (b *Buffer) Buffer() []byte {
	return b.buffer
}

var reMultipleLineBreak = regexp.MustCompile(`(?m)^\s*\r?\n`)

func (b *Buffer) RemoveEmptyLine() {
	b.buffer = reMultipleLineBreak.ReplaceAll(b.buffer, []byte{})
}
