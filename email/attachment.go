package email

import "net/textproto"

// Attachment
// is email attachement param type
// Based on the mime/multipart.FileHeader struct
// Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment struct {
	FileName string               // fileName
	Header   textproto.MIMEHeader // header
	Content  []byte               // content
}
