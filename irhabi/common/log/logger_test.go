// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var output bytes.Buffer

func TestNew(t *testing.T) {
	var output bytes.Buffer

	l := New()
	assert.IsType(t, &logrus.Logger{}, l)

	l = New("LOGGER")
	l.Out = &output
	defer teardown()

	l.Info("show my debug")
	assert.Contains(t, output.String(), "LOGGER")
}

func TestDebug(t *testing.T) {
	setup()
	defer teardown()

	Debug("show my debug!")
	Debug("%s, %d", "show my debug!", 5)
	assert.Empty(t, output.String())

	Log.Level = logrus.DebugLevel
	Debug("show my debug!")
	assert.NotEmpty(t, output.String())
	Debug("%s, %d", "show my debug!", 5)
	assert.NotEmpty(t, output.String())
}

func TestInfo(t *testing.T) {
	setup()
	defer teardown()

	Info("show my info!")
	Info("%s, %d", "show my info!", 5)
	assert.NotEmpty(t, output.String())

	Log.Level = logrus.DebugLevel
	Info("show my debug!")
	Info("%s, %d", "show my debug!", 5)
	assert.NotEmpty(t, output.String())
}

func TestWarning(t *testing.T) {
	setup()
	defer teardown()

	Warning("show my warning!")
	Warning("%s, %d", "show my warning!", 5)
	assert.NotEmpty(t, output.String())
}

func TestError(t *testing.T) {
	setup()
	defer teardown()

	Error(errors.New("show my errors!"))
	assert.NotEmpty(t, output.String())
}

func TestNoColor(t *testing.T) {
	f := LogFormatter{
		Tty:    false,
		Prefix: "LOGGER",
	}

	l := logrus.Logger{
		Out:       &output,
		Formatter: &f,
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	l.WithFields(logrus.Fields{
		"prefix": "OK",
		"file":   1,
		"error":  errors.New("error here"),
	}).Info("OK")
	assert.Contains(t, output.String(), "LOGGER")
}

func TestFormater(t *testing.T) {
	var output bytes.Buffer

	f := LogFormatter{
		Tty:         false,
		Prefix:      "LOGGER",
		ForceColors: false,
	}

	log := &logrus.Logger{
		Out:       &output,
		Formatter: &f,
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	log.WithFields(logrus.Fields{
		"prefix": "OK",
		"file":   1,
		"time":   1,
		"msg":    1,
		"level":  1,
		"logger": 1,
		"error":  errors.New("error here"),
	}).Error("OK")

	assert.Contains(t, output.String(), "fields.time")
	assert.Contains(t, output.String(), "fields.msg")
	assert.Contains(t, output.String(), "fields.level")
	assert.Contains(t, output.String(), "fields.logger")
}

func setup() {
	f := LogFormatter{
		Tty:         true,
		Prefix:      "LOGGER",
		ForceColors: true,
	}

	Log = &logrus.Logger{
		Out:       &output,
		Formatter: &f,
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
}

func teardown() {
	Log.Out = os.Stdout
}
