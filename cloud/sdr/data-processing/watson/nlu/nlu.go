package nlu

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/open-horizon/examples/cloud/sdr/data-processing/wutil"
)

// SentimentRequest is the input body for the NLU POST /analyze REST API
type SentimentRequest struct {
	Text     string           `json:"text"`
	Features SentimentFeature `json:"features"`
}

// SentimentFeature is what we want it to do and return
type SentimentFeature struct {
	Entities SentimentEntities `json:"entities"`
}

// SentimentEntities will recognize nouns and do and return sentiment for each
type SentimentEntities struct {
	Sentiment bool `json:"sentiment"`
}

// Sentiment uses Watson NLU to determine sentiments of recognizable nouns
func Sentiment(text, username, password string) (sentiments string, err error) {
	fmt.Println("using Watson NLU to get sentiments from text...")
	// Watson NLU API: https://www.ibm.com/watson/developercloud/natural-language-understanding/api/v1/#post-analyze
	apiURL := "https://gateway.watsonplatform.net/natural-language-understanding/api/v1/analyze"
	headers := []wutil.Header{{Key: "Content-Type", Value: "application/json"}}
	sentimentReq := SentimentRequest{Text: text, Features: SentimentFeature{Entities: SentimentEntities{Sentiment: true}}}
	json, err := json.Marshal(sentimentReq)
	if err != nil {
		panic(err)
	}

	err = wutil.HTTPPost(apiURL, username, password, headers, bytes.NewReader(json), &sentiments)
	return
}
