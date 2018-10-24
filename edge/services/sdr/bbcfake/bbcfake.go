package bbcfake

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	mp3 "github.com/hajimehoshi/go-mp3"
)

func downsample(input []byte) (out []byte) {
	out = make([]byte, len(input)/6+1)
	for i := 0; i < len(input); i += 12 {
		out[i/6] = input[i]
		out[(i/6)+1] = input[i+1]
	}
	return
}

// DownloadAndSplit fetches audio from a URL and returns a slice of chunks.
func DownloadAndSplit(url string) (chunks [][]byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	d, err := mp3.NewDecoder(resp.Body)
	if err != nil {
		return
	}
	defer d.Close()

	var buf = make([]byte, 512)

Loop:
	for {
		var chunk bytes.Buffer
		var writer = bufio.NewWriter(&chunk)
		for i := 0; i < 10997; i++ {
			_, err = d.Read(buf)
			if err != nil {
				break Loop
			}
			writer.Write(downsample(buf))
		}
		chunks = append(chunks, chunk.Bytes())
	}
	err = nil
	return
}

func uniqueStrings(inputs []string) (set map[string]bool) {
	set = map[string]bool{}
	for _, input := range inputs {
		set[input] = true
	}
	return
}

// ListMp3Urls links the urls in the page
func ListMp3Urls(url string) (urls map[string]bool, err error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	page, err := ioutil.ReadAll(resp.Body)
	//page, err := ioutil.ReadFile("/tmp/page")
	if err != nil {
		fmt.Println(err)
	}
	r, err := regexp.Compile("http://open.live.bbc.co.uk/mediaselector/6/redir/version/2.0/mediaset/audio-nondrm-download-low/proto/http/vpid/([a-z]|[0-9]){8}.mp3")
	if err != nil {
		fmt.Println(err)
	}
	urls = uniqueStrings(r.FindAllString(string(page), -1))
	return
}

// ListLinks found in the page
func ListLinks() (links map[string]bool, err error) {
	resp, err := http.Get("https://www.bbc.co.uk/worldserviceradio")
	if err != nil {
		fmt.Println(err)
	}
	page, err := ioutil.ReadAll(resp.Body)
	//page, err := ioutil.ReadFile("/tmp/page")
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := regexp.Compile("http://www.bbc.co.uk/programmes/([a-z]|[0-9]){8}")
	if err != nil {
		fmt.Println(err)
		return
	}
	candidates := r.FindAllString(string(page), -1)
	links = map[string]bool{}
	for _, input := range candidates {
		links[input] = true
	}
	return
}

// NewFakeRadio creates a new fake radio which will give chunks of audio from BBC news.
func NewFakeRadio() FakeRadio {
	return FakeRadio{
		urlsPointer:   0,
		urls:          []string{},
		chunksPointer: 0,
		chunks:        [][]byte{},
	}
}

// FakeRadio is holds the state for the radio.
type FakeRadio struct {
	urlsPointer   int
	urls          []string
	chunksPointer int
	chunks        [][]byte
}

func (fr *FakeRadio) refreshUrls() {
	fmt.Println("fetching a new set of URLs")

	links, err := ListLinks()
	if err != nil {
		panic(err)
	}
	for link := range links {
		mp3s, err := ListMp3Urls(link)
		if err != nil {
			panic(err)
		}
		for l := range mp3s {
			fr.urls = append(fr.urls, l)
		}
	}
	if len(fr.urls) == 0 {
		panic("got 0 URLs, something is very wrong, please contact the developers.")
	}
	fr.urlsPointer = len(fr.urls)
	fmt.Println("got", len(fr.urls), "urls")
}

func (fr *FakeRadio) getNextURL() string {
	if fr.urlsPointer == 0 {
		fr.refreshUrls()
	}
	fr.urlsPointer--
	return fr.urls[fr.urlsPointer]
}

func (fr *FakeRadio) refreshChunks() {
	var err error
	fr.chunks, err = DownloadAndSplit(fr.getNextURL())
	if err != nil {
		panic(err)
	}
	fr.chunksPointer = len(fr.chunks)
	fmt.Println("got", len(fr.chunks), "chunks of fake audio")
	return
}

// GetNextChunk may fetch a new audio file
func (fr *FakeRadio) GetNextChunk() (chunk []byte) {
	if fr.chunksPointer == 0 {
		fr.refreshChunks()
	}
	fr.chunksPointer--
	return fr.chunks[fr.chunksPointer]
}
