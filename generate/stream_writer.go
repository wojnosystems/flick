package generate

import (
	"fmt"
	"io"
)

type byteCounter struct {
	written int
}

func (c *byteCounter) record(bytesWritten int, err error) error {
	if err != nil {
		return err
	}
	c.written += bytesWritten
	return nil
}

func (c byteCounter) Count() int {
	return c.written
}

type stringWriter struct {
	counter byteCounter
	stream  io.Writer
}

func (w *stringWriter) write(v string) error {
	return w.counter.record(w.stream.Write([]byte(v)))
}

func (w *stringWriter) writeLn(v string) error {
	return w.write(v + "\n")
}

func (w *stringWriter) write2Ln(v string) error {
	return w.write(v + "\n\n")
}

func (w *stringWriter) writeLnF(format string, values ...interface{}) error {
	return w.writeLn(fmt.Sprintf(format, values...))
}
