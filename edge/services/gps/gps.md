# Horizon GPS Microservice

The GPS microservice provides location coordinates and satellite data to Horizon workload clients.

## Environment Variables

These environment variables are used by the GPS microservice container (and are passed by Horizon):

* **HZN_USE_GPS==true**: the device owner will allow exact gps coords from gps sensor to be shared
    - **HZN_LAT**, **HZN_LON** can be optionally set as fallback when the gps sensor is not available
    - **HZN_LOCATION_ACCURACY_KM** is ignored (because sharing exact coords was approved)
    - **HZN_USER_PROVIDED_COORDS** is also ignored (because we do not care if coords are estimated or not, because we are not applying an accuracy)
* **HZN_USE_GPS==false** (default): do not share exact gps coords from gps sensor
    - **HZN_LAT**, **HZN_LON** must be set (either entered by device owner, or estimated from IP)
    - **HZN_USER_PROVIDED_COORDS==true** (default): device owner entered lat/lon
        - **HZN_LOCATION_ACCURACY_KM**: obfuscate lat/lon by this much
    - **HZN_USER_PROVIDED_COORDS==false**: lat/lon was estimated from IP
        - **HZN_LOCATION_ACCURACY_KM** is ignored (forced to 0) because IP estimate is already inaccurate

## RESTful API

### **API:** GET /gps/location
---

#### Parameters:
none

#### Response:

code: 
* 200 -- success

body:


| Name | Type | Description |
| ---- | ---- | ---------------- |
| latitude | float | the latitude of the current location |
| longitude | float | the longitude of the current location |
| elevation | float | the elevation of the current location in meters |
| accuracy_km | float | the location accuracy in kilometers |
| loc_source | string | one of: Manual, Estimated, GPS, or Searching |
| loc_last_update | float | the timestamp the location was read (UTC) |


#### Example:
```
curl -sS -w "%{http_code}" http://gps:31779/v1/gps/location | jq .
{
  "latitude": 42.052577333,
  "longitude": -73.960314,
  "elevation": 33,
  "accuracy_km": 0,
  "loc_source": "GPS",
  "loc_last_update": 1498070507
}
200
```

### **API:** GET /gps/satellites
---

#### Parameters:
none

#### Response:

code: 
* 200 -- success

body:


| Name | Type | Description |
| ---- | ---- | ---------------- |
| satellites | json | array of data for satellites |
| satellites.PRN | int | PRN ID of the satellite. 1-63 are GNSS satellites, 64-96 are GLONASS satellites, 100-164 are SBAS satellites |
| satellites.az | int | azimuth, degrees from true north |
| satellites.el | int | elevation in degrees |
| satellites.ss | int | signal strength in dB |
| satellites.used | int | used in current location solution? (SBAS/WAAS/EGNOS satellites may be flagged used if the solution has corrections from them, but not all drivers make this information available) |


#### Example:
```
curl -sS -w "%{http_code}" http://gps:31779/v1/gps/satellites | jq .
{
  "satellites": [
    {
      "PRN": 1,
      "az": 63,
      "el": 58,
      "ss": 23,
      "used": true
    },
    {
      "PRN": 3,
      "az": 125,
      "el": 20,
      "ss": 27,
      "used": true
    },
    ...
  ]
}
```
