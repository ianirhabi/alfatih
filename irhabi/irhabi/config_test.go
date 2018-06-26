package irhabi

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("APP_DEBUGMODE", "")
	os.Setenv("APP_JWT_SECRET", "")
	os.Setenv("APP_GZIP", "")
	os.Setenv("SERVER_HOST", "")
	os.Setenv("DB_ENGINE", "")
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_NAME", "")
	os.Setenv("DB_USER", "")
	os.Setenv("DB_PASS", "")

	c := loadConfig()
	var tests = []struct {
		value    interface{}
		expected interface{}
	}{
		{c.DebugMode, true},
		{c.JwtSecret, "V3ryF*ck1ngS3cur3"},
		{c.GzipEnable, false},
		{c.Host, "0.0.0.0:8080"},
		{c.DbEngine, "mysql"},
		{c.DbHost, "0.0.0.0:3306"},
		{c.DbName, "konektifa_app"},
		{c.DbUser, "root"},
		{c.DbPassword, ""},
	}

	for _, test := range tests {
		assert.Equal(t, test.value, test.expected)
	}
}

func TestIsDebugging(t *testing.T) {
	Config.DebugMode = true
	assert.True(t, IsDebug())

	Config.DebugMode = false
	assert.False(t, IsDebug())
}
