// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailer

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// header type represents an request header
type header map[string][]string

// Encoding represents a MIME encoding scheme like quoted-printable or base64.
type Encoding string

const (
	// QuotedPrintable represents the quoted-printable encoding as defined in
	// RFC 2045.
	QuotedPrintable Encoding = "quoted-printable"
	// Base64 represents the base64 encoding as defined in RFC 2045.
	Base64 Encoding = "base64"
	// Unencoded can be used to avoid encoding the body of an email. The headers
	// will still be encoded using quoted-printable encoding.
	Unencoded Encoding = "8bit"
)

// A MessageSetting can be used as an argument in NewMessage to configure an
// email.
type MessageSetting func(m *Message)

// SetCharset is a message setting to set the charset of the email.
func SetCharset(charset string) MessageSetting {
	return func(m *Message) {
		m.charset = charset
	}
}

// SetEncoding is a message setting to set the encoding of the email.
func SetEncoding(enc Encoding) MessageSetting {
	return func(m *Message) {
		m.encoding = enc
	}
}

// ParseTemplate perform template parsing from path into template html
func ParseTemplate(filename string) *template.Template {
	// fuck this, we load an template from directory
	// that was defined on environment variable
	// because if we dont set the full path
	// the test will be failed cause tmp directory
	// that created by test
	tf := filepath.Join(os.Getenv("EMAIL_TEMPLATE_DIR"), filename)

	t, err := template.ParseFiles(tf)
	if err != nil {
		panic("mailer: Error when parsing template, " + err.Error())
	}

	return t
}

func hasSpecials(text string) bool {
	for i := 0; i < len(text); i++ {
		switch c := text[i]; c {
		case '(', ')', '<', '>', '[', ']', ':', ';', '@', '\\', ',', '.', '"':
			return true
		}
	}

	return false
}

func newCopier(s string) func(io.Writer) error {
	return func(w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	}
}
