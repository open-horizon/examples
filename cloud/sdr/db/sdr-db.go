// Example for reading/writing the SDR Compose Postgresql DB using go.
// See README.md for setup requirements.

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	//"strconv"
	"database/sql"

	_ "github.com/lib/pq"
)

func usage(exitCode int) {
	fmt.Printf("Usage: %s\n\nEnvironment Variables: SDR_DB_PASSWORD, SDR_DB_USER, SDR_DB_HOST, SDR_DB_PORT, SDR_DB_NAME\n", os.Args[0])
	os.Exit(exitCode)
}

func main() {
	// The lib/pq package is a postgresql driver for the standard database/sql package. See https://godoc.org/github.com/lib/pq for details.

	// Get cmd line args passed in, or default them
	noun, sentiment := getArgs()

	// Get the env var values we need to connect to our db
	pw := requiredEnvVar("SDR_DB_PASSWORD", "")
	host := requiredEnvVar("SDR_DB_HOST", "")
	port := requiredEnvVar("SDR_DB_PORT", "")
	user := requiredEnvVar("SDR_DB_USER", "admin")
	dbName := requiredEnvVar("SDR_DB_NAME", "sdr")

	// Connect to db
	connStr := "postgres://" + user + ":" + pw + "@" + host + ":" + port + "/" + dbName + "?sslmode=require"
	db, err := sql.Open("postgres", connStr)
	exitOnErr(err)
	defer db.Close()

	// Increment the numberofmentions for noun==wedding
	// stmt, err := db.Prepare("update globalnouns set numberofmentions = numberofmentions + 1, timeupdated = CURRENT_TIMESTAMP where noun=$1")
	stmt, err := db.Prepare("INSERT INTO globalnouns VALUES ($1, $2, 1, CURRENT_TIMESTAMP) ON CONFLICT (noun) DO UPDATE SET sentiment = ((globalnouns.sentiment * globalnouns.numberofmentions) + $2) / (globalnouns.numberofmentions + 1), numberofmentions = globalnouns.numberofmentions + 1, timeupdated = CURRENT_TIMESTAMP")
	exitOnErr(err)
	defer stmt.Close()
	// res, err := stmt.Exec("wedding")
	res, err := stmt.Exec(noun, sentiment)
	exitOnErr(err)
	affect, err := res.RowsAffected()
	exitOnErr(err)
	fmt.Printf("inserted/updated %d row(s)\n", affect)

	/* Try updating with QueryRow(). The update actually happens, but the row isn't returned, and there isn't a good way to get errors....
	var noun, sentiment, timeupdated string
	var numberofmentions int
	err = db.QueryRow("update nouns set numberofmentions = numberofmentions + 1 where noun=$1", "wedding").Scan(&noun, &sentiment, &numberofmentions, &timeupdated)
	exitOnErr(err)
	fmt.Printf("updated row: %s, %s, %d, %s\n", noun, sentiment, numberofmentions, timeupdated)
	*/

	// Could also insert a row like:
	// err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "foo", "2012-12-09").Scan(&lastInsertId)

	// Read db
	rows, err := db.Query("SELECT * FROM globalnouns")
	exitOnErr(err)
	defer rows.Close()
	for rows.Next() {
		var noun, timeupdated string
		var sentiment float64
		var numberofmentions int
		err := rows.Scan(&noun, &sentiment, &numberofmentions, &timeupdated)
		exitOnErr(err)
		fmt.Printf("queried row: %s, %f, %d, %s\n", noun, sentiment, numberofmentions, timeupdated)
	}
	exitOnErr(rows.Err())
}

var verboseBool bool

func verbose(msg string, args ...interface{}) {
	if !verboseBool {
		return
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, "[verbose] "+msg, args...) // send to stderr so it doesn't mess up stdout if they are piping that to jq or something like that
}

// Get the cmd line args, or default them
func getArgs() (noun string, sentiment float64) {
	noun = "wedding"
	sentiment = 0.2
	if len(os.Args) > 1 {
		noun = os.Args[1]
	}
	if len(os.Args) > 2 {
		var err error
		sentiment, err = strconv.ParseFloat(os.Args[2], 64)
		exitOnErr(err)
	}
	fmt.Printf("setting noun: %s, sentiment: %f\n", noun, sentiment)
	return
}

// requiredEnvVar gets an env var value. If a default value is not supplied and the env var is not defined, a fatal error is displayed.
func requiredEnvVar(name, defaultVal string) string {
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

func exitOnErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(2)
	}
}
