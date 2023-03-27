package email

import (
	"errors"
	"net/textproto"
)

const (
	MaxLineLength = 76                             // MaxLineLength is the maximum line length pre RFC 2045
	ContentType   = "text/plain; charset=us-ascii" // email.ContentType is email default Content-Type according to RFC 2045, section 5.2

	StrSubject = "Subject"
	StrTo      = "To"
	StrCc      = "Cc"
	StrBcc     = "Bcc"
	StrFrom    = "From"
)

var (
	// ErrMissingBoundary is returned when there is no boundary given for multipart entity
	ErrMissingBoundary = errors.New("no boundary found for multipart entity")

	// ErrMissingContentType is returned when there is no "Content-Type" header for a MIME entity
	ErrMissingContentType = errors.New("no Content-Type found for MIME entity")
)

// Email
// used for email message type
type Email struct {
	From        string               //from
	To          []string             // to
	Bcc         []string             // bcc
	Cc          []string             // cc
	Subject     string               // subject
	Text        []byte               // text plaintext message (optional)
	Html        []byte               // html message (optional)
	Headers     textproto.MIMEHeader // headers
	Attachments []*Attachment        // attachments
	ReadReceipt []string             // read receipt
}

// part
// copyable representation of a multipart.Part
type part struct {
	header textproto.MIMEHeader
	body   []byte
}

// New
// create an Email, and returns the pointer to it.
func New() *Email {
	return &Email{
		Headers: textproto.MIMEHeader{},
	}
}

// NewEmailFromReader reads a stream of bytes from an io.Reader, r,
// and returns an email struct containing the parsed data.
// This function expects the data in RFC 5322 format.
/* func NewEmailFromReader(r io.Reader) (*Email, error) {
	em := New()
	tp := TpReader(r)
	// parse tp headers
	hdrs, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	// set subject, to, cc, bcc and from
	for h, v := range hdrs {
		switch h {
		case StrSubject:
			em.Subject = v[0]
		case StrTo:
			em.To = v
		case StrCc:
			em.Cc = v
		case StrBcc:
			em.Bcc = v
		case StrFrom:
			em.From = v[0]
		}
		delete(hdrs, h)
	}
	em.Headers = hdrs
	body := tp.R
	ps, err := parseMIMEParts(em.Headers, body)
	if err != nil {
		return em, err
	}
	for _, p := range ps {
		headerContentType := p.header.Get("Content-Type")
		if headerContentType == "" {
			return em, ErrMissingContentType
		}
		ct, _, err := mime.ParseMediaType(headerContentType)
		if err != nil {
			return em, err
		}
		switch ct {
		case "text/plain":
			em.Text = p.body
			case ""
		}
	}
}*/
