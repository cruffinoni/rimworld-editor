package generator

import (
	"sort"
	"strings"
)

type buffer struct {
	writtenHeaders map[string]bool
	header         string
	imp            []string
	body           string
	footer         string
}

const (
	customTypesPath = "xml/types"
	headerEmbedded  = "xml/types/primary"
	saverPath       = "xml/saver"
	xmlAttributes   = "xml/attributes"
	headerXml       = "xml"
)

func (b *buffer) writeImport(imp ...string) {
	for _, v := range imp {
		if h, ok := b.writtenHeaders[v]; ok && h {
			return
		}
		b.writtenHeaders[v] = true
		b.imp = append(b.imp, `"github.com/cruffinoni/rimworld-editor/`+v+`"`+"\n")
	}
}

func (b *buffer) writeToHeader(s string) {
	b.header += s
}

func (b *buffer) writeToBody(s string) {
	b.body += s
}

func (b *buffer) writeToFooter(s string) {
	b.footer += s
}

func (b *buffer) bytes() []byte {
	builder := strings.Builder{}
	builder.WriteString(b.header)
	if len(b.imp) > 1 {
		sort.Sort(sort.StringSlice(b.imp))
		builder.WriteString("\nimport (\n")
		for _, v := range b.imp {
			builder.WriteString(v)
		}
		builder.WriteString("\n)\n")
	} else {
		builder.WriteString("\nimport " + b.imp[0] + "\n")
	}
	builder.WriteString(b.body)
	builder.WriteString(b.footer)
	return []byte(builder.String())
}
