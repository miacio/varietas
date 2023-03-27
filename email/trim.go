package email

import (
	"bufio"
	"bytes"
	"io"
	"net/textproto"
	"unicode"
)

// trimReader
// a custom io.Reader that will trim any leading
// whitespace, as this can cause email imports to fail.
type trimReader struct {
	rd      io.Reader
	trimmed bool
}

// Read
// trims off any unicode whitespace from the originating reader
func (tr *trimReader) Read(buf []byte) (int, error) {
	n, err := tr.rd.Read(buf)
	if err != nil {
		return n, err
	}
	if !tr.trimmed {
		t := bytes.TrimLeftFunc(buf[:n], unicode.IsSpace)
		tr.trimmed = true
		n = copy(buf, t)
	}
	return n, err
}

// TpReader
// textproto reader io
func TpReader(r io.Reader) *textproto.Reader {
	s := &trimReader{rd: r}
	return textproto.NewReader(bufio.NewReader(s))
}
