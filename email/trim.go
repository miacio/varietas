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
	rd io.Reader
}

// Read
// trims off any unicode whitespace from the originating reader
func (tr trimReader) Read(buf []byte) (int, error) {
	n, err := tr.rd.Read(buf)
	t := bytes.TrimLeftFunc(buf[:n], unicode.IsSpace)
	n = copy(buf, t)
	return n, err
}

// TpReader
// textproto reader io
func TpReader(r io.Reader) *textproto.Reader {
	s := trimReader{rd: r}
	return textproto.NewReader(bufio.NewReader(s))
}
