//
// hgps -- a resource class for caching Horizon GPS data
//
// Cache your Horizon GPS data in a "smart" GpsCache instance.  Cache features:
// - (optional) obfuscate location within specified distance from actual loc
//     (note that this works correctly for any location anywhere on the globe)
// - (optional) can discover elevation when not specified (i.e., ground level)
// - caches GPS fixes -- once a fix is found, coordinates are never again zero
// Note that GpsCache instances are also thread safe.
//
// Usage example:
//     var p *GpsCache = hgps.New()
//     err := p.SetLocationSourceAndAccuracyInKm(source, gps, accuracy_km)
//     t := p.SetLocation(lat, lon, -1.0) // Requests elevation discovery
//     lat, lon, elev := p.GetLocation()
//     json_str := p.GetAsJSON()
//
// Written by Glen Darling (glendarling@us.ibm.com), Oct. 2016
//
package hgps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	// "strconv"
	// "strings"
	"sync"
	"time"

	// Local modules
	"gpsdc"
	"logutil"

	// For location obfuscation
	"github.com/kellydunn/golang-geo"
)

// DEBUG:
//   0 : No debug output from this module
//   1 : Trace all state changes
//   2 : Also dump full state at each state change
//   3 : Also show location obfuscation and altitude discovery
const DEBUG = 0

// How long to delay between retries when seeking the public IP address
const GET_PUBLIC_IP_RETRY_SEC = 10

// How long to delay between retries for location estimate from IP address
const LOC_ESTIMATE_RETRY_SEC = 10

type SourceType string

const (
	MANUAL    SourceType = "Manual"
	ESTIMATED SourceType = "Estimated"
	SEARCHING SourceType = "Searching"
	GPS       SourceType = "GPS"
)

// JSON struct for location data
type LocationData struct {
	Latitude   float64    `json:"latitude" description:"Location latitude"`
	Longitude  float64    `json:"longitude" description:"Location longitude"`
	ElevationM float64    `json:"elevation" description:"Location elevation in meters"`
	AccuracyKM float64    `json:"accuracy_km" description:"Location accuracy in kilometers"`
	LocSource  SourceType `json:"loc_source" description:"Location source (one of: Manual, Estimated, GPS, or Searching)"`
	LastUpdate int64      `json:"loc_last_update" description:"Time of most recent location update (UTC)."`
}

// JSON struct for satellite data
type SatellitesData struct {
	Satellites []SatelliteData `json:"satellites" description:"Array of data for each of the attached satellites."`
}

// JSON struct for data about a single satellite
type SatelliteData struct {
	PRN  float64 `json:"PRN"  description:"PRN ID of the satellite. 1-63 are GNSS satellites, 64-96 are GLONASS satellites, 100-164 are SBAS satellites"`
	Az   float64 `json:"az"   description:"Azimuth, degrees from true north."`
	El   float64 `json:"el"   description:"Elevation in degrees."`
	Ss   float64 `json:"ss"   description:"Signal strength in dB."`
	Used bool    `json:"used" description:"Used in current solution? (SBAS/WAAS/EGNOS satellites may be flagged used if the solution has corrections from them, but not all drivers make this information available.)
"`
}

// JSON struct for Horizon GPS data
type GpsData struct {
	Location   LocationData    `json:"location" description:"Location"`
	Satellites []SatelliteData `json:"satellites" description:"Data about discovered satellites."`
}

// The GpsCache "class"
type GpsCache struct {
	data  GpsData
	mutex sync.RWMutex
}

// GpsCache "constructor"
func New() *GpsCache {
	n := &GpsCache{}
	print(n)
	return n
}

// Set the source of the location data (Manual, Estimated, GPS, or Searching)
func (h *GpsCache) SetLocationSource(loc_source SourceType) {
	h.mutex.Lock()
	logutil.LogDebug("SetLocationSource(loc_source=%v)\n", loc_source)
	// The policy of not obfuscating if source is gps is enforced when env vars are read, so we do not need to do it here
	// if GPS == loc_source || SEARCHING == loc_source {
	// 	if 0.0 != h.data.Location.AccuracyKM {
	// 		logutil.Logf("Forcing accuracy_km to 0.0 for location source %q.\n", loc_source)
	// 	}
	// 	h.data.Location.AccuracyKM = 0.0
	// }
	h.data.Location.LocSource = loc_source
	print(h)
	h.mutex.Unlock()
	return
}

// Get the source of the location data (Manual, Estimated, GPS, or Searching)
func (h *GpsCache) GetLocationSource() (loc_source SourceType) {
	h.mutex.Lock()
	loc_source = h.data.Location.LocSource
	h.mutex.Unlock()
	return
}

