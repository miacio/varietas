package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	MaxLineLength      = 76                             // MaxLineLength is the maximum line length pre RFC 2045
	DefaultContentType = "text/plain; charset=us-ascii" // email.ContentType is email default Content-Type according to RFC 2045, section 5.2

	StrReplyTo     = "Reply-To"
	StrSubject     = "Subject"
	StrTo          = "To"
	StrCc          = "Cc"
	StrBcc         = "Bcc"
	StrFrom        = "From"
	StrDate        = "Date"
	StrMessageId   = "Message-Id"
	StrMimeVersion = "MIME-Version"
)

var (
	// ErrMissingBoundary is returned when there is no boundary given for multipart entity
	ErrMissingBoundary = errors.New("no boundary found for multipart entity")

	// ErrMissingContentType is returned when there is no "Content-Type" header for a MIME entity
	ErrMissingContentType = errors.New("no Content-Type found for MIME entity")

	// ErrMustSpecifyMessage is returned when there is email param exist empty
	ErrMustSpecifyMessage = errors.New("must specify at least one from address and one to address")
)

// Email
// used for email message type
type Email struct {
	ReplyTo     []string             // reply to
	From        string               //from
	To          []string             // to
	Bcc         []string             // bcc
	Cc          []string             // cc
	Subject     string               // subject
	Text        []byte               // text plaintext message (optional)
	Html        []byte               // html message (optional)
	Sender      string               // override from as SMTP envelope sender (optional)
	Headers     textproto.MIMEHeader // headers
	Attachments []*Attachment        // attachments
	ReadReceipt []string             // read receipt

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
func NewEmailFromReader(r io.Reader) (*Email, error) {
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
			subj, err := (&mime.WordDecoder{}).DecodeHeader(em.Subject)
			if err == nil && len(subj) > 0 {
				em.Subject = subj
			}
		case StrTo:
			em.To = handleAddressList(v)
		case StrCc:
			em.Cc = handleAddressList(v)
		case StrBcc:
			em.Bcc = handleAddressList(v)
		case StrReplyTo:
			em.ReplyTo = handleAddressList(v)
		case StrFrom:
			em.From = v[0]
			fr, err := (&mime.WordDecoder{}).DecodeHeader(em.From)
			if err == nil && len(fr) > 0 {
				em.From = fr
			}
		}
		delete(hdrs, h)
	}
	em.Headers = hdrs
	body := tp.R
	// recursively parse the MIME parts
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

		if cd := p.header.Get("Content-Disposition"); cd != "" {
			cd, params, err := mime.ParseMediaType(p.header.Get("Content-Disposition"))
			if err != nil {
				return em, err
			}
			fileName, fileNameDefined := params["fileName"]
			if cd == "attachment" || (cd == "inline" && fileNameDefined) {
				_, err := em.Attach(bytes.NewReader(p.body), fileName, ct)
				if err != nil {
					return em, err
				}
				continue
			}
		}

		switch ct {
		case "text/plain":
			em.Text = p.body
		case "text/html":
			em.Html = p.body
		}
	}
	return em, nil
}

// Attach is used to attach content from an io.Reader to the email.
// Required params include an io.Reader, the desired fileName for the attachment, and the Content-Type
// the func will return the create Attachment for reference, as well as nil for the error, if successful.
func (e *Email) Attach(r io.Reader, fileName string, contentType string) (a *Attachment, err error) {
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, r); err != nil {
		return
	}
	at := &Attachment{
		FileName:    fileName,
		ContentType: contentType,
		Header:      textproto.MIMEHeader{},
		Content:     buffer.Bytes(),
	}
	e.Attachments = append(e.Attachments, at)
	return at, nil
}

// AttachFile is used to attach content to the email.
// it attempts to open the file referenced by fileName and, if successful, creates an attachment.
// this attachment is then appended to the slice of e.Attachments.
// the func will then return the Attachment for reference, as well as nil for the error, if successful.
func (e *Email) AttachFile(fileName string) (*Attachment, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ct := mime.TypeByExtension(filepath.Ext(fileName))
	basename := filepath.Base(fileName)
	return e.Attach(f, basename, ct)
}

