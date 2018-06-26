// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/mgutz/ansi"
	"golang.org/x/crypto/ssh/terminal"
)

var timestampFormat = "2006/01/02 - 15:04:05"

type LogFormatter struct {
	// Using colors or not.
	Tty bool

	// Prefix formater
	Prefix string

	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Whether the logger's out is to a terminal
	isTerminal   bool
	terminalOnce sync.Once
}

// Format implement logrus.Formatter
// The format will be determined by f.Tty & f.ForceColors
// If tty is false then jsonFormater will be used.
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.terminalOnce.Do(func() {
		if entry.Logger != nil {
			f.isTerminal = f.checkIfTerminal(entry.Logger.Out)
		}
	})

	isColorTerminal := f.isTerminal && (runtime.GOOS != "windows")
	isColored := (f.ForceColors || isColorTerminal) && f.Tty

	if isColored {
		return f.textFormater(entry)
	}

	return f.jsonFormater(entry)
}

// textFormater formating text with colored and custom format.
func (f *LogFormatter) textFormater(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		if k != "prefix" {
			keys = append(keys, k)
		}
	}

	prefixFieldClashes(entry.Data)

	var levelColor string
	var levelText string

	prefixColor := ansi.ColorCode("white+b:green")
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = ansi.Green
	case logrus.WarnLevel:
		levelColor = ansi.Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		prefixColor = ansi.ColorCode("white+b:red")
		levelColor = ansi.Red
	default:
		levelColor = ansi.Blue
	}

	if entry.Level != logrus.WarnLevel {
		levelText = strings.ToUpper(entry.Level.String())
	} else {
		levelText = "WARN"
	}

	prefix := ""
	message := entry.Message

	if prefixValue, ok := entry.Data["prefix"]; ok {
		prefix = fmt.Sprintf("%s %s %s ", prefixColor, prefixValue.(string), ansi.Reset)
	} else {
		prefixValue, trimmedMsg := extractPrefix(entry.Message)
		if len(prefixValue) > 0 {
			prefix = fmt.Sprintf("%s %s %s ", prefixColor, prefixValue, ansi.Reset)
			message = trimmedMsg
		}
	}

	format := "%s[%s]%s %s%+5s%s %s%s"
	if f.Prefix != "" {
		format = "[" + f.Prefix + "]" + format
	}
	fmt.Fprintf(b, format, ansi.LightBlack, entry.Time.Format(timestampFormat), ansi.Reset, levelColor, levelText, ansi.Reset, prefix, message)
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " %s%s%s=%+v", levelColor, k, ansi.Reset, v)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// jsonFormater make log with format json.
func (f *LogFormatter) jsonFormater(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case string:
			data[k] = v
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	prefixFieldClashes(data)

	data["time"] = entry.Time.Format(timestampFormat)
	data["msg"] = entry.Message
	data["level"] = entry.Level.String()

	if f.Prefix != "" {
		data["logger"] = f.Prefix
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

// checkIfTerminal check if this caller is on terminal
func (f *LogFormatter) checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

// NewFormater return new instance of LogFormatter.
func NewFormater(debug bool, prefix ...string) *LogFormatter {
	pfx := ""
	if len(prefix) > 0 {
		pfx = prefix[0]
	}

	return &LogFormatter{
		Prefix: pfx,
		Tty:    debug,
	}
}

// extractPrefix return key of prefix with value
// from string.
func extractPrefix(msg string) (string, string) {
	prefix := ""
	regex := regexp.MustCompile("^\\[(.*?)\\]")
	if regex.MatchString(msg) {
		match := regex.FindString(msg)
		prefix, msg = match[1:len(match)-1], strings.TrimSpace(msg[len(match):])
	}
	return prefix, msg
}

// prefixFieldClashes resolve if logrus.Fields has duplicate
// key with defined keys.
func prefixFieldClashes(data logrus.Fields) {
	_, ok := data["time"]
	if ok {
		data["fields.time"] = data["time"]
	}
	_, ok = data["msg"]
	if ok {
		data["fields.msg"] = data["msg"]
	}
	_, ok = data["level"]
	if ok {
		data["fields.level"] = data["level"]
	}
	_, ok = data["logger"]
	if ok {
		data["fields.logger"] = data["logger"]
	}
	_, ok = data["prefix"]
	if ok {
		data["fields.prefix"] = data["prefix"]
	}
}
