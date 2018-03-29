# Horizon Location Workload

Reports the GPS coordinates of the Horizon edge node and satellite information to the data ingest service. This workload by the other Horizon POCs to add location info to the data an edge node is producing.

## Environment Variables

These environment variables can be set in the location workload policy file deployment string to override the default behavior:

| Name | Description |
| ---- | ---------------- |
| REPORTING_INTERVAL | Number of seconds to wait between successive publications to the central MQTT |
| SKIP_NUM_REPEAT_LOC_READINGS | If the GPS data continues to be the same, skip this many readings before sending to MQTT |
| SKIP_NUM_REPEAT_SAT_READINGS | If the statellite data continues to be the same, skip this many readings before sending to MQTT |
| MAX_REGISTRATION_ATTEMPTS | Number of times to try registering with the central MQTT before giving up |
| SECONDS_BETWEEN_REG_ATTEMPTS | Number of seconds to wait between registration attempts |
| REG_SECONDS_BEFORE_STREAMING | Number of seconds to wait between successful registration and streaming data to MQTT |

## Data Sent to Data Ingest Service (MQTT)

The GPS data JSON structure sent to topic `/applications/in/<agreement-id>/public/h/<device-id>/0/`:

```
{
  "t": 1498250140,
  "r": {
    "lat": 42.052304,
    "lon": -73.9601765,
    "alt": 93.3
  }
}
```

The statellite data JSON structure sent to topic `/applications/in/<agreement-id>/public/h/<device-id>/6/`:

```
{
  "t": 1498250140,
  "d": [
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