// Set location source data, accuracy distance, and enforce policy constraints
func (h *GpsCache) SetLocationAccuracyInKm(accuracy_km float64) {
	h.mutex.Lock()
	logutil.LogDebug("SetLocationAccuracyInKm(accuracy_km=%f)\n", accuracy_km)
	// The policy of not obfuscating if source is gps is enforced when env vars are read, so we do not need to do it here
	// loc_source := h.data.Location.LocSource
	// if 0.0 != accuracy_km && (GPS == loc_source || SEARCHING == loc_source) {
	// 	err = errors.New(fmt.Sprintf("ERROR: Nonzero location accuracy (%f) is not permitted when location data source is %q.", accuracy_km, loc_source))
	// } else {
	h.data.Location.AccuracyKM = accuracy_km
	print(h)
	// 	err = nil
	// }
	h.mutex.Unlock()
	return
}

// Set satellite data
func (h *GpsCache) SetSatellites(report *gpsdc.SKYReport) {
	sats := report.Satellites
	h.mutex.Lock()
	if DEBUG > 0 {
		logutil.LogDebug("SetSatellites(%d satellites found):", len(sats))
		for _, sat := range sats {
			logutil.LogDebug("  PRN:%.1f, Az=%.1f, El=%.1f, Ss=%.1f, Used=%v", sat.PRN, sat.Az, sat.El, sat.Ss, sat.Used)
		}
	}

	// Clear the existing slice, then fill it from the gosdc data
	h.data.Satellites = nil
	for _, sat := range sats {
		var one SatelliteData
		one.PRN = sat.PRN
		one.Az = sat.Az
		one.El = sat.El
		one.Ss = sat.Ss
		one.Used = sat.Used
		h.data.Satellites = append(h.data.Satellites, one)
	}

	h.mutex.Unlock()
	print(h)
}

// Set location (latitude, longitude, elevation)
// NOTE: if location lat, lon, or elev are pure zero, the cache is not updated
// NOTE: if location provided matches the cached location, the LastUpdate
//       value is not changed.
// NOTE: Location randomization is done only when set, to prevent clients from
//       receiving many randomizations of a single actual location (enabling
//       them to estimate the actual location ever more accurately as more
//       randomizations are received.
// NOTE: See also the note above function obfuscate_lat_lon(), below.
func (h *GpsCache) SetLocation(lat float64, lon float64, elev float64) (update_time int64) {
	h.mutex.Lock()
	logutil.LogDebug("SetLocation(lat=%f, lon=%f, elev=%f)\n", lat, lon, elev)

	// Is the passed location nonzero? (zeros mean the GPS has no fix)
	if 0.0 != lat && 0.0 != lon && elev != 0.0 {

		// Is the passed location different from the cache?
		if lat != h.data.Location.Latitude || lon != h.data.Location.Longitude || elev != h.data.Location.ElevationM {

			// Location has changed.  Update the cache, and timestamp
			h.data.Location.Latitude, h.data.Location.Longitude = obfuscate_lat_lon(lat, lon, h.data.Location.AccuracyKM)
			if elev < 0.0 {
				elev = get_elevation_in_meters(lat, lon)
			}
			h.data.Location.ElevationM = elev

			// This is appropriate when loc_source==GPS, but when loc_source==Manual or Estimated we will update the timestamp when location is retrieved as json
			h.data.Location.LastUpdate = time.Now().UTC().Unix()
		}
	}
	h.mutex.Unlock()
	print(h)
	return h.data.Location.LastUpdate
}

// Return whether or not GPS hardware is available
func (h *GpsCache) HasGPS() (has_gps bool) {
	h.mutex.RLock()
	has_gps = GPS == h.data.Location.LocSource || SEARCHING == h.data.Location.LocSource
	h.mutex.RUnlock()
	return
}

// Return the current configuration details
func (h *GpsCache) GetConfiguration() (config string) {
	h.mutex.RLock()
	loc_source := h.data.Location.LocSource
	source := fmt.Sprintf("Source: %q", loc_source)
	accuracy_km := h.data.Location.AccuracyKM
	accuracy := fmt.Sprintf("obfuscation up to %f km", accuracy_km)
	if 0.0 == accuracy_km {
		accuracy = "no obfuscation"
	}
	config = fmt.Sprintf("%s -- %s", source, accuracy)
	h.mutex.RUnlock()
	return
}

// Return true if lat and lon have been set yet
func (h *GpsCache) IsLocationSet() (isSet bool) {
	h.mutex.RLock()
	isSet = (h.data.Location.Latitude != 0 && h.data.Location.Longitude != 0)
	h.mutex.RUnlock()
	return
}

