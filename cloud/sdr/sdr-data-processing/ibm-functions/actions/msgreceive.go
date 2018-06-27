package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "io/ioutil"
    "encoding/json"
    "encoding/base64"
    "net/http"
)

const (
    VERBOSE = true
)

type ActionMsg struct {
    Value string `json:"value"`
    Key string `json:"key"`
    Topic string `json:"topic"`
    Partition int `json:"partition"`
    Offset int `json:"offset"`
}

type ActionArg struct {
    Messages []ActionMsg `json:"messages"`
    WatsonSttUsername string `json:"watsonSttUsername"`
    WatsonSttPassword string `json:"watsonSttPassword"`
}

type ModelList struct {
    Models []map[string]interface{}
}

func main() {
    //program receives one argument: the JSON object as a string
    arg := os.Args[1]

    // can optionally log to stdout (or stderr)
    fmt.Printf("msgreceive arg: %s\n", arg)

    // unmarshal the string to a JSON object
    // var obj map[string]interface{}
    obj := ActionArg{}
    err := json.Unmarshal([]byte(arg), &obj)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error unmarshalling the arg: %v\n", err)
        os.Exit(2)
    }

    // Echo the messages we received
    fmt.Println("messages:")
    for _, m := range obj.Messages {
        var msg string
        var err error
        if msg, err = strconv.Unquote(m.Value); err != nil {     // for some reason the msg comes with escaped double quotes around it
            fmt.Fprintf(os.Stderr, "Message '%s' not escaped, using original message. Err: %v\n", m.Value, err)
            msg = m.Value
        }
        // else msg has the unescaped value from Unquote()
        fmt.Printf("Msg from topic %s, partition %d: %s\n", m.Topic, m.Partition, msg)
    }

    var models ModelList
    actionResult := make(map[string]string)
    httpCode := HttpGet("https://stream.watsonplatform.net/speech-to-text/api/v1/models", obj.WatsonSttUsername+":"+obj.WatsonSttPassword, &models)
    if httpCode == 200 && len(models.Models) > 0 {
        fmt.Printf("First Models: %v\n", models.Models[0])

        // last line of stdout is the result JSON object as a string
        actionResult["msg"] = "msgreceive successful"
    } else {
        actionResult["msg"] = "could not get models from Watson STT"
    }

    res, _ := json.Marshal(actionResult)
    fmt.Println(string(res))
}

func HttpGet(url string, credentials string, structure interface{}) (httpCode int) {
    apiMsg := http.MethodGet + " " + url
    Verbose(apiMsg)
    httpClient := &http.Client{}
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "New request failed for GET %s: %v\n", apiMsg, err)
    }
    req.Header.Add("Accept", "application/json")
    if credentials != "" {
        req.Header.Add("Authorization", fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(credentials))))
    }
    resp, err := httpClient.Do(req)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error running GET %s: %v\n", apiMsg, err)
    }
    defer resp.Body.Close()
    httpCode = resp.StatusCode
    Verbose("HTTP code: %d", httpCode)
    if httpCode != 200 {
        fmt.Fprintf(os.Stderr, "Error: bad HTTP code from %s: %d", apiMsg, httpCode)
    } else {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Fprintf(os.Stderr, "failed to read body response from %s: %v", apiMsg, err)
        }
        switch s := structure.(type) {
        case *string:
            // Just return the unprocessed response body
            *s = string(bodyBytes)
        default:
            // Put the response body in the specified struct
            err = json.Unmarshal(bodyBytes, structure)
            if err != nil {
                fmt.Fprintf(os.Stderr, "failed to unmarshal body response from %s: %v", apiMsg, err)
            }
        }
    }
    return
}


func Verbose(msg string, args ...interface{}) {
    if !VERBOSE {
        return
    }
    if !strings.HasSuffix(msg, "\n") {
        msg += "\n"
    }
    //fmt.Fprintf(os.Stderr, "[verbose] "+msg, args...) // send to stderr so it doesn't mess up stdout if they are piping that to jq or something like that
    fmt.Printf("[verbose] "+msg, args...)
}
