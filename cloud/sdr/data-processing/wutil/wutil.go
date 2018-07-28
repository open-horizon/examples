package wutil

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Header is 1 http header that should be added to the http call
type Header struct {
	Key   string
	Value string
}

// HTTPPost makes a REST post request and returns the response
func HTTPPost(url, username, password string, headers []Header, requestBody io.Reader, response interface{}) (err error) {
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, url, requestBody)
	r.SetBasicAuth(username, password)
	for _, h := range headers {
		r.Header.Add(h.Key, h.Value)
	}
	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.New("rest call to " + url + " got status " + strconv.Itoa(resp.StatusCode))
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if len(bytes) > 0 && response != nil {
		switch r := response.(type) {
		case *string:
			// They gave us a string, so just returned the unprocessed response
			*r = string(bytes)
		default:
			// They gave a struct, so put the json in it
			if err = json.Unmarshal(bytes, response); err != nil {
				return
			}
		}
	}
	return
}

// MarshallIndent converts the struct to json and then to a string, for quick output.
func MarshalIndent(myStruct interface{}) string {
	json, err := json.MarshalIndent(myStruct, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(json)
}
