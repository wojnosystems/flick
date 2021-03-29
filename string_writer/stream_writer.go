package string_writer

import (
	"fmt"
	"io"
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
	Counter ByteCounter
	stream  io.Writer
}

func New(stream io.Writer) *Type {
	return &Type{
		stream: stream,
	}
}

func (w *Type) Write(v string) error {
	return w.Counter.Record(w.stream.Write([]byte(v)))
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
