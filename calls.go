package bumble

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// A BumbleCall is a pre-recorded bumble API call.
type BumbleCall struct {
	URL      string
	Headers  map[string]string
	PostBody string
}

// recreate creates a new http.Request mimicking the
// original request.
func (b *BumbleCall) recreate() (*http.Request, error) {
	req, err := http.NewRequest("POST", b.URL, bytes.NewReader([]byte(b.PostBody)))
	if err != nil {
		return nil, errors.Wrap(err, "generate request")
	}
	for k, v := range b.Headers {
		req.Header.Add(k, v)
	}
	req.Header.Del("Accept-Encoding")
	return req, nil
}

func (b *BumbleCall) replacedReq(replacements map[string]interface{}) (*http.Request, error) {
	var err error
	bc := *b
	bc.PostBody, err = replacePostBody(bc.PostBody, replacements)
	if err != nil {
		return nil, errors.Wrap(err, "generate request")
	}
	return bc.recreate()
}

type GetEncountersCall struct {
	BumbleCall
}

func (g *GetEncountersCall) Request() (*http.Request, error) {
	return g.recreate()
}

type DislikeCall struct {
	BumbleCall
}

func (d *DislikeCall) Request(userID string) (*http.Request, error) {
	return d.replacedReq(map[string]interface{}{
		"person_id": userID,
		"vote":      3,
	})
}

type UpdateLocationCall struct {
	BumbleCall
}

func (u *UpdateLocationCall) Request(lat, lon float64) (*http.Request, error) {
	return u.replacedReq(map[string]interface{}{
		"longitude": lon,
		"latitude":  lat,
	})
}

// replacePostBody replaces named fields in a JSON POST
// body with new values.
// It recursively searches the posted structure to find
// the named fields.
func replacePostBody(body string, replacements map[string]interface{}) (string, error) {
	var obj interface{}
	if err := json.Unmarshal([]byte(body), &obj); err != nil {
		return "", err
	}
	replacePostBodyObj(obj, replacements)
	newData, _ := json.Marshal(obj)
	return string(newData), nil
}

func replacePostBodyObj(obj interface{}, replacements map[string]interface{}) {
	switch obj := obj.(type) {
	case []interface{}:
		for _, x := range obj {
			replacePostBodyObj(x, replacements)
		}
	case map[string]interface{}:
		for k, v := range obj {
			if rep, ok := replacements[k]; ok {
				obj[k] = rep
			} else {
				replacePostBodyObj(v, replacements)
			}
		}
	}
}
