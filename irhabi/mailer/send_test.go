package mailer

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTo1  = "to1@example.com"
	testTo2  = "to2@example.com"
	testFrom = "from@example.com"
	testBody = "Test message"
	testMsg  = "To: " + testTo1 + ", " + testTo2 + "\r\n" +
		"From: " + testFrom + "\r\n" +
		"Mime-Version: 1.0\r\n" +
		"Date: Wed, 25 Jun 2014 17:46:00 +0000\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"Content-Transfer-Encoding: quoted-printable\r\n" +
		"\r\n" +
		testBody
)

type mockSender SendFunc

func (s mockSender) Send(from string, to []string, msg io.WriterTo) error {
	return s(from, to, msg)
}

type mockSendCloser struct {
	mockSender
	close func() error
}

func (s *mockSendCloser) Close() error {
	return s.close()
}

func TestSend(t *testing.T) {
	s := &mockSendCloser{
		mockSender: stubSend(t, testFrom, []string{testTo1, testTo2}, testMsg),
		close: func() error {
			t.Error("Close() should not be called in Send()")
			return nil
		},
	}
	err := Send(s, getTestMessage())
	assert.NoError(t, err)
}

func getTestMessage() *Message {
	m := NewMessage()
	m.SetHeader("From", testFrom)
	m.SetHeader("To", testTo1, testTo2)
	m.SetBody("text/plain", testBody)

	return m
}

func stubSend(t *testing.T, wantFrom string, wantTo []string, wantBody string) mockSender {
	return func(from string, to []string, msg io.WriterTo) error {
		assert.Equal(t, wantFrom, from)
		assert.True(t, reflect.DeepEqual(to, wantTo), fmt.Sprintf("invalid to, got %v, want %v", to, wantTo))

		buf := new(bytes.Buffer)
		_, err := msg.WriteTo(buf)
		assert.NoError(t, err)

		compareBodies(t, buf.String(), wantBody)
		return nil
	}
}
