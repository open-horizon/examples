package logutil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Set in main via an env var
var GPS_DEBUG int = 0

// Log the passed string (maybe do more with this function someday?)
func Log(args ...interface{}) {
	log.Println(args...)
}

func Logf(formatStr string, args ...interface{}) {
	log.Printf(formatStr, args...)
}

// Log if GPS_DEBUG is non-zero (the lowest level of debug)
func LogDebug(formatStr string, args ...interface{}) {
	if GPS_DEBUG == 0 { return }
	log.Printf("DEBUG: "+formatStr, args...)
}

// Log REST requests
func LogRestRequest(method string, message string) {
	Log(fmt.Sprintf("REST --> %s: %s", method, message))
}

// Log REST responses
func LogRestResponse(message string) {
	Log(fmt.Sprintf("REST <-- %s", message))
}

// Dump and entire file into the log
func LogFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		Log(fmt.Sprintf("LogFile: Unable to open file: %q.", path))
		return
	}
	defer file.Close()

	Log(fmt.Sprintf("Start of file: %q:", path))
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Log(scanner.Text())
	}
	Log(fmt.Sprintf("End of file: %q.", path))
}

// Log JSON data
func LogJsonBytes(message []byte) {
	content := json_prettyprint(message, "  ", "  ")
	slice := strings.Split(content, "\n")
	for _, line := range slice {
		Log(line)
	}
}

// Log a JSON byte array, formatted for human consumption
func json_prettyprint(in []byte, prefix string, indent string) string {
	var out bytes.Buffer
	pretty_err := json.Indent(&out, in, prefix, indent)
	if pretty_err != nil {
		return "(pretty-printing error)"
	}
	return string(out.Bytes())
}
