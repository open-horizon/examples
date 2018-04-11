package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	// "net"
	"net/http"
	// "os"

	// Local modules
	"hgps"
	"logutil"
)

// DEBUG:
//   0 : No debug output
//   1 : Trace REST service transactions
// const DEBUG = 0

// The gps cache
var cache *hgps.GpsCache = nil

// Format JSON for the human eye
func json_prettyprint(in []byte, prefix string, indent string) string {
	var out bytes.Buffer
	pretty_err := json.Indent(&out, in, prefix, indent)
	if pretty_err != nil {
		return "(pretty-printing error)"
	}
	return string(out.Bytes())
}

func add_json_header(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
}

// Add anti-caching headers to the http.ResponseWriter object
func add_no_cache_header(w http.ResponseWriter) {
	// Note that all of these are required
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Expires", "0")
}

// Set the specified HTTP status code in the http.ResponseWriter object
func set_response_status(w http.ResponseWriter, s int) {
	// Argument s should be one of the RFC2616 constants defined here:
	//     https://golang.org/src/net/http/status.go
	w.WriteHeader(s)
}

// Send an error response, and also log it
func respond_with_error(w http.ResponseWriter,
	err int,
	method string,
	path string,
	details ...string) {
	footer := ""
	if len(details) > 0 {
		footer = fmt.Sprintf("  Details: %s", details[0])
	}
	message := fmt.Sprintf("Received: %s, %s.%s", method, path, footer)
	if logutil.GPS_DEBUG > 0 {
		logutil.LogRestResponse(fmt.Sprintf("%s (%d).  %s", http.StatusText(err), err, message))
	}
	add_json_header(w)
	set_response_status(w, err)
	fmt.Fprintf(w, "{\"error\":%q}\n", message)
}

// Send a success response, and also log it
func respond_with_success(w http.ResponseWriter,
	status int,
	message string) {
	detail := fmt.Sprintf("Success (%d). Sending response: %s", status, message)
	if logutil.GPS_DEBUG > 0 {
		logutil.LogRestResponse(detail)
	}
	add_json_header(w)
	set_response_status(w, status)
	fmt.Fprintf(w, "%s\n", message)
}

// Handle requests to unrecognized URIs
func handle_bad_request(w http.ResponseWriter, r *http.Request) {
	respond_with_error(w, http.StatusBadRequest, r.Method, r.URL.Path)
}

// Handler for requests to URI "/gps/location"
func handle_uri_location(w http.ResponseWriter, r *http.Request) {

	if logutil.GPS_DEBUG > 0 {
		// Log the request before handling it
		logutil.LogRestRequest(r.Method, r.URL.Path)
	}

	// Accept only GET requests to this URI
	if r.Method != "GET" {
		respond_with_error(w, http.StatusMethodNotAllowed, r.Method, r.URL.Path)
		return
	}

	// Add the standard JSON header to the response
	// add_json_header(w)  // respond_with_success does this

	// Prevent this data from being cached (client must always revalidate)
	add_no_cache_header(w)

	// Construct the JSON response string
	response_string := string(cache.GetLocationAsJSON())

	// Send the response (and log it)
	respond_with_success(w, http.StatusOK, response_string)
}

// Handler for requests to URI "/gps/satellites"
func handle_uri_satellites(w http.ResponseWriter, r *http.Request) {

	if logutil.GPS_DEBUG > 0 {
		// Log the request before handling it
		logutil.LogRestRequest(r.Method, r.URL.Path)
	}

	// Accept only GET requests to this URI
	if r.Method != "GET" {
		respond_with_error(w, http.StatusMethodNotAllowed, r.Method, r.URL.Path)
		return
	}

	// Add the standard JSON header to the response
	// add_json_header(w)  // respond_with_success does this

	// Prevent this data from being cached (client must always revalidate)
	add_no_cache_header(w)

	// Construct the JSON response string
	response_string := string(cache.GetSatellitesAsJSON())

	// Send the response (and log it)
	respond_with_success(w, http.StatusOK, response_string)
}

// Handler for requests to URI "/gps"
func handle_uri_gps(w http.ResponseWriter, r *http.Request) {

	if logutil.GPS_DEBUG > 0 {
		// Log the request before handling it
		logutil.LogRestRequest(r.Method, r.URL.Path)
	}

	// Accept only GET requests to this URI
	if r.Method != "GET" {
		respond_with_error(w, http.StatusMethodNotAllowed, r.Method, r.URL.Path)
		return
	}

	// Add the standard JSON header to the response
	// add_json_header(w)  // respond_with_success does this

	// Prevent this data from being cached (client must always revalidate)
	add_no_cache_header(w)

	// Construct the JSON response string
	response_string := string(cache.GetAsJSON())

	// Send the response (and log it)
	respond_with_success(w, http.StatusOK, response_string)
}

func StartWebServer(port int, c *hgps.GpsCache) {

	// Save the cache reference to a global for the REST handlers to use
	cache = c

	// This code used to check whether the specified port is in use, but net.Dial() hangs while the net.Dial in gpsdc.go is
	// waiting to connect to gpsd (which once in a while takes a long time). Beside, http.ListenAndServe() below will fail appropriately
	// if the port is in use.
	// net_dial_target := fmt.Sprintf("localhost:%d", port)
	// _, dial_err := net.Dial("tcp", net_dial_target)
	// if dial_err == nil {
	// 	logutil.Log(fmt.Sprintf("TCP address and port %q are already in use.", net_dial_target))
	// 	os.Exit(1)
	// }

	// Add handlers for the REST APIs supported
	http.HandleFunc("/v1/gps/location", handle_uri_location)
	http.HandleFunc("/v1/gps/satellites", handle_uri_satellites)
	http.HandleFunc("/v1/gps", handle_uri_gps)

	// All others will fall through to this generic "bad request" handler
	http.HandleFunc("/", handle_bad_request)

	// Announce launch
	logutil.Log(fmt.Sprintf("REST server is launching on port %d.", port))

	// Launch server (binding it to all interfaces)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
