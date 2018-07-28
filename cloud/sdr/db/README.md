# SDR POC Database

The Watson insights from SDR audio analysis are stored an IBM Compose Postgresql DB.

## Schema (in progress)

- Nouns:
    - noun (key, string)
    - edgeNode (key, foreign key, string) - horizon org/nodeid
    - frequency (key, foreign key, float) - station frequency
    - sentiment (float) - sentiment score from -1.0 (full negative) to 1.0 (full positive). A running average of all sentiments we've received from this node and station.
    - numberOfMentions (integer) - running count
    - timeUpdated (date) - most recent update
- Stations:
    - edgeNode (key, string) - horizon org/nodeid that received data from this station
    - frequency (key, float) - station frequency
    - dataQualityMetric (float) - as determined/reported by the data_broker service
    - Phase 2:
        - callLetters (when broadcast)
        - (maybe some other way to determine that stations from nearby nodes are actually the same)
- EdgeNodes:
    - edgeNode (key, string) - horizon org/nodeid
    - latitude (float) - latitude of the node
    - longitude (float) - longitude of the node

## Manually Connect to the DB
```
export SDR_DB_PASSWORD='<pw>'
export SDR_DB_USER=admin
export SDR_DB_HOST='<host>'
export SDR_DB_PORT=<port>
export SDR_DB_NAME=sdr

psql "sslmode=require host=$SDR_DB_HOST port=$SDR_DB_PORT dbname=$SDR_DB_NAME user=$SDR_DB_USER password=$SDR_DB_PASSWORD"
```

## Example Commands to Create/Modify Tables
```
CREATE TABLE nouns(
   noun TEXT PRIMARY KEY NOT NULL,
   sentiment TEXT NOT NULL,
   numberofmentions INT NOT NULL,
   timeupdated timestamp with time zone
);

INSERT INTO nouns VALUES ('wedding', 'positive', 2, '2018-06-23 10:05:00');
INSERT INTO nouns VALUES ('trump', 'negative', 100, '2018-07-23 11:05:00');
INSERT INTO nouns VALUES ('foo', 'positive', 100, '2018-07-24 11:05:00');

# If you need to manually update a row or change a column definition:
UPDATE nouns SET timeupdated = '2018-06-23 14:00' WHERE noun = 'wedding';
ALTER TABLE nouns alter column timeupdated type timestamp with time zone;
```

## Run Example Go Code to Write and Read DB
```
# set same env vars as above
go get github.com/lib/pq
make
./sdr-db
```
