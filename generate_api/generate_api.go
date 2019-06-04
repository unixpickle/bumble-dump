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

	bumble "github.com/unixpickle/bumble-dump"
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
	PostData PostData  `json:"postData"`
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
	var obj HAR
	essentials.Must(json.NewDecoder(os.Stdin).Decode(&obj))

	var api bumble.BumbleAPI
	var foundLoc, foundGetEnc, foundDislike bool
	for _, entry := range obj.Log.Entries {
		call := bumble.BumbleCall{
			URL:      entry.Request.URL,
			Headers:  map[string]string{},
			PostBody: entry.Request.PostData.Text,
		}
		for _, header := range entry.Request.Headers {
			call.Headers[header.Name] = header.Value
		}
		if call.URL == "https://bumble.com/unified-api.phtml?SERVER_UPDATE_LOCATION" {
			api.UpdateLocationCall = bumble.UpdateLocationCall{BumbleCall: call}
			foundLoc = true
		} else if call.URL == "https://bumble.com/unified-api.phtml?SERVER_GET_ENCOUNTERS" {
			api.GetEncountersCall = bumble.GetEncountersCall{BumbleCall: call}
			foundGetEnc = true
		} else if call.URL == "https://bumble.com/unified-api.phtml?SERVER_ENCOUNTERS_VOTE" {
			api.DislikeCall = bumble.DislikeCall{BumbleCall: call}
			foundDislike = true
		}
	}
	if !foundLoc || !foundGetEnc || !foundDislike {
		essentials.Die("missing a request (did you remember to swipe someone?)")
	}

	json.NewEncoder(os.Stdout).Encode(api)
}
