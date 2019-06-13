// Command generate_api asks the user to paste in HAR
// information after using the Bumble website and liking
// or disliking one user.
//
// The requests in the HAR data are converted into a
// reusable API for the Bumble service.
//
// HAR is fed to standard input, and an encoded API is fed
// to standard output.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
		if strings.Contains(call.URL, "unified-api.phtml?SERVER_UPDATE_LOCATION") {
			api.UpdateLocationCall = bumble.UpdateLocationCall{BumbleCall: call}
			foundLoc = true
		} else if strings.Contains(call.URL, "unified-api.phtml?SERVER_GET_ENCOUNTERS") {
			api.GetEncountersCall = bumble.GetEncountersCall{BumbleCall: call}
			foundGetEnc = true
		} else if strings.Contains(call.URL, "unified-api.phtml?SERVER_ENCOUNTERS_VOTE") {
			api.DislikeCall = bumble.DislikeCall{BumbleCall: call}
			foundDislike = true
		}
	}
	if !foundLoc || !foundGetEnc || !foundDislike {
		if !foundLoc {
			fmt.Fprintln(os.Stderr, "Missing location update request. Try updating your")
			fmt.Fprintln(os.Stderr, "location in settings and allowing your browser to")
			fmt.Fprintln(os.Stderr, "provide your location to the website.")
		}
		if !foundGetEnc {
			fmt.Fprintln(os.Stderr, "Missing encounters request.")
		}
		if !foundDislike {
			fmt.Fprintln(os.Stderr, "Missing dislike request. Make sure to swipe someone.")
		}
		essentials.Die("Cannot generate API due to missing request.")
	}

	json.NewEncoder(os.Stdout).Encode(api)
}