// Return cached, and possibly already obfuscated, location (lat, lon, elev)
func (h *GpsCache) GetLocation() (lat float64, lon float64, elev float64, accuracy_km float64) {
	h.mutex.RLock()
	lat = h.data.Location.Latitude
	lon = h.data.Location.Longitude
	elev = h.data.Location.ElevationM
	accuracy_km = h.data.Location.AccuracyKM
	h.mutex.RUnlock()
	return
}

// Return cached location data as JSON-formatted bytes
func (h *GpsCache) GetLocationAsJSON() (json_bytes []byte) {
	h.mutex.RLock()
	if h.data.Location.LocSource == MANUAL || h.data.Location.LocSource == ESTIMATED {
		// copy the location struct and update the timestamp, otherwise this location could look really old
		loc := h.data.Location
		loc.LastUpdate = time.Now().UTC().Unix()
		json_bytes, _ = json.Marshal(loc)
	} else {
		json_bytes, _ = json.Marshal(h.data.Location)
	}
	h.mutex.RUnlock()
	return
}

// Return cached satellite data as JSON-formatted bytes
func (h *GpsCache) GetSatellitesAsJSON() (json_bytes []byte) {
	h.mutex.RLock()
	var sat_only SatellitesData
	sat_only.Satellites = h.data.Satellites
	json_bytes, _ = json.Marshal(sat_only)
	h.mutex.RUnlock()
	return
}

// Return the entire Horizon GPS data as a JSON-formatted bytes
func (h *GpsCache) GetAsJSON() (json_bytes []byte) {
	h.mutex.RLock()
	json_bytes, _ = json.Marshal(h.data)
	h.mutex.RUnlock()
	return
}

// Estimate node location using its public IP address for geo-location
func (h *GpsCache) EstimateLocation() (lat float64, lon float64, err error) {
	err, ip_address := get_public_address()
	for nil != err {
		logutil.Logf("%v\n", err)
		logutil.Logf("INFO: Pausing for %d seconds before retrying public address...\n", GET_PUBLIC_IP_RETRY_SEC)
		time.Sleep(time.Duration(GET_PUBLIC_IP_RETRY_SEC) * time.Second)
		err, ip_address = get_public_address()
	}
	logutil.Logf("INFO: Discovered IP address: %v\n", ip_address)
	lat, lon, err = get_location_from_ip_address(ip_address)
	for nil != err {
		logutil.Logf("%v\n", err)
		logutil.Logf("INFO: Pausing for %d seconds before retrying location estimate...\n", LOC_ESTIMATE_RETRY_SEC)
		time.Sleep(time.Duration(LOC_ESTIMATE_RETRY_SEC) * time.Second)
		lat, lon, err = get_location_from_ip_address(ip_address)
	}
	logutil.Logf("INFO: Estimated location: lat=%f, lon=%f\n", lat, lon)
	return
}

// Internal routine to dump the passed instance as json
// Please note that this routine expects that a lock is already held on the
// passed GpsCache object!
func print(h *GpsCache) {
	if DEBUG < 2 {
		return
	}
	j, _ := json.Marshal(h.data)
	var formatted bytes.Buffer
	_ = json.Indent(&formatted, j, "DEBUG: ", "  ")
	logutil.LogDebug("%s\n", string(formatted.Bytes()))
}

// Utility routine to obfuscate a location (latitude and longitude only)
// to some other location on the surface of the earth placed at some random
// bearing, and at some random distance along that bearing within the
// specified accuracy_km distance (in kilometers) from the actual location.
// Note that this obfuscation should only be done once for any location
// otherwise an adversary could receuve many obfuscated locations derived from
// a single original location and use a Monte Carlo technique to derive the
// original location: https://en.wikipedia.org/wiki/Monte_Carlo_integration
func obfuscate_lat_lon(lat float64, lon float64, accuracy_km float64) (obfuscated_lat float64, obfuscated_lon float64) {

	actual_location := geo.NewPoint(lat, lon)
	obfuscated_location := actual_location

	// Obfuscate within the specified accuracy (if nonzero)
	if 0.0 != accuracy_km {
		// Inefficiently create a new random number generator each time
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		random_distance_kilometers := accuracy_km * r.Float64()
		random_bearing_degrees := 360.0 * r.Float64()
		obfuscated_location = actual_location.PointAtDistanceAndBearing(random_distance_kilometers, random_bearing_degrees)
		if DEBUG > 2 {
			logutil.LogDebug("Obfuscating: From: (lat=%f,lon=%f), To: (lat=%f,lon=%f)\n", lat, lon, obfuscated_location.Lat(), obfuscated_location.Lng())
		}
	}

	return obfuscated_location.Lat(), obfuscated_location.Lng()
}

