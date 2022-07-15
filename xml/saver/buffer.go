package saver

import (
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"os"
	"strings"
)

type Flag uint

type Buffer struct {
	buffer          []byte
	indentation     int
	depth           int
	size            int
	flag            Flag
	hasMultipleLine bool
}

const (
	FlagNone         Flag = 0
	FlagWriteOpenTag      = 1 << iota
	FlagWriteCloseTag
	FlagWriteIndent
	FlagWriteEmptyTag
	FlagWriteNewLine
)

func NewBuffer() *Buffer {
	b := &Buffer{
		buffer: make([]byte, 0),
		// indentation starts at -1 because the first tag is not indented.
		indentation: -1,
		depth:       0,
	}
	b.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))
	return b
}

func (b *Buffer) AddFlag(flag Flag) {
	b.flag |= flag
}

func (b *Buffer) RemoveFlag(flag Flag) {
	b.flag &^= flag
}

func (b *Buffer) Write(p []byte) {
	b.size += len(p)
	//b.hasMultipleLine = strings.Index(string(p), "\n") != -1
	b.buffer = append(b.buffer, p...)
}

func (b *Buffer) LastWriteWithNewLine() bool {
	return b.hasMultipleLine
}

func (b *Buffer) WriteWithIndent(p []byte) {
	b.Write([]byte(strings.Repeat("\t", b.indentation)))
	b.Write(p)
}

func (b *Buffer) OpenTag(tag string, attr attributes.Attributes) {
	b.IncreaseDepth()
	if attr != nil {
		b.WriteWithIndent([]byte("<" + tag + ` ` + attr.Join(" ") + ">"))
	} else {
		b.WriteWithIndent([]byte("<" + tag + ">"))
	}
}

func (b *Buffer) WriteEmptyTag(tag string, attr attributes.Attributes) {
	if attr != nil {
		b.Write([]byte("<" + tag + ` ` + attr.Join(" ") + " />"))
	} else {
		b.Write([]byte("<" + tag + " />"))
	}
}

func (b *Buffer) CloseTag(tag string) {
	//if b.buffer[b.size-1] != '\n' {
	//	b.Write([]byte("\n"))
	//}
	//b.Write([]byte(strings.Repeat("\t", b.indentation)))
	b.Write([]byte(`</` + tag + ">"))
	b.DecreaseDepth()
}

func (b *Buffer) CloseTagWithIndent(tag string) {
	b.Write([]byte(strings.Repeat("\t", b.indentation)))
	b.CloseTag(tag)
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

func (b *Buffer) Buffer() []byte {
	return b.buffer
}
