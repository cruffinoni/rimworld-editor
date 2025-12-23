package writer

import (
	"sort"
	"strings"

	"github.com/cruffinoni/rimworld-editor/internal/codegen/importpath"
)

type buffer struct {
	writtenHeaders map[string]bool
	header         strings.Builder
	imp            []string
	body           strings.Builder
	footer         strings.Builder
}

func (b *buffer) writeImport(imp ...string) {
	for _, v := range imp {
		if h, ok := b.writtenHeaders[v]; ok && h {
			return
		}
		b.writtenHeaders[v] = true
		b.imp = append(b.imp, `"`+importpath.CodePackage+v+`"`+"\n")
	}
}

func (b *buffer) writeToHeader(s string) {
	b.header.WriteString(s)
}

func (b *buffer) writeToBody(s string) {
	b.body.WriteString(s)
}

func (b *buffer) writeToFooter(s string) {
	b.footer.WriteString(s)
}

func (b *buffer) bytes() []byte {
	builder := strings.Builder{}
	builder.WriteString(b.header.String())
	if len(b.imp) > 1 {
		sort.Strings(b.imp)
		builder.WriteString("\nimport (\n")
		for _, v := range b.imp {
			builder.WriteString(v)
		}
		builder.WriteString("\n)\n")
	} else if len(b.imp) == 1 {
		builder.WriteString("\nimport " + b.imp[0] + "\n")
	}
	builder.WriteString(b.body.String())
	builder.WriteString(b.footer.String())
	return []byte(builder.String())
}
