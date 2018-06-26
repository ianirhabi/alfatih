// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailer

import "io"

type file struct {
	Name     string
	Header   map[string][]string
	CopyFunc func(w io.Writer) error
}

func (f *file) setHeader(field, value string) {
	f.Header[field] = []string{value}
}

// A FileSetting can be used as an argument in Message.Attach or Message.Embed.
type FileSetting func(*file)

// SetHeader is a file setting to set the MIME header of the message part that
// contains the file content.
//
// Mandatory headers are automatically added if they are not set when sending
// the email.
func SetHeader(h map[string][]string) FileSetting {
	return func(f *file) {
		for k, v := range h {
			f.Header[k] = v
		}
	}
}

// Rename is a file setting to set the name of the attachment if the name is
// different than the filename on disk.
func Rename(name string) FileSetting {
	return func(f *file) {
		f.Name = name
	}
}

// SetCopyFunc is a file setting to replace the function that runs when the
// message is sent. It should copy the content of the file to the io.Writer.
//
// The default copy function opens the file with the given filename, and copy
// its content to the io.Writer.
func SetCopyFunc(f func(io.Writer) error) FileSetting {
	return func(fi *file) {
		fi.CopyFunc = f
	}
}

type part struct {
	contentType string
	copier      func(io.Writer) error
	encoding    Encoding
}

// A PartSetting can be used as an argument in Message.SetBody,
// Message.AddAlternative or Message.AddAlternativeWriter to configure the part
// added to a message.
type PartSetting func(*part)

// SetPartEncoding sets the encoding of the part added to the message. By
// default, parts use the same encoding than the message.
func SetPartEncoding(e Encoding) PartSetting {
	return PartSetting(func(p *part) {
		p.encoding = e
	})
}
