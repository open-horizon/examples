# SDR POC Database

The Watson insights from SDR audio analysis are stored an IBM Compose Postgresql DB.

## Schema (in progress)

- GlobalNouns: summary of the nouns and sentiments from all nodes and stations
    - noun (key, string)
    - sentiment (float64) - sentiment score from -1.0 (full negative) to 1.0 (full positive). A running average of all sentiments we've received for this noun from all nodes/stations.
    - numberOfMentions (integer) - running count
    - timeUpdated (date) - most recent update
- NodeNouns: summary of the nouns and sentiments from this node (all stations)
    - noun (key, string)
    - edgeNode (key, foreign key, string) - horizon org/nodeid
    - sentiment (float64) - sentiment score from -1.0 (full negative) to 1.0 (full positive). A running average of all sentiments we've received from this node (for all stations).
    - numberOfMentions (integer) - running count
    - timeUpdated (date) - most recent update
- Stations: list of the nodes and their stations that we've received audio clips from
    - edgeNode (key, string) - horizon org/nodeid that received data from this station
    - frequency (key, float32) - station frequency
    - numberofclips (int) - number of audio clips received from this station from this edge node
    - dataQualityMetric (float32) - as determined/reported by the data_broker service
    - Phase 2:
        - callLetters (when broadcast)
        - (maybe some other way to determine that stations from nearby nodes are actually the same)
- EdgeNodes:
    - edgeNode (key, string) - horizon org/nodeid
    - latitude (float32) - latitude of the node
    - longitude (float32) - longitude of the node

## Manually Connect to the DB
```
export SDR_DB_PASSWORD='<pw>'
export SDR_DB_USER=admin
export SDR_DB_HOST='<host>'
export SDR_DB_PORT=<port>
export SDR_DB_NAME=sdr

psql "sslmode=require host=$SDR_DB_HOST port=$SDR_DB_PORT dbname=$SDR_DB_NAME user=$SDR_DB_USER password=$SDR_DB_PASSWORD"
```

## Example SQL Statements to Manually Create/Modify Tables
```
CREATE TABLE globalnouns(noun TEXT PRIMARY KEY NOT NULL, sentiment DOUBLE PRECISION NOT NULL, numberofmentions BIGINT NOT NULL, timeupdated timestamp with time zone);

CREATE TABLE nodenouns(noun TEXT NOT NULL, edgenode TEXT NOT NULL, sentiment DOUBLE PRECISION NOT NULL, numberofmentions BIGINT NOT NULL, timeupdated timestamp with time zone, PRIMARY KEY(noun, edgenode) );

CREATE TABLE stations(edgenode TEXT NOT NULL, frequency REAL NOT NULL, numberofclips BIGINT NOT NULL, dataqualitymetric REAL, timeupdated timestamp with time zone, PRIMARY KEY(edgenode, frequency) );

CREATE TABLE edgenodes(edgenode TEXT PRIMARY KEY NOT NULL, latitude REAL NOT NULL, longitude REAL NOT NULL, timeupdated timestamp with time zone);

# Add rows to the globalnouns table:
INSERT INTO globalnouns VALUES ('wedding', 0.99, 2, '2018-06-23 10:05:00');
INSERT INTO globalnouns VALUES ('trump', -0.25, 100, '2018-07-23 11:05:00');
INSERT INTO globalnouns VALUES ('foo', 0.0, 100, '2018-08-01 11:05:00');

# Update a row:
UPDATE globalnouns SET sentiment = 0.25, timeupdated = '2018-06-23 14:00' WHERE noun = 'wedding';

# Upsert a row (insert if not there, update if there):
INSERT INTO globalnouns VALUES ('wedding', 0.5, 1, CURRENT_TIMESTAMP) ON CONFLICT (noun) DO UPDATE SET sentiment = ((globalnouns.sentiment * globalnouns.numberofmentions) + 0.5) / (globalnouns.numberofmentions + 1), numberofmentions = globalnouns.numberofmentions + 1, timeupdated = CURRENT_TIMESTAMP;

# If you need to change a column definition:
ALTER TABLE globalnouns alter column timeupdated type timestamp with time zone;
```

## Run Example Go Code to Write and Read DB
```
# set same env vars as above
go get github.com/lib/pq
make
./sdr-db
```
