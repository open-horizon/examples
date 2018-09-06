// Example for consuming  messages from IBM Cloud Message Hub (kafka) using go.
// See README.md for setup requirements.

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"database/sql"

	"github.com/Shopify/sarama"             // doc: https://godoc.org/github.com/Shopify/sarama
	cluster "github.com/bsm/sarama-cluster" // doc: http://godoc.org/github.com/bsm/sarama-cluster
	_ "github.com/lib/pq"

	"github.com/open-horizon/examples/cloud/sdr/data-ingest/example-go-clients/util"
	"github.com/open-horizon/examples/cloud/sdr/data-processing/watson/nlu"
	"github.com/open-horizon/examples/cloud/sdr/data-processing/watson/stt"
	"github.com/open-horizon/examples/cloud/sdr/data-processing/wutil"
	"github.com/open-horizon/examples/edge/msghub/sdr2msghub/audiolib"
)

func usage(exitCode int) {
	fmt.Printf("Usage: %s [-t <topic>] [-h] [-v]\n\nEnvironment Variables: MSGHUB_API_KEY, MSGHUB_BROKER_URL, MSGHUB_TOPIC\n", os.Args[0])
	os.Exit(exitCode)
}

const minConfidence = 0.5

func main() {
	sttUsername := util.RequiredEnvVar("STT_USERNAME", "")
	sttPassword := util.RequiredEnvVar("STT_PASSWORD", "")

	nluUsername := util.RequiredEnvVar("NLU_USERNAME", "")
	nluPassword := util.RequiredEnvVar("NLU_PASSWORD", "")

	// Get all of the input options
	var topic string
	flag.StringVar(&topic, "t", "", "topic")
	var help bool
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&util.VerboseBool, "v", false, "verbose")
	flag.Parse()
	if help {
		usage(1)
	}

	db := connect2DB() // this will read necessary env vars
	defer db.Close()

	// Get msg hub info
	apiKey := util.RequiredEnvVar("MSGHUB_API_KEY", "")
	username := apiKey[:16]
	password := apiKey[16:]
	util.Verbose("username: %s, password: %s\n", username, password)
	brokerStr := util.RequiredEnvVar("MSGHUB_BROKER_URL", "kafka01-prod02.messagehub.services.us-south.bluemix.net:9093,kafka02-prod02.messagehub.services.us-south.bluemix.net:9093,kafka03-prod02.messagehub.services.us-south.bluemix.net:9093,kafka04-prod02.messagehub.services.us-south.bluemix.net:9093,kafka05-prod02.messagehub.services.us-south.bluemix.net:9093")
	brokers := strings.Split(brokerStr, ",")
	if topic == "" {
		topic = util.RequiredEnvVar("MSGHUB_TOPIC", "sdr-audio")
	}

	util.Verbose("starting message hub consuming example...")
	if util.VerboseBool {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	// Initialize msg hub consumer
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	err := util.PopulateConfig(&config.Config, username, password, apiKey) // add creds and tls info
	util.ExitOnErr(err)
	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", []string{topic}, config)
	util.ExitOnErr(err)
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			if util.VerboseBool {
				log.Printf("Rebalanced: %+v\n", ntf)
			}
		}
	}()

	// consume messages, watch signals
	for {
		fmt.Printf("Listening for messages produced to %s...\n", topic)
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				// Got an audio clip, convert it to text
				audioMsg := &audiolib.AudioMsg{}
				err = json.Unmarshal(msg.Value, &audioMsg)
				if err != nil {
					log.Println(err)
					continue
				}
				timeStamp := time.Unix(audioMsg.Ts, 0)
				audio, err := base64.StdEncoding.DecodeString(audioMsg.Audio)
				if err != nil {
					log.Println(err)
					continue
				}
				fmt.Println("got audio from device:", audioMsg.DevID, "on station:", audioMsg.Freq)
				transcript, err := stt.Transcribe(audio, sttUsername, sttPassword)
				fatalIfErr(err)
				if util.VerboseBool {
					fmt.Println("STT:", wutil.MarshalIndent(transcript.Results))
				} else {
					fmt.Println("STT:", transcript.Results)
				}

				//todo: make edgenode reference in nodenouns and stations a foreign key?

				// Send each string of text that has good confidence to NLU
				for _, r := range transcript.Results {
					altNum := 0 //todo: we only seem to get 1 alternative, not sure if it will always be that way
					if r.Final && r.Alternatives[altNum].Confidence > minConfidence {
						sentiments, err := nlu.Sentiment(r.Alternatives[altNum].Transcript, nluUsername, nluPassword)
						if err != nil {
							fmt.Println(err)
							continue
						}
						if util.VerboseBool {
							fmt.Println("NLU:", wutil.MarshalIndent(sentiments))
						} else {
							fmt.Println("NLU:", sentiments)
						}
						addSentimentsToDB(db, sentiments, timeStamp, audioMsg.DevID)
					} else {
						util.Verbose("Skipping: Final: %v, Confidence: %f, Text: %s\n", r.Final, r.Alternatives[altNum].Confidence, r.Alternatives[altNum].Transcript)
					}
				}

				// Record node and station info in db
				addNodeStationToDB(db, timeStamp, audioMsg.DevID, audioMsg.Freq, audioMsg.Lat, audioMsg.Lon, audioMsg.ExpectedValue)

				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}

func fatalIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// addSentimentsToDB adds the nouns, sentiments, and node id to DB tables
func addSentimentsToDB(db *sql.DB, sentiments nlu.AnalyzeResponse, timeStamp time.Time, nodeID string) {
	fmt.Println("adding the nouns, sentiments, and node id to DB tables...")
	dups := make(map[string]bool)                                 // so we can avoid duplicates from entities and keywords
	nouns := append(sentiments.Entities, sentiments.Keywords...)  // concat the 2 lists
	ts := timeStamp.Format("2006-01-02 15:04:05.999999999 -0700") // Time.String() adds text for the time zone (e.g. EDT), which postgres doesn't accept

	// Loop thru the entities and keywords, adding their nouns and sentiments to the DB
	for _, e := range nouns {
		noun := e.Text
		sentiment := e.Sentiment.Score
		if _, ok := dups[noun]; ok {
			continue // we already processed this noun
		}
		dups[noun] = true
		util.Verbose("adding noun %s with sentiment score %f to db...", noun, sentiment)

		// Add this noun/sentiment to the globalnouns table
		// This is the postgres way to upsert a row (insert if not there, update if there)
		stmt, err := db.Prepare("INSERT INTO globalnouns VALUES ($1, $2, 1, $3) ON CONFLICT (noun) DO UPDATE SET sentiment = ((globalnouns.sentiment * globalnouns.numberofmentions) + $2) / (globalnouns.numberofmentions + 1), numberofmentions = globalnouns.numberofmentions + 1, timeupdated = $3")
		fatalIfErr(err)
		defer stmt.Close()
		_, err = stmt.Exec(noun, sentiment, ts)
		fatalIfErr(err)
		// affect, err := res.RowsAffected()
		// fatalIfErr(err)

		// Add this noun/sentiment to the nodenouns table
		stmt, err = db.Prepare("INSERT INTO nodenouns VALUES ($1, $4, $2, 1, $3) ON CONFLICT ON CONSTRAINT nodenouns_pkey DO UPDATE SET sentiment = ((nodenouns.sentiment * nodenouns.numberofmentions) + $2) / (nodenouns.numberofmentions + 1), numberofmentions = nodenouns.numberofmentions + 1, timeupdated = $3")
		fatalIfErr(err)
		defer stmt.Close() // defer will evaluate and save the new value of stmt: https://golang.org/ref/spec#Defer_statements
		_, err = stmt.Exec(noun, sentiment, ts, nodeID)
		fatalIfErr(err)
	}
}

// addNodeStationToDB adds the node and station info to DB tables. This is only called once per msg hub msg.
func addNodeStationToDB(db *sql.DB, timeStamp time.Time, nodeID string, stationFreq, latitude, longitude, expectedValue float32) {
	fmt.Println("adding the node and station info to DB tables...")
	ts := timeStamp.Format("2006-01-02 15:04:05.999999999 -0700") // Time.String() adds text for the time zone (e.g. EDT), which postgres doesn't accept

	// Add station info to the stations table
	stmt, err := db.Prepare("INSERT INTO stations VALUES ($1, $2, 1, $3, $4) ON CONFLICT ON CONSTRAINT stations_pkey DO UPDATE SET numberofclips = stations.numberofclips + 1, dataqualitymetric =$3, timeupdated = $4")
	fatalIfErr(err)
	defer stmt.Close()
	_, err = stmt.Exec(nodeID, stationFreq, expectedValue, ts) //todo: not sure what to do with expectedValue
	fatalIfErr(err)

	// Add node info to the edgenodes table
	stmt, err = db.Prepare("INSERT INTO edgenodes VALUES ($1, $2, $3, $4) ON CONFLICT (edgenode) DO UPDATE SET latitude = $2, longitude = $3, timeupdated = $4")
	fatalIfErr(err)
	defer stmt.Close() // defer will evaluate and save the new value of stmt: https://golang.org/ref/spec#Defer_statements
	_, err = stmt.Exec(nodeID, latitude, longitude, ts)
	fatalIfErr(err)
}

func connect2DB() *sql.DB {
	// The lib/pq package is a postgresql driver for the standard database/sql package. See https://godoc.org/github.com/lib/pq for details.

	// Get the env var values we need to connect to our db
	pw := util.RequiredEnvVar("SDR_DB_PASSWORD", "")
	host := util.RequiredEnvVar("SDR_DB_HOST", "")
	port := util.RequiredEnvVar("SDR_DB_PORT", "")
	user := util.RequiredEnvVar("SDR_DB_USER", "admin")
	dbName := util.RequiredEnvVar("SDR_DB_NAME", "sdr")

	// Connect to db
	connStr := "postgres://" + user + ":" + pw + "@" + host + ":" + port + "/" + dbName + "?sslmode=require"
	db, err := sql.Open("postgres", connStr)
	fatalIfErr(err)
	fmt.Printf("connected to %s\n", connStr)
	return db
}
