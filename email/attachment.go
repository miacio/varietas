package email

import (
	"fmt"
	"net/textproto"
)

// Attachment
// is email attachement param type
// Based on the mime/multipart.FileHeader struct
// Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment struct {
	FileName    string               // fileName
	ContentType string               // content type
	Header      textproto.MIMEHeader // header
	Content     []byte               // content
	HTMLRelated bool
}

func (at *Attachment) setDefaultHeaders() {
	contentType := "application/octet-stream"
	if len(at.ContentType) > 0 {
		contentType = at.ContentType
	}
	at.Header.Set("Content-Type", contentType)

	if len(at.Header.Get("Content-Disposition")) == 0 {
		disposition := "attachment"
		if at.HTMLRelated {
			disposition = "inline"
		}
		at.Header.Set("Content-Disposition", fmt.Sprintf("%s;\r\n filename=\"%s\"", disposition, at.FileName))
	}
	if len(at.Header.Get("Content-ID")) == 0 {
		at.Header.Set("Content-ID", fmt.Sprintf("<%s>", at.FileName))
	}
	if len(at.Header.Get("Content-Transfer-Encoding")) == 0 {
		at.Header.Set("Content-Transfer-Encoding", "base64")
	}
}
