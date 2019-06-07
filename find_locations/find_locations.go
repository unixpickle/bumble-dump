// Command find_locations populates the locations
// collection with geocoordinates for every location in
// every downloaded user profile.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/unixpickle/bumble-dump"
	"github.com/unixpickle/essentials"
)

func main() {
	db, err := bumble.OpenDatabase(bumble.GetConfig())
	essentials.Must(err)

	locs, err := db.AllUserLocations(context.Background())
	essentials.Must(err)

	for _, loc := range locs {
		if _, err := db.GetLocation(loc); err == nil {
			continue
		}
		log.Println("looking up:", loc)
		lat, lon, err := lookupLocation(loc)
		if err != nil {
			log.Println("error:", err)
			continue
		}
		essentials.Must(db.AddLocation(&bumble.Location{Name: loc, Lat: lat, Lon: lon}))
	}
}

func lookupLocation(name string) (lat, lon float64, err error) {
	dataStr := "address=" + url.QueryEscape(name)
	body := bytes.NewReader([]byte(dataStr))
	req, err := http.NewRequest("POST", "https://www.mapdevelopers.com/data.php?operation=geocode",
		body)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Referer", "https://www.mapdevelopers.com/geocode_tool.php")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var obj struct {
		Data struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lng"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return 0, 0, err
	}
	return obj.Data.Lat, obj.Data.Lon, nil
}
