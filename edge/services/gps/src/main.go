/*
Horizon GPS Data REST API

This REST server enables workload access to static location or dynamic
GPS receiver location and satellite data.

Currently the data available through this API includes only:

* the surface location (latitude and longitude) of the node using the
  "decimal coordinates" system.  GPS decimal coordinates have the
  following properties:
  - latitudes are always in the range between -90.0 deg to 90.0 deg, where:
    - latitudes between 0 deg and 90 deg are Northern hemisphere, and
    - latitudes between 0 deg and -90 deg are Southern hemisphere.
  - longitudes are always in the range between -180.0 deg to 180.0 deg, where:
    - longitudes between 0 deg and 180 deg are East of Greenwich meridian, and
    - longitudes between 0 deg and -180 deg are West of Greenwich meridian.

* the vertical location (elevation, or altitude) of the node relative to
  global sea level.  When not provided, this is computed using the USGS
  public elevation query REST API.

* location accuracy information.  This is always 0 when GPS hardware is
  being used (GPS data is never obfuscated by this code) or when the
  location is estimated using the public IP address of the node.  When
  location obfuscation of a static location is requested, the accuracy
  must be specified in kilometers.  Note that location obfuscation is
  always computed only once, creating a new and persistent (obfuscated)
  static location for the node.  As a result the gps REST API cannot be
  used repeatedly in concert with a Monte Carlo technique to derive a
  more precise location defeating the obfuscation.

* the source of the location data (i.e., from "GPS", or a specified
  "Manual" location, or an "Estimated" location derived from the public
  IP address of the node).  Note that when GPS hardware is being used, and
  a GPS *fix* has not been achieved, the source will be given as
  "Searching".  If previously aquired location data is available, that
  old data is provided while searching for a "fix".

* data about the discovered satellites (only if using GPS, of course)

Notice that the location of the node as provided by this API may be dynamic
(provided by GPS hardware), or it may be static.  It may be relatively
precise (when dynamic with a "fix", or when precisely provided statically),
or it may be approximate (derived from the public IP address of the node,
or deliberately obfuscated within a specified distance from a statically
provided location), or it may be completely fictional (which is useful
for testing and other purposes).

The various different ways that this service operates are controlled by
runtime configuration provided by variables in the process environment.
For example, this microservice could be started as follows:

    $ env HZN_GPS_PORT=12345 \
          HZN_USE_GPS=false \
          HZN_LAT=37.0 \
          HZN_LON=-121.0 \
          HZN_LOCATION_ACCURACY_KM=10.0 \
          go run ...

The example above configures manually specified (lat/lon) coordinates and
disables the use of any GPS hardware that might be present.  It also
configures location obfuscation and sets the location accuracy within 10km
of the statically provided location.

Given the above server configuration, a typical request/response
interaction with the service should be similar to what is shown below:

    $ curl -s ...:12345/v1/gps | jq
    {
      "location": {
        "latitude": 36.973370118640865,
        "longitude": -121.00122383303518,
        "elevation": 183.69,
        "accuracy_km": 10,
        "loc_source": "Estimated",
        "loc_last_update": 1486401329
      },
      "satellites": null
    }
    $

If you have the appropriate GPS hardware atached you could instead run the
server as follows to enable the gps microservice to use it:

    $ env HZN_GPS_PORT=12345 \
          HZN_USE_GPS=true
          go run ...

Note that the environment variables HZN_LON, HZN_LAT, HZN_LOCATION_ACCURACY
are ignored when HZN_USE_GPS is set, so these variables are not being set
in this example.

Since the above configuration tells the server that GPS hardware is
available, the expected response would typically be similar to
that shown below (assuming the GPS hardware has acquired a "fix"):

    $ curl -s ...:12345/v1/gps | jq
    {
      "location": {
        "latitude": 37.273474,
        "longitude": -121.880242,
        "elevation": 37.3,
        "accuracy_km": 0,
        "loc_source": "GPS",
        "loc_last_update": 1486416035
      },
      "satellites": [
        {
          "PRN": 1,
          "az": 195,
          "el": 0,
          "ss": 0,
          "used": false
        },
        {
          "PRN": 3,
          "az": 175,
          "el": 56,
          "ss": 18,
          "used": true
        },

  ... <likely with many more satellites listed>

        {
        {
          "PRN": 23,
          "az": 1,
          "el": 70,
          "ss": 15,
          "used": true
        }
      ]
    }
    $

Written by Glen Darling (glendarling@us.ibm.com), Oct. 2016
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	"time"
	"path/filepath"
	"strings"

	// Local modules
	"envutil"
	"gpsdc"
	"hgps"
	"logutil"
	"web"
)

// DEBUG:
//   0 : No debug output
//   1 : Trace GPS monitor fix status
const DEBUG = 1

var DEV_FILE_PATTERNS = []string{"/dev/ttyACM*", "/dev/ttyAMA*", "/dev/cu.usb*", "/dev/tty.usb*"}

// Configuration environment variable names
const (
	GPS_PORT             = "HZN_GPS_PORT"
	USE_GPS              = "HZN_USE_GPS"
	LAT                  = "HZN_LAT"
	LON                  = "HZN_LON"
	LOCATION_ACCURACY_KM = "HZN_LOCATION_ACCURACY_KM"
	GPS_DEBUG            = "GPS_DEBUG"
)

// Configuration variable default values
const (
	DEFAULT_GPS_PORT             = 80
	DEFAULT_USE_GPS              = false
	DEFAULT_LATITUDE             = 0.0
	DEFAULT_LONGITUDE            = 0.0
	DEFAULT_LOCATION_ACCURACY_KM = 0.0
	DEFAULT_GPS_DEBUG            = 0
)

// GPS daemon configuration (port 2947 is the default)
const (
	DAEMON_PORT = 2947
	DAEMON_PATH = "/usr/sbin/gpsd"
)

// Given a list of device path patterns, return the 1st device match on this system, to start gpsd with.
// Returns the empty string if none found.
func findDevFile(devFilePatterns []string) string {
	for _, devPattern := range devFilePatterns {
		logutil.LogDebug("Looking for %v", devPattern)
		if matches, _ := filepath.Glob(devPattern); len(matches) > 0 {
			logutil.LogDebug("Found '%v'", matches[0])
			return matches[0]
		}
	}
	return ""
}

// Run an OS cmd and return any error (exit code, etc.)
func runCmd(commandString string, args ...string) error {
	cmd := exec.Command(commandString, args...)
	if err := cmd.Run(); err != nil { return err }
	return nil
}

// Connect to the gpsd socket using the monitor library, then loop forever
func run_gps_monitor(cache *hgps.GpsCache) {
	logutil.Log("Looking for GPS device file...")
	var devFile string
	var estimate int = 0
	searchInterval := 5 	// seconds
	logEvery := 100    // searches
	searchNum := logEvery   // start it off at the limit so we log the 1st failure
	for {
		devFile = findDevFile(DEV_FILE_PATTERNS)
		if devFile != "" {
			logutil.Logf("Found GPS device file %v", devFile)
			break
		}

		// can not start gpsd because no sensor. If we have already stored the static location, we can serve that for now
		if !cache.IsLocationSet() {
			logutil.Log("ERROR: Did not find a GPS /dev file that we support. Can not fall back to location from specified environment variables, because they were not set. Exiting.")
			
			// Since the default is to use the gps hardware, when it's not found (perhaps on mac)
			// the program will then fall back on the estimated coordinated from IP address
			lat, lon, err := cache.EstimateLocation()
			
			if nil != err {
				logutil.Logf("Houston, we have a problem...")
				logutil.Logf("ERROR: Did not find a GPS /dev file that we support and location estimation using public IP address failed.")
				os.Exit(2)
			}

			envutil.Cfg.LAT = lat
			envutil.Cfg.LON = lon
			envutil.Cfg.LOCATION_ACCURACY_KM = DEFAULT_LOCATION_ACCURACY_KM
			
			cache.SetLocation(envutil.Cfg.LAT, envutil.Cfg.LON, -1.0)
			cache.SetLocationSource(hgps.ESTIMATED)
			estimate = 1
			break
		}
		if searchNum >= logEvery {
			logutil.Logf("Did not find a GPS /dev file that we support. Will continue to look every %v seconds, and serve location from the specified environment variables for now...", searchInterval)
			searchNum = 1
		} else { searchNum++ }
		time.Sleep(time.Duration(searchInterval) * time.Second)
	}
	// runCmd blocks till the cmd returns, but gpsd goes into daemon mode and returns quickly
	logutil.Log("Starting gpsd...")
	if err := runCmd(DAEMON_PATH, devFile); err != nil {
		logutil.Logf("GPS daemon failed to start: %v. Can not serve location from GPS sensor.", err)
		if cache.IsLocationSet() {
			logutil.Log("Serving location from specified environment variables.")
			return
		} else {
			logutil.Log("ERROR: Can not fall back to location from specified environment variables, because they were not set. Exiting.")
			os.Exit(2)
		}
	}
	if (estimate != 1){
		cache.SetLocationSource(hgps.SEARCHING)
	}
	logutil.Logf("CACHE: %s", cache.GetConfiguration())

	logutil.Log("Starting GPS monitor...")
	var gps_monitor *gpsdc.Session
	for {
		var dial_err error
		// Note this dial takes a while to time out if it can not connect
		gps_monitor, dial_err = gpsdc.Dial(fmt.Sprintf("localhost:%d", DAEMON_PORT))
		if dial_err == nil { break }
		logutil.Log("GPS monitor was unable to connect to the gpsd daemon socket.  Sleeping before retrying...")
		time.Sleep(5000 * time.Millisecond)
	}
	logutil.Log("GPS monitor has connected to the gpsd daemon socket.")

	// var received_fix bool = false
	gps_monitor.AddFilter("SKY", func(r interface{}) {
		report := r.(*gpsdc.SKYReport)
		cache.SetSatellites(report) 	// The debug output is inside this method
	})
	gps_monitor.AddFilter("TPV", func(r interface{}) {
		report := r.(*gpsdc.TPVReport)
		// Ignore any 0.0, 0.0, 0.0 location
		if 0.0 != report.Lat && 0.0 != report.Lon && report.Alt > 0.0 {
			if cache.GetLocationSource() != hgps.GPS {
				// changed from no fix to fix
				// received_fix = true
				logutil.Logf("GPS received fix: Time:%q, Lat=%f, Lon=%f, Elev=%f", report.Time, report.Lat, report.Lon, report.Alt)
				cache.SetLocationSource(hgps.GPS)
			} else {
				logutil.LogDebug("GPS received fix: Time:%q, Lat=%f, Lon=%f, Elev=%f", report.Time, report.Lat, report.Lon, report.Alt)
			}
			cache.SetLocation(report.Lat, report.Lon, report.Alt)
		} else {
			// GPS hardware lost its fix
			if cache.GetLocationSource() != hgps.SEARCHING {
				// changed from fix to no fix
				// received_fix = false
				logutil.Logf("GPS lost fix.  Searching...")
				cache.SetLocationSource(hgps.SEARCHING)
			} else {
				logutil.LogDebug("GPS does not have fix, searching...")
			}
		}
	})

	gps_monitor.SendCommand(gpsdc.DefaultWatchCommand)
	done := gps_monitor.Listen()
	<-done

	// This code cannot be reached
	logutil.Log("GPS monitoring loop unexpectedly ended.")
	os.Exit(2)
}

func main() {
	/* 
	  Variables set in the process environment configure the behavior of
	  the "gps" microservice, as described below:

	  * HZN_USE_GPS==true: this means GPS hardware access is enabled
	    - HZN_LAT, HZN_LON can be optionally set as fallback when the
	      gps sensor is not available
	    - HZN_LOCATION_ACCURACY_KM is ignored in this situation (because
	      they approved sharing exact coords)

	  * HZN_USE_GPS==false (the default): GPS hardware access is disabled
	    - HZN_LAT, HZN_LON may be set (provided in the process environment)
	      + if they were provided, and either or both were nonzero, then
	        if HZN_LOCATION_ACCURACY_KM nonzero the latitude and longitude
                provided will be obfuscated within a radiius onf this size
	    - if HZN_LAT, HZN_LON are not set, or are 0.0, then the location
              will be estimated by using the public Internet IP address of the
              node and HZN_LOCATION_ACCURACY_KM will be ignored (the location is
              already obfuscated by estimation from the IP address in this case)
	*/
	envutil.Cfg.GPS_PORT = envutil.GetInt(GPS_PORT, DEFAULT_GPS_PORT, false)
	envutil.Cfg.USE_GPS = envutil.GetBool(USE_GPS, DEFAULT_USE_GPS, true)
	envutil.Cfg.LAT = envutil.GetFloat(LAT, DEFAULT_LATITUDE, false)
	envutil.Cfg.LON = envutil.GetFloat(LON, DEFAULT_LONGITUDE, false)
	envutil.Cfg.LOCATION_ACCURACY_KM = envutil.GetFloat(LOCATION_ACCURACY_KM, DEFAULT_LOCATION_ACCURACY_KM, false)
	logutil.GPS_DEBUG = int(envutil.GetInt(GPS_DEBUG, DEFAULT_GPS_DEBUG, false))

	// Create the Horizon GPS cache that the REST server reads from
	cache := hgps.New()

	// Configure based upon the location data source
	source := hgps.MANUAL
	if envutil.Cfg.USE_GPS {
		source = hgps.SEARCHING
		// Obfuscation is prohibited when source is GPS hardware
		envutil.Cfg.LOCATION_ACCURACY_KM = DEFAULT_LOCATION_ACCURACY_KM
	} else {
		// No GPS.  If no static location is provided, estimate it:
		if envutil.Cfg.LAT == DEFAULT_LATITUDE || envutil.Cfg.LON == DEFAULT_LONGITUDE {
			logutil.Log("Estimating node location using public IP address.")
			source = hgps.ESTIMATED
			lat, lon, err := cache.EstimateLocation()
			if nil != err {
				logutil.Logf("ERROR: GPS hardware is not enabled, valid static latitude and longitude were not provided, and location estimation using public IP address failed: %v", err)
				os.Exit(2)
			}
		        envutil.Cfg.LAT = lat
			envutil.Cfg.LON = lon

			// Obfuscation is prohibited when location is estimated
			envutil.Cfg.LOCATION_ACCURACY_KM = DEFAULT_LOCATION_ACCURACY_KM
		}
	}

	// Log the configuration to aid in development, test and debugging
	logutil.Log("Blue Horizon gps service configuration:")
	logutil.Logf("  %v=%v", GPS_PORT, envutil.Cfg.GPS_PORT)
	logutil.Logf("  %v=%v", USE_GPS, envutil.Cfg.USE_GPS)
	logutil.Logf("  %v=%v", LAT, envutil.Cfg.LAT)
	logutil.Logf("  %v=%v", LON, envutil.Cfg.LON)
	logutil.Logf("  %v=%v", LOCATION_ACCURACY_KM, envutil.Cfg.LOCATION_ACCURACY_KM)
	logutil.Logf("  %v: %v", "(location source)", source)
	logutil.Logf("  %v: %v", "(/dev patterns)", strings.Join(DEV_FILE_PATTERNS, ", "))

	// Pre-populate the cache with source and accuracy
	cache.SetLocationSource(source)
        cache.SetLocationAccuracyInKm(envutil.Cfg.LOCATION_ACCURACY_KM)

        // If location was statically provided or estimated from IP address
	// cache it now (and only once) for the REST server to consume.
        if envutil.Cfg.LAT != DEFAULT_LATITUDE && envutil.Cfg.LON != DEFAULT_LONGITUDE {
		// Passing a negative elevation here causes the elevation
		// to be discovered using the USGS web service.
		cache.SetLocation(envutil.Cfg.LAT, envutil.Cfg.LON, -1.0)

		// Since this location was cached, place it into the log too
		cached_lat, cached_lon, cached_elev_m, cached_accuracy_km := cache.GetLocation()
		logutil.Logf("CACHED: Static location: Latitude=%f, Longitude=%f, Elevation=%fm, Accuracy=%fkm", cached_lat, cached_lon, cached_elev_m, cached_accuracy_km)
	}

	// Log the cache configuration
        logutil.Logf("CACHED: %s", cache.GetConfiguration())

	// If statically configured, log the whole cache
	if !envutil.Cfg.USE_GPS {
		logutil.Logf("CACHED: %s", string(cache.GetAsJSON()))
	}

	// If GPS hardware access is enabled, start up a goroutine to
	// continuously read updates from the gpsd daemon and repeatedly
	// cache the latest location for the REST server to consume.
	if envutil.Cfg.USE_GPS {
		// Interact with gpsd and continuously cache its location data
		go run_gps_monitor(cache)
	}
	
	// Set up a signal handler for a cleaner exit on SIGINT
	sig_channel := make(chan os.Signal, 1)
	signal.Notify(sig_channel, syscall.SIGINT)
	go func() {
		sig := <-sig_channel
		fmt.Printf("\n") // Newline to follow the ^C on the console
		logutil.Log(fmt.Sprintf("Received signal %q.  Exiting.", sig))
		// @@@ What to cleanup before exit?
		os.Exit(0)
	}()

	// And finally, start a REST service to provide access to the GPS data
	web.StartWebServer(envutil.Cfg.GPS_PORT, cache)
}
