package printer

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Writer struct {
	out *os.File
	in  *os.File
	err *os.File
}

func NewPrint() *Writer {
	w := &Writer{
		out: os.Stdout,
		in:  os.Stdin,
		err: os.Stderr,
	}
	//w.out, _ = os.OpenFile("out.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	return w
}

func (l *Writer) formatColor(buffer []byte) []byte {
	//log.Printf("Buffer: '%s'", buffer)
	f := colorFinderRegex.FindAllSubmatch(buffer, -1)
	//log.Printf("f: %q", f)
	if f == nil {
		return buffer
	}
	for _, i := range f {
		replacement := []byte("\x1b[")
		//log.Printf("i: %q", i)
		composed := bytes.Split(i[1], []byte(","))
		//log.Printf("Composed: %q", composed)
		for _, c := range composed {
			if bytes.HasPrefix(c, []byte("B_")) {
				color := bytes.TrimPrefix(c, []byte("B_"))
				if col, ok := colorValues[strings.ToLower(string(color))]; ok {
					replacement = append(replacement, []byte(strconv.FormatInt(int64(col+BackgroundBlack), 10))...)
				} else {
					replacement = append(replacement, []byte("%B_COLOR_NOT_FOUND%")...)
					replacement = append(replacement, c...)
					replacement = append(replacement, '%')
				}
			} else if bytes.HasPrefix(c, []byte("F_")) {
				color := bytes.TrimPrefix(c, []byte("F_"))
				if col, ok := colorValues[strings.ToLower(string(color))]; ok {
					replacement = append(replacement, []byte(strconv.FormatInt(int64(col+ForegroundBlack), 10))...)
				} else {
					replacement = append(replacement, []byte("%F_COLOR_NOT_FOUND%")...)
					replacement = append(replacement, c...)
					replacement = append(replacement, '%')
				}
			} else {
				if opt, ok := colorOptions[strings.ToLower(string(c))]; ok {
					replacement = append(replacement, []byte(strconv.FormatInt(int64(opt), 10))...)
				} else {
					replacement = append(replacement, []byte("%NOT_FOUND%")...)
					replacement = append(replacement, c...)
					replacement = append(replacement, '%')
				}
			}
			replacement = append(replacement, ';')
		}
		replacement = replacement[:len(replacement)-1]
		replacement = append(replacement, 'm')
		buffer = bytes.ReplaceAll(buffer, i[0], replacement)
	}
	buffer = append(buffer, []byte("\x1b[0m")...)
	return buffer
}

func (l *Writer) WriteToError(b []byte) {
	b = append([]byte("{-F_RED,BOLD}Error:{-RESET} "), b...)
	l.write(b, l.err)
}

func (l *Writer) WriteToStd(b []byte) {
	l.write(b, l.out)
}

func (l *Writer) write(b []byte, out *os.File) {
	b = l.formatColor(b)
	_, err := out.Write(b)
	if !bytes.HasSuffix(b, []byte("\n")) {
		_, err = out.Write([]byte("\n"))
	}
	//time.Sleep(50 * time.Millisecond)
	if err != nil {
		panic(err)
	}
}

func (l *Writer) WriteToStdf(format string, a ...any) {
	b := []byte(fmt.Sprintf(format, a...))
	l.write(b, l.out)
}

func (l *Writer) WriteToErrf(format string, a ...any) {
	b := []byte(fmt.Sprintf(format, a...))
	l.WriteToError(b)
}
