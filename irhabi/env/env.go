// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package env

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Try to load .env files in main directory.
func init() {
	AppPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configFile := filepath.Join(AppPath, ".env")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("No Config File Loaded, Using Default Config ...")
	} else {
		Load(configFile)
	}
}

// GetString retrieves the value of the environment variable named
// by the key. If the variable is present and not nil in the environment the
// value is returned as string.
// Otherwise will environment variable will be set as given defaultValue.
func GetString(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	os.Setenv(key, defaultValue)

	return defaultValue
}

// GetInt retrieves the value of the environment variable named
// by the key. If the variable is present and not nil in the environment the
// value is returned as int.
// Otherwise will environment variable will be set as given defaultValue.
func GetInt(key string, defaultValue int) int {
	p, _ := strconv.ParseInt(os.Getenv(key), 10, 32)
	if val := int(p); val != 0 {
		return val
	}
	os.Setenv(key, strconv.Itoa(defaultValue))

	return defaultValue
}

// GetBool retrieves the value of the environment variable named
// by the key. If the variable is present and not nil in the environment the
// value is returned as bool.
// Otherwise will environment variable will be set as given defaultValue.
func GetBool(key string, defaultValue bool) bool {
	if v := os.Getenv(key); v != "" {
		if v == "true" {
			return true
		}

		return false
	}
	os.Setenv(key, strconv.FormatBool(defaultValue))

	return defaultValue
}

// Load will read your env file(s) and load them into ENV for this process.
// Call this function as close as possible to the start of your program (ideally in main)
// If you call Load without any args it will default to loading .env in the current path
// You can otherwise tell it which files to load (there can be more than one) like
// env.Load("fileone", "filetwo")
// It's important to note that it WILL NOT OVERRIDE an env variable that already exists -
// consider the .env file to set dev vars or sensible defaults
func Load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return // return early on a spazout
		}
	}
	return
}

func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	for key, value := range envMap {
		if os.Getenv(key) == "" || overload {
			os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	envMap = make(map[string]string)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			key, value, err := parseLine(fullLine)

			if err == nil {
				envMap[key] = value
			}
		}
	}
	return
}

func parseLine(line string) (key string, value string, err error) {
	if len(line) == 0 {
		err = errors.New("zero length string")
		return
	}

	// ditch the comments (but keep quoted hashes)
	if strings.Contains(line, "#") {
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		var segmentsToKeep []string
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}

			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}

		line = strings.Join(segmentsToKeep, "#")
	}

	// now split key from value
	splitString := strings.SplitN(line, "=", 2)

	if len(splitString) != 2 {
		// try yaml mode!
		splitString = strings.SplitN(line, ":", 2)
	}

	if len(splitString) != 2 {
		err = errors.New("Can't separate key from value")
		return
	}

	// Parse the key
	key = splitString[0]
	if strings.HasPrefix(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}
	key = strings.Trim(key, " ")

	// Parse the value
	value = splitString[1]
	// trim
	value = strings.Trim(value, " ")

	// check if we've got quoted values
	if strings.Count(value, "\"") == 2 || strings.Count(value, "'") == 2 {
		// pull the quotes off the edges
		value = strings.Trim(value, "\"'")

		// expand quotes
		value = strings.Replace(value, "\\\"", "\"", -1)
		// expand newlines
		value = strings.Replace(value, "\\n", "\n", -1)
	}

	return
}

func isIgnoredLine(line string) bool {
	trimmedLine := strings.Trim(line, " \n\t")
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}
