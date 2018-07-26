# Horizon GPS Service

The GPS service provides location coordinates and satellite data to other Horizon services.

## Input Values

These values can be passed to the GPS service thru the `global` section of the input file given to `hzn register`:
```
  "global": [
    {
      "type": "LocationAttributes",
      "variables": {
        "lat": 45.421530,     /* this is passed to each container as HZN_LAT */
        "lon": -75.697193,    /* this is passed to each container as HZN_LON */
        "use_gps": false,    /* true if you have, and want to use, an attached GPS sensor. Passed to each container as HZN_USE_GPS. */
        "location_accuracy_km": 0.0   /* Make the node location inaccurate by this number of KM to protect privacy. */
      }
    }
  ],
```

If you don't specify the values above, the GPS service will default to using the edge node's public IP to estimate its location.

## RESTful API

Other Horizon services can use the GPS service by requiring it in its own service definition, and then in its code accessing the GPS REST APIs with the URL:
```
http://gps:31779/v1/<api-from-the-list-below>
```

### **API:** GET /gps/location
---

#### Parameters:
none

#### Response:

code: 
* 200 -- success
* other http codes TBD

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
* other http codes TBD

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