// Get the elevation in meters of a (lat, lon) surface location, using the
// free USGS Elevation Query Web Service (at http://ned.usgs.gov) which when
// invoked with a uri like this (x=lon, y=lat):
//     http://ned.usgs.gov/epqs/pqs.php?x=-121&y=37&units=Meters&output=json
// returns data in this form:
//     {"USGS_Elevation_Point_Query_Service":{"Elevation_Query":{"x":-121,"y":37,"Data_Source":"3DEP 1\/3 arc-second","Elevation":183.69,"Units":"Meters"}}}
func get_elevation_in_meters(lat float64, lon float64) (elevation_m float64) {
	// USGS returns -1000000 on error, so I'm co-opting that for all errors
	elevation_m = -1000000
	resp, get_err := http.Get(fmt.Sprintf("http://ned.usgs.gov/epqs/pqs.php?x=%f&y=%f&units=Meters&output=json", lon, lat))
	if get_err != nil {
		logutil.Logf("get_elevation_in_meters: error getting elevation from ned.usgs.gov: %v\n", get_err)
		return
	}
	defer resp.Body.Close()
	body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		logutil.Logf("get_elevation_in_meters: error reading response body from ned.usgs.gov: %v\n", read_err)
		return
	}
	// The output from ned.usgs.gov is like: {"USGS_Elevation_Point_Query_Service":{"Elevation_Query":{"x":-73.959494,"y":42.214607,"Data_Source":"3DEP 1\/3 arc-second","Elevation":139.13,"Units":"Meters"}}}
	var parsed map[string]map[string]map[string]interface{}
	unm_err := json.Unmarshal(body, &parsed)
	if unm_err != nil {
		logutil.Logf("get_elevation_in_meters: error parsing response from ned.usgs.gov: %v\n", unm_err)
		return
	}
	data := parsed["USGS_Elevation_Point_Query_Service"]["Elevation_Query"]["Elevation"]
	e, ok := data.(float64)
	if !ok {
		logutil.Logf("get_elevation_in_meters: error converting elevation to float64, elevation=%v\n", data)
		return
	}
	elevation_m = e
	if DEBUG > 2 {
		logutil.LogDebug("Discovered altitude for (lat=%f,lon=%f): %f meters\n", lat, lon, elevation_m)
	}
	return
}

// Get the public Internet IP address of this node (outside any firewalled
// and NATted LAN).  Typically the node wil not be reachable at this IP
// address due to the firewall, but it is convenient to use the public IP
// address to estimate location when no other source of location data is
// available.  This function is hard-coded to use "http://ifconfig.co".
// It will need to be rewritten if we change to using a different provider.
func get_public_address() (err error, ip_address string) {
	err = nil
	// NOTE: Using 'https' sometimes fails with a 509 due to
	// certificate errors.  So using 'http' is more reliable.
	provider_url := "http://ifconfig.co"
	resp, rest_err := http.Get(provider_url)
	if rest_err != nil {
		ip_address = ""
		err = errors.New(fmt.Sprintf("ERROR: REST call to %s failed: %v", provider_url, rest_err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// Parse out the IP address from the body
	ip_address = string(body[:len(body)-1])
	return
}

type IPGPSCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	// there are many more fields, but we don't care about them
}

// Estimate the location of the node by using its IP address.  This function
// uses the FreeGeoIP service, which enables up to 10,000 queries per hour
// (after which it returns HTTP 403 responses).  See "http://freegeoip.net/".
func get_location_from_ip_address(ip_address string) (lat float64, lon float64, err error) {
	apikey := "166f0d44636d9ed71c6a01530f687fdb"
	url := "http://api.ipstack.com/" + ip_address + "?access_key=" + apikey + "&format=1"
	resp, rest_err := http.Get(url)
	if rest_err != nil {
		err = errors.New(fmt.Sprintf("ERROR: Request to %s failed: %v", url, rest_err))
		return
	}
	defer resp.Body.Close()
	httpCode := resp.StatusCode
	if httpCode != 200 {
		err = errors.New(fmt.Sprintf("ERROR: Bad http code %d from %s\n", httpCode, url))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if DEBUG > 2 {
		logutil.LogDebug("Response from %s:\n%s\n", url, string(body))
	}
	coordinates := IPGPSCoordinates{}
	err = json.Unmarshal(body, &coordinates)
	if err != nil {
		err = errors.New(fmt.Sprintf("ERROR: failed to unmarshal body response from %s: %v", url, err))
		return
	}
	lat = coordinates.Latitude
	lon = coordinates.Longitude
	return
}
