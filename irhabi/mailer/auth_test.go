package mailer

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testUser = "username"
	testPwd  = "password"
	testHost = "mail.qasico.com"
)

type authTest struct {
	auths      []string
	challenges []string
	tls        bool
	wantData   []string
	wantError  bool
}

func TestNoAdvertisement(t *testing.T) {
	testLoginAuth(t, &authTest{
		auths:     []string{},
		tls:       false,
		wantError: true,
	})
}

func TestNoAdvertisementTLS(t *testing.T) {
	testLoginAuth(t, &authTest{
		auths:      []string{},
		challenges: []string{"Username:", "Password:"},
		tls:        true,
		wantData:   []string{"", testUser, testPwd},
	})
}

func TestLogin(t *testing.T) {
	testLoginAuth(t, &authTest{
		auths:      []string{"PLAIN", "LOGIN"},
		challenges: []string{"Username:", "Password:"},
		tls:        false,
		wantData:   []string{"", testUser, testPwd},
	})
}

func TestLoginTLS(t *testing.T) {
	testLoginAuth(t, &authTest{
		auths:      []string{"LOGIN"},
		challenges: []string{"Username:", "Password:"},
		tls:        true,
		wantData:   []string{"", testUser, testPwd},
	})
}

func TestAuthError(t *testing.T) {
	auth := &Auth{
		username: testUser,
		password: testPwd,
		host:     testHost,
	}
	server := &smtp.ServerInfo{
		Name: "wrongHost",
		TLS:  false,
		Auth: []string{"LOGIN"},
	}
	_, _, err := auth.Start(server)
	assert.Error(t, err)

	x, _ := auth.Next([]byte(""), false)
	assert.Nil(t, x)

	_, err = auth.Next([]byte(""), true)
	assert.Error(t, err)
}

func testLoginAuth(t *testing.T, test *authTest) {
	auth := &Auth{
		username: testUser,
		password: testPwd,
		host:     testHost,
	}
	server := &smtp.ServerInfo{
		Name: testHost,
		TLS:  test.tls,
		Auth: test.auths,
	}
	proto, toServer, err := auth.Start(server)

	if test.wantError {
		assert.Error(t, err, fmt.Sprintf("Auth.Start(): %v", err))
	} else {
		assert.NoError(t, err, fmt.Sprintf("Auth.Start(): %v", err))
	}

	// noting to test any more because its error :)
	if err != nil && test.wantError {
		return
	}

	var i = 0
	var got = string(toServer)
	assert.Equal(t, "LOGIN", proto)
	assert.Equal(t, test.wantData[i], got)

	for _, challenge := range test.challenges {
		i++
		assert.False(t, i >= len(test.wantData), fmt.Sprintf("unexpected challenge: %q", challenge))

		toServer, err = auth.Next([]byte(challenge), true)
		assert.NoError(t, err, fmt.Sprintf("Auth.Auth(): %v", err))

		got = string(toServer)
		assert.Equal(t, test.wantData[i], got)
	}
}
