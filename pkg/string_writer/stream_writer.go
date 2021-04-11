package string_writer

import (
	"fmt"
	"io"
	"strings"
)

type ByteCounter struct {
	written int
}

func (c *ByteCounter) Record(bytesWritten int, err error) error {
	if err != nil {
		return err
	}
	c.written += bytesWritten
	return nil
}

func (c ByteCounter) Count() int {
	return c.written
}

type Type struct {
	Counter      ByteCounter
	stream       io.Writer
	indent       string
	indentLevel  uint
	cachedIndent string
}

func New(stream io.Writer, indent string) *Type {
	return &Type{
		stream: stream,
		indent: indent,
	}
}

func (w *Type) Write(v string) error {
	return w.Counter.Record(w.stream.Write([]byte(w.cachedIndent + v)))
}

func (w *Type) WriteLn(v string) error {
	return w.Write(v + "\n")
}

func (w *Type) Write2Ln(v string) error {
	return w.Write(v + "\n\n")
}

func (w *Type) WriteLnF(format string, values ...interface{}) error {
	return w.WriteLn(fmt.Sprintf(format, values...))
}

func (w *Type) In(callback func(out *Type) error) error {
	w.indentLevel++
	w.cachedIndent = strings.Repeat(w.indent, int(w.indentLevel))
	err := callback(w)
	w.indentLevel--
	w.cachedIndent = strings.Repeat(w.indent, int(w.indentLevel))
	return err
}
