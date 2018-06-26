// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailer

import (
	"git.qasico.com/cuxs/env"
)

// Config represents all configurable mailer data smtp credentials
var Config *configMailer

// configMailer type to store mailer configuration
type configMailer struct {
	SMTPIdentity string
	SMTPUsername string
	SMTPPassword string
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderName   string
}

// ReadEnv set all configurable data from env variable
func ReadEnv() {
	Config = &configMailer{
		SMTPIdentity: env.GetString("SMTP_IDENTITY", ""),
		SMTPUsername: env.GetString("SMTP_USERNAME", "euvoriaMail"),
		SMTPPassword: env.GetString("SMTP_PASSWORD", "euvoriacom3S15cwXPm"),
		SMTPHost:     env.GetString("SMTP_HOST", "smtp.sendgrid.net"),
		SMTPPort:     env.GetString("SMTP_PORT", "587"),
		SenderEmail:  env.GetString("SMTP_SENDER", "noreply@konektifa.com"),
		SenderName:   env.GetString("SMTP_SENDER_NAME", "Konektifa System"),
	}
}

func init() {
	ReadEnv()
}
