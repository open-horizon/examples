package nlu

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/open-horizon/examples/cloud/sdr/data-ingest/example-go-clients/util"
	"github.com/open-horizon/examples/cloud/sdr/data-processing/wutil"
)

// AnalyzeRequest is the input body for the NLU POST /analyze REST API
type AnalyzeRequest struct {
	Text     string        `json:"text"`
	Features EntityFeature `json:"features"`
}

// EntityFeature is what we want it to do and return
type EntityFeature struct {
	Entities SentimentOption `json:"entities"` // proper nouns
	Keywords SentimentOption `json:"keywords"` //todo: other keywords. This gives us more results, but we have to eliminate duplicates from what was returned in the entities list
}

// SentimentOption will recognize nouns and return sentiment for each
type SentimentOption struct {
	Sentiment bool `json:"sentiment"`
}

// AnalyzeResponse is the response from the NLU POST /analyze REST API
type AnalyzeResponse struct {
	Usage    UsageDetails     `json:"usage"`
	Entities []EntityResponse `json:"entities"`
	Keywords []EntityResponse `json:"keywords"` // this is a subset of EntityResponse so we can reuse the struct
}

// UsageDetails summarizes how many watson service resources we used
type UsageDetails struct {
	TextUnits      int `json:"text_units"`
	TextCharacters int `json:"text_characters"`
	Features       int `json:"features"`
}

// EntityResponse is a proper noun that was recognized
type EntityResponse struct {
	Type      string            `json:"type"`
	Text      string            `json:"text"` // the proper noun
	Sentiment SentimentResponse `json:"sentiment"`
	Relevance float64           `json:"relevance"` // Value between 0 and 1. Doc didn't help understand exactly what this means
	Count     int               `json:"count"`
}

// SentimentResponse holds the sentiment and score
type SentimentResponse struct {
	Score float64 `json:"score"` // how strong this sentiment is. Value seems to be between -1.0 and 1.0
	Label string  `json:"label"` // this is the sentiment: positive, neutral, or negative
}

// Sentiment uses Watson NLU to determine sentiments of recognizable nouns
func Sentiment(text, username, password string) (sentiments AnalyzeResponse, err error) {
	fmt.Println("using Watson NLU to get sentiments from text...")
	// util.Verbose("text: %s", text)
	// Watson NLU API: https://www.ibm.com/watson/developercloud/natural-language-understanding/api/v1/#post-analyze
	analyzeVersion := "version=2018-03-19" //todo: this is required, but not sure how often this will change
	apiURL := "https://gateway.watsonplatform.net/natural-language-understanding/api/v1/analyze?" + analyzeVersion
	headers := []wutil.Header{{Key: "Content-Type", Value: "application/json"}}
	sentimentReq := AnalyzeRequest{Text: text, Features: EntityFeature{Entities: SentimentOption{Sentiment: true}, Keywords: SentimentOption{Sentiment: true}}}
	json, err := json.Marshal(sentimentReq)
	if err != nil {
		panic(err)
	}
	util.Verbose("request body: %s", string(json))

	err = wutil.HTTPPost(apiURL, username, password, headers, bytes.NewReader(json), &sentiments)
	return
}
