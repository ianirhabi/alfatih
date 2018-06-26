// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("SMTP_IDENTITY", "")
	os.Setenv("SMTP_USERNAME", "username")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("SMTP_HOST", "mail.qasico.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_SENDER", "noreply@qasico.com")
	os.Setenv("SMTP_SENDER_NAME", "System Qasico")

	ReadEnv()

	assert.Equal(t, "", Config.SMTPIdentity)
	assert.Equal(t, "username", Config.SMTPUsername)
	assert.Equal(t, "password", Config.SMTPPassword)
	assert.Equal(t, "mail.qasico.com", Config.SMTPHost)
	assert.Equal(t, "587", Config.SMTPPort)
	assert.Equal(t, "noreply@qasico.com", Config.SenderEmail)
	assert.Equal(t, "System Qasico", Config.SenderName)
}