// messageHeaders merges the email's various fields and custom headers together in a standards
// create a MIMEHeader to be used in the result message.
// it does not alter e.Headers.
// "e"'s fields To, Cc, From, Subject will be used unless they are present in e.Headers.
// Unless set in e.Headers, "Date" will filled with the current time.
func (e *Email) messageHeaders() (textproto.MIMEHeader, error) {
	res := make(textproto.MIMEHeader, len(e.Headers)+6)
	if e.Headers != nil {
		for _, h := range []string{StrReplyTo, StrTo, StrCc, StrFrom, StrSubject, StrDate, StrMessageId, StrMimeVersion} {
			if v, ok := e.Headers[h]; ok {
				res[h] = v
			}
		}
	}
	// Set headers if there are values.
	if _, ok := res[StrReplyTo]; !ok && len(e.ReplyTo) > 0 {
		res.Set(StrReplyTo, strings.Join(e.ReplyTo, ", "))
	}
	if _, ok := res[StrTo]; !ok && len(e.To) > 0 {
		res.Set(StrTo, strings.Join(e.To, ", "))
	}
	if _, ok := res[StrCc]; !ok && len(e.Cc) > 0 {
		res.Set(StrCc, strings.Join(e.Cc, ", "))
	}
	if _, ok := res[StrSubject]; !ok && e.Subject != "" {
		res.Set(StrSubject, e.Subject)
	}
	if _, ok := res[StrMessageId]; !ok {
		id, err := generateMessageID()
		if err != nil {
			return nil, err
		}
		res.Set(StrMessageId, id)
	}
	// Date and From are required headers.
	if _, ok := res[StrFrom]; !ok {
		res.Set(StrFrom, e.From)
	}
	if _, ok := res[StrDate]; !ok {
		res.Set(StrDate, time.Now().Format(time.RFC1123Z))
	}
	if _, ok := res[StrMimeVersion]; !ok {
		res.Set(StrMimeVersion, "1.0")
	}
	for field, vals := range e.Headers {
		if _, ok := res[field]; !ok {
			res[field] = vals
		}
	}
	return res, nil
}

func (e *Email) categorizeAttachments() (htmlRelated, others []*Attachment) {
	for _, a := range e.Attachments {
		if a.HTMLRelated {
			htmlRelated = append(htmlRelated, a)
		} else {
			others = append(others, a)
		}
	}
	return
}

// Bytes converts the Email object to []byte
// including all needed MIMEHeaders, boundaries, etc.
func (e *Email) Bytes() ([]byte, error) {
	// TODO: better guess buffer size
	buff := bytes.NewBuffer(make([]byte, 0, 4096))

	headers, err := e.messageHeaders()
	if err != nil {
		return nil, err
	}

	htmlAttachments, otherAttachments := e.categorizeAttachments()
	if len(e.Html) == 0 && len(htmlAttachments) > 0 {
		return nil, errors.New("there are HTML attachments, but no HTML body")
	}

	var (
		isMixed       = len(otherAttachments) > 0
		isAlternative = len(e.Text) > 0 && len(e.Html) > 0
		isRelated     = len(e.Html) > 0 && len(htmlAttachments) > 0
	)

	var w *multipart.Writer
	if isMixed || isAlternative || isRelated {
		w = multipart.NewWriter(buff)
	}
	switch {
	case isMixed:
		headers.Set("Content-Type", "multipart/mixed;\r\n boundary="+w.Boundary())
	case isAlternative:
		headers.Set("Content-Type", "multipart/alternative;\r\n boundary="+w.Boundary())
	case isRelated:
		headers.Set("Content-Type", "multipart/related;\r\n boundary="+w.Boundary())
	case len(e.Html) > 0:
		headers.Set("Content-Type", "text/html; charset=UTF-8")
		headers.Set("Content-Transfer-Encoding", "quoted-printable")
	default:
		headers.Set("Content-Type", "text/plain; charset=UTF-8")
		headers.Set("Content-Transfer-Encoding", "quoted-printable")
	}
	headerToBytes(buff, headers)
	_, err = io.WriteString(buff, "\r\n")
	if err != nil {
		return nil, err
	}

	// Check to see if there is a Text or HTML field
	if len(e.Text) > 0 || len(e.Html) > 0 {
		var subWriter *multipart.Writer

		if isMixed && isAlternative {
			// Create the multipart alternative part
			subWriter = multipart.NewWriter(buff)
			header := textproto.MIMEHeader{
				"Content-Type": {"multipart/alternative;\r\n boundary=" + subWriter.Boundary()},
			}
			if _, err := w.CreatePart(header); err != nil {
				return nil, err
			}
		} else {
			subWriter = w
		}
		// Create the body sections
		if len(e.Text) > 0 {
			// Write the text
			if err := writeMessage(buff, e.Text, isMixed || isAlternative, "text/plain", subWriter); err != nil {
				return nil, err
			}
		}
		if len(e.Html) > 0 {
			messageWriter := subWriter
			var relatedWriter *multipart.Writer
			if (isMixed || isAlternative) && len(htmlAttachments) > 0 {
				relatedWriter = multipart.NewWriter(buff)
				header := textproto.MIMEHeader{
					"Content-Type": {"multipart/related;\r\n boundary=" + relatedWriter.Boundary()},
				}
				if _, err := subWriter.CreatePart(header); err != nil {
					return nil, err
				}

				messageWriter = relatedWriter
			} else if isRelated && len(htmlAttachments) > 0 {
				relatedWriter = w
				messageWriter = w
			}
			// Write the HTML
			if err := writeMessage(buff, e.Html, isMixed || isAlternative || isRelated, "text/html", messageWriter); err != nil {
				return nil, err
			}
			if len(htmlAttachments) > 0 {
				for _, a := range htmlAttachments {
					a.setDefaultHeaders()
					ap, err := relatedWriter.CreatePart(a.Header)
					if err != nil {
						return nil, err
					}
					// Write the base64Wrapped content to the part
					base64Wrap(ap, a.Content)
				}

				if isMixed || isAlternative {
					relatedWriter.Close()
				}
			}
		}
		if isMixed && isAlternative {
			if err := subWriter.Close(); err != nil {
				return nil, err
			}
		}
	}
	// Create attachment part, if necessary
	for _, a := range otherAttachments {
		a.setDefaultHeaders()
		ap, err := w.CreatePart(a.Header)
		if err != nil {
			return nil, err
		}
		// Write the base64Wrapped content to the part
		base64Wrap(ap, a.Content)
	}
	if isMixed || isAlternative || isRelated {
		if err := w.Close(); err != nil {
			return nil, err
		}
	}
	return buff.Bytes(), nil
}

