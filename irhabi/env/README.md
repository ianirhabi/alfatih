# cuxs/env

[![build status](https://git.qasico.com/cuxs/env/badges/master/build.svg)](https://git.qasico.com/cuxs/env/commits/master)
[![coverage report](https://git.qasico.com/cuxs/env/badges/master/coverage.svg)](https://git.qasico.com/cuxs/env/commits/master)

## ENV

```go
func GetString(key string, defaultValue string)
// GetString retrieves the value of the environment variable named
// by the key. If the variable is present and not nil in the environment the
// value is returned as string.

func GetInt(key string, defaultValue int)
// GetInt retrieves the value of the environment variable named
// by the key. If the variable is present and not nil in the environment the
// value is returned as int.

func GetBool(key string, defaultValue bool)
// GetBool retrieves the value of the environment variable named
// by the key. If the variable is present and not nil in the environment the
// value is returned as bool.

func Load(filenames ...string) (err error)
// Load will read your env file(s) and load them into ENV for this process.
// Call this function as close as possible to the start of your program (ideally in main)
// If you call Load without any args it will default to loading .env in the current path
// You can otherwise tell it which files to load (there can be more than one) like
// env.Load("fileone", "filetwo")
// It's important to note that it WILL NOT OVERRIDE an env variable that already exists -
// consider the .env file to set dev vars or sensible defaults
```

## Example .env file

```go
### APP
USERNAME=Asep
PASSWORD=123
VISIBILITY=true
```

## Example load .env file

```go
package main

import "git.qasico.com/cuxs/env"
// import env from cuxs

type config struct {
	Username     string
	Password     int
	Visibility   bool
}

func main() {

	var a config

	env.Load(".env")
	// load .env file

	a.Username = env.GetString("USERNAME", "Udin")
	// load USERNAME from .env file using func GetString
	// if in .env file variable USERNAME is nil, default USERNAME is Udin

	a.Password = env.GetInt("PASSWORD", 456)
	// load PASSWORD from .env file using func GetInt
	// if in .env file variable PASSWORD is nil, default PASSWORD is 456

	a.Visibility = env.GetBool("VISIBILITY", false)
	// load VISIBILITY from .env file using func GetBool
	// if in .env file variable VISIBILITY is nil, default VISIBILITY is false

}
```
