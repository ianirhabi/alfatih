package mailer

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"time"

	"git.qasico.com/cuxs/common/log"
)

// Message represents an email.
type Message struct {
	header      header
	parts       []*part
	attachments []*file
	embedded    []*file
	charset     string
	encoding    Encoding
	hEncoder    mimeEncoder
	buf         bytes.Buffer
}

// NewMessage creates a new message. It uses UTF-8 and quoted-printable encoding
// by default.
func NewMessage(settings ...MessageSetting) *Message {
	m := &Message{
		header:   make(header),
		charset:  "UTF-8",
		encoding: QuotedPrintable,
	}

	m.applySettings(settings)

	if m.encoding == Base64 {
		m.hEncoder = bEncoding
	} else {
		m.hEncoder = qEncoding
	}

	// Set From data Header from env variable
	m.SetAddressHeader("From", Config.SenderEmail, Config.SenderName)

	return m
}

// SetRecipient sets an list of recipient of this messages,
// it can be set multiple recipient, if you need to set email and name of recipient
// use FormatAddress instead of normal string.
func (m *Message) SetRecipient(address ...string) {
	m.encodeHeader(address)
	m.header["To"] = address
}

// SetSubject sets an value of subject email messages.
func (m *Message) SetSubject(subject ...string) {
	m.encodeHeader(subject)
	m.header["Subject"] = subject
}

// SetHeader sets a value to the given header field.
func (m *Message) SetHeader(field string, value ...string) {
	m.encodeHeader(value)
	m.header[field] = value
}

// SetHeaders sets the message headers.
func (m *Message) SetHeaders(h map[string][]string) {
	for k, v := range h {
		m.SetHeader(k, v...)
	}
}

// SetAddressHeader sets an address to the given header field.
func (m *Message) SetAddressHeader(field, address, name string) {
	m.header[field] = []string{m.FormatAddress(address, name)}
}

// SetDateHeader sets a date to the given header field.
func (m *Message) SetDateHeader(field string, date time.Time) {
	m.header[field] = []string{m.FormatDate(date)}
}

// SetBody sets the body of the message. It replaces any content previously set
// by SetBody, AddAlternative or AddAlternativeWriter.
func (m *Message) SetBody(contentType, body string, settings ...PartSetting) {
	m.parts = []*part{m.newPart(contentType, newCopier(body), settings)}
}

// GetHeader gets a header field.
func (m *Message) GetHeader(field string) []string {
	return m.header[field]
}

// FormatAddress formats an address and a name as a valid RFC 5322 address.
func (m *Message) FormatAddress(address, name string) string {
	if name == "" {
		return address
	}

	enc := m.encodeString(name)
	if enc == name {
		m.buf.WriteByte('"')
		for i := 0; i < len(name); i++ {
			b := name[i]
			if b == '\\' || b == '"' {
				m.buf.WriteByte('\\')
			}
			m.buf.WriteByte(b)
		}
		m.buf.WriteByte('"')
	} else if hasSpecials(name) {
		m.buf.WriteString(bEncoding.Encode(m.charset, name))
	} else {
		m.buf.WriteString(enc)
	}
	m.buf.WriteString(" <")
	m.buf.WriteString(address)
	m.buf.WriteByte('>')

	addr := m.buf.String()
	m.buf.Reset()
	return addr
}

// FormatDate formats a date as a valid RFC 5322 date.
func (m *Message) FormatDate(date time.Time) string {
	return date.Format(time.RFC1123Z)
}

// FormatHTML formats an html template with data interface that will be used as body
func (m *Message) FormatHTML(t *template.Template, data interface{}) string {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		panic("mailer: Error when compiling template, " + err.Error())
	}

	return buf.String()
}

// AddAlternative adds an alternative part to the message.
//
// It is commonly used to send HTML emails that default to the plain text
// version for backward compatibility. AddAlternative appends the new part to
// the end of the message. So the plain text part should be added before the
// HTML part. See http://en.wikipedia.org/wiki/MIME#Alternative
func (m *Message) AddAlternative(contentType, body string, settings ...PartSetting) {
	m.AddAlternativeWriter(contentType, newCopier(body), settings...)
}

// AddAlternativeWriter adds an alternative part to the message. It can be
// useful with the text/template or html/template packages.
func (m *Message) AddAlternativeWriter(contentType string, f func(io.Writer) error, settings ...PartSetting) {
	m.parts = append(m.parts, m.newPart(contentType, f, settings))
}

// Attach attaches the files to the email.
func (m *Message) Attach(filename string, settings ...FileSetting) {
	m.attachments = m.appendFile(m.attachments, filename, settings)
}

// Embed embeds the images to the email.
func (m *Message) Embed(filename string, settings ...FileSetting) {
	m.embedded = m.appendFile(m.embedded, filename, settings)
}

// Reset resets the message so it can be reused. The message keeps its previous
// settings so it is in the same state that after a call to NewMessage.
func (m *Message) Reset() {
	for k := range m.header {
		delete(m.header, k)
	}
	m.parts = nil
	m.attachments = nil
	m.embedded = nil
}

// Send initialing new dialer with the messages and sending the email.
func (m *Message) Send() (err error) {
	d := NewDialer()
	if err = d.DialAndSend(m); err != nil {
		log.Error(err)
	}
	return
}

func (m *Message) applySettings(settings []MessageSetting) {
	for _, s := range settings {
		s(m)
	}
}

func (m *Message) encodeHeader(values []string) {
	for i := range values {
		values[i] = m.encodeString(values[i])
	}
}

func (m *Message) encodeString(value string) string {
	return m.hEncoder.Encode(m.charset, value)
}

func (m *Message) newPart(contentType string, f func(io.Writer) error, settings []PartSetting) *part {
	p := &part{
		contentType: contentType,
		copier:      f,
		encoding:    m.encoding,
	}

	for _, s := range settings {
		s(p)
	}

	return p
}

func (m *Message) appendFile(list []*file, name string, settings []FileSetting) []*file {
	f := &file{
		Name:   filepath.Base(name),
		Header: make(map[string][]string),
		CopyFunc: func(w io.Writer) error {
			h, err := os.Open(name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, h); err != nil {
				h.Close()
				return err
			}
			return h.Close()
		},
	}

	for _, s := range settings {
		s(f)
	}

	if list == nil {
		return []*file{f}
	}

	return append(list, f)
}
