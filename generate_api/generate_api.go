// Command generate_api asks the user to paste in HAR
// information after using the Bumble website and liking
// or disliking one user.
//
// The requests in the HAR data are converted into a
// reusable API for the Bumble service.
package main

import (
	"encoding/json"
	"os"

	"github.com/unixpickle/essentials"
)

type HAR struct {
	Log struct {
		Entries []*LogEntry `json:"entries"`
	} `json:"log"`
}

type LogEntry struct {
	Request *Request `json:"request"`
}

type Request struct {
	Method   string    `json:"method"`
	URL      string    `json:"url"`
	Headers  []*Header `json:"headers"`
	PostData *PostData `json:"postData"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PostData struct {
	MIMEType string `json:"mimeType"`
	Text     string `json:"text"`
}

func main() {
	var obj *HAR
	essentials.Must(json.NewDecoder(os.Stdin).Decode(&obj))

	// TODO: extract APIs here.
}
