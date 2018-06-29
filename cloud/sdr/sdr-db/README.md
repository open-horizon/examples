# SDR POC Database

The Watson insights from SDR audio analysis are stored an IBM Compose Postgresql DB.

## Schema (in progress)

- Nouns: noun (key), localStation (key and foreign key), edgeNodeId (key and foreign key), sentiment, numberOfMentions, timeUpdated
- LocalStations: edgeNodeId (key), frequency, dataQualityMetric
- EdgeNode: id (key), latitude, longitude
- Phase 2: Stations: callLetters, localStation

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

UPDATE nouns SET timeupdated = '2018-06-23 14:00' WHERE noun = 'wedding';
```

## Run Example Go Code to Write and Read DB
```
# set same env vars as above
go get github.com/lib/pq
make
./sdr-db
```
