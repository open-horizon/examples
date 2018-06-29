// Example for reading/writing the SDR Compose Postgresql DB using go.
// See README.md for setup requirements.

package main

import (
	"fmt"
	"os"
	"strings"
	//"strconv"
	"database/sql"
	_ "github.com/lib/pq"
)

func Usage(exitCode int) {
	fmt.Printf("Usage: %s\n\nEnvironment Variables: SDR_DB_PASSWORD, SDR_DB_USER, SDR_DB_HOST, SDR_DB_PORT, SDR_DB_NAME\n", os.Args[0])
	os.Exit(exitCode)
}

func main() {
	// The lib/pq package is a postgresql driver for the standard database/sql package. See https://godoc.org/github.com/lib/pq for details.

	// Get the env var values we need to connect to our db
	pw := RequiredEnvVar("SDR_DB_PASSWORD", "")
	host := RequiredEnvVar("SDR_DB_HOST", "")
	port := RequiredEnvVar("SDR_DB_PORT", "")
	user := RequiredEnvVar("SDR_DB_USER", "admin")
	dbName := RequiredEnvVar("SDR_DB_NAME", "sdr")

	// Connect to db
	connStr := "postgres://"+user+":"+pw+"@"+host+":"+port+"/"+dbName+"?sslmode=require"
	db, err := sql.Open("postgres", connStr)
	ExitOnErr(err)
	defer db.Close()

	// Increment the numberofmentions for noun==wedding
	stmt, err := db.Prepare("update nouns set numberofmentions = numberofmentions + 1, timeupdated = CURRENT_TIMESTAMP where noun=$1")
	ExitOnErr(err)
	defer stmt.Close()
	res, err := stmt.Exec("wedding")
	ExitOnErr(err)
	affect, err := res.RowsAffected()
	ExitOnErr(err)
	fmt.Printf("updated %d row(s)\n", affect)

	/* Try updating with QueryRow(). The update actually happens, but the row isn't returned, and there isn't a good way to get errors....
	var noun, sentiment, timeupdated string
	var numberofmentions int
	err = db.QueryRow("update nouns set numberofmentions = numberofmentions + 1 where noun=$1", "wedding").Scan(&noun, &sentiment, &numberofmentions, &timeupdated)
	ExitOnErr(err)
	fmt.Printf("updated row: %s, %s, %d, %s\n", noun, sentiment, numberofmentions, timeupdated)
	*/

	// Could also insert a row like:
	// err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "foo", "2012-12-09").Scan(&lastInsertId)

	// Read db
	rows, err := db.Query("SELECT * FROM nouns")
	ExitOnErr(err)
	defer rows.Close()
	for rows.Next() {
		var noun, sentiment, timeupdated string
		var numberofmentions int
		err := rows.Scan(&noun, &sentiment, &numberofmentions, &timeupdated)
		ExitOnErr(err)
		fmt.Printf("queried row: %s, %s, %d, %s\n", noun, sentiment, numberofmentions, timeupdated)
	}
	ExitOnErr(rows.Err())
}



var VerboseBool bool

func Verbose(msg string, args ...interface{}) {
	if !VerboseBool {
		return
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, "[verbose] "+msg, args...) // send to stderr so it doesn't mess up stdout if they are piping that to jq or something like that
}

// RequiredEnvVar gets an env var value. If a default value is not supplied and the env var is not defined, a fatal error is displayed.
func RequiredEnvVar(name, defaultVal string) string {
	v := os.Getenv(name)
	if defaultVal != "" {
		v = defaultVal
	}
	if v == "" {
		fmt.Printf("Error: environment variable '%s' must be defined.\n", name)
		os.Exit(2)
	}
	return v
}

func ExitOnErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(2)
	}
}

