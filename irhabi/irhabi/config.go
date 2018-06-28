package irhabi

import (
	"io"
	"os"

	"github.com/alfatih/irhabi/env"
)

var (
	// Config store current configuration values.
	Config = loadConfig()

	// DefaultWriter Default writer of applications
	DefaultWriter io.Writer = os.Stdout
)

// Application configuration variable
type config struct {
	DebugMode  bool   // Switch debug mode for production or development
	JwtSecret  string // Secret key for Json web token algorithm
	GzipEnable bool   // State of gzip compression
	Host       string // IP Application will run, default is 0.0.0.0:8080
	DbEngine   string // Database engines
	DbHost     string // IP Database server, default is 0.0.0.0:3306
	DbName     string // Database name will be used
	DbUser     string // Database username
	DbPassword string // Database password
}

// loadConfig set config value from environment variable.
// If not exists, it will have a default values.
func loadConfig() *config {
	c := new(config)

	c.DebugMode = env.GetBool("APP_DEBUGMODE", true)
	c.JwtSecret = env.GetString("APP_JWT_SECRET", "V3ryF*ck1ngS3cur3")
	c.GzipEnable = env.GetBool("APP_GZIP", false)
	c.Host = env.GetString("SERVER_HOST", "0.0.0.0:8500")
	c.DbEngine = env.GetString("DB_ENGINE", "mysql")
	c.DbHost = env.GetString("DB_HOST", "0.0.0.0:3306")
	//c.DbName = env.GetString("DB_NAME", "konektifa_app")
	c.DbUser = env.GetString("DB_USER", "root")
	c.DbPassword = env.GetString("DB_PASS", "")

	return c
}

// IsDebug returns true if the framework is running in debug mode.
// set environtment variable to release for disable debug.
func IsDebug() bool {
	return Config.DebugMode
}
