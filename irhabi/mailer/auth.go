package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
)

// Auth is an smtp.Auth that implements the LOGIN authentication mechanism.
type Auth struct {
	username string
	password string
	host     string
}

// Start initialing authentication mechanism to the server.
func (a *Auth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		advertised := false
		for _, mechanism := range server.Auth {
			if mechanism == "LOGIN" {
				advertised = true
				break
			}
		}
		if !advertised {
			return "", nil, errors.New("mailer: unencrypted connection")
		}
	}
	if server.Name != a.host {
		return "", nil, errors.New("mailer: wrong host name")
	}
	return "LOGIN", nil, nil
}

// Next using previous authentication mechanism
func (a *Auth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil
	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("mailer: unexpected server challenge: %s", fromServer)
	}
}