// Select and parse an SMTP envelope sender address.  Choose Email.Sender if set, or fallback to Email.From.
func (e *Email) parseSender() (string, error) {
	if e.Sender != "" {
		sender, err := mail.ParseAddress(e.Sender)
		if err != nil {
			return "", err
		}
		return sender.Address, nil
	} else {
		from, err := mail.ParseAddress(e.From)
		if err != nil {
			return "", err
		}
		return from.Address, nil
	}
}

// Send an email using the given host and SMTP auth (optional), returns any error thrown by smtp.SendMail
// This function merges the To, Cc, and Bcc fields and calls the smtp.SendMail function using the Email.Bytes() output as the message
func (e *Email) Send(addr string, a smtp.Auth) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(e.To)+len(e.Cc)+len(e.Bcc))
	to = append(append(append(to, e.To...), e.Cc...), e.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if e.From == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	sender, err := e.parseSender()
	if err != nil {
		return err
	}
	raw, err := e.Bytes()
	if err != nil {
		return err
	}
	return smtp.SendMail(addr, a, sender, to, raw)
}

// SendWithTLS sends an email with an optional TLS config.
// This is helpful if you need to connect to a host that is used an untrusted
// certificate.
func (e *Email) SendWithTLS(addr string, auth smtp.Auth, tlsConfig *tls.Config) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(e.To)+len(e.Cc)+len(e.Bcc))
	to = append(append(append(to, e.To...), e.Cc...), e.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if e.From == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	sender, err := e.parseSender()
	if err != nil {
		return err
	}
	raw, err := e.Bytes()
	if err != nil {
		return err
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, tlsConfig.ServerName)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("localhost"); err != nil {
		return err
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(sender); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// SendWithStartTLS sends an email over TLS using STARTTLS with an optional TLS config.
//
// The TLS Config is helpful if you need to connect to a host that is used an untrusted
// certificate.
func (e *Email) SendWithStartTLS(addr string, auth smtp.Auth, tlsConfig *tls.Config) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(e.To)+len(e.Cc)+len(e.Bcc))
	to = append(append(append(to, e.To...), e.Cc...), e.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if e.From == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	sender, err := e.parseSender()
	if err != nil {
		return err
	}
	raw, err := e.Bytes()
	if err != nil {
		return err
	}

	// Taken from the standard library
	// https://github.com/golang/go/blob/master/src/net/smtp/smtp.go#L328
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("localhost"); err != nil {
		return err
	}
	// Use TLS if available
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(tlsConfig); err != nil {
			return err
		}
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(sender); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
