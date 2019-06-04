// Command scan automatically dumps Bumble profiles as
// JSON.
package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/unixpickle/bumble-dump"
	"github.com/unixpickle/essentials"
)

const ErrBackoff = time.Minute

func main() {
	if len(os.Args) != 2 {
		essentials.Die("Usage: scan <api.json>")
	}

	var api bumble.BumbleAPI
	f, err := os.Open(os.Args[1])
	essentials.Must(err)
	err = json.NewDecoder(f).Decode(&api)
	f.Close()
	essentials.Must(err)

	enc := json.NewEncoder(os.Stdout)

SearchLoop:
	for {
		lat, lon := randomLocation()
		log.Printf("searching at location: %f,%f", lat, lon)
		if err := api.UpdateLocation(lat, lon); err != nil {
			log.Print(err)
			time.Sleep(ErrBackoff)
			continue
		}

		var numResults int
		for numResults < 1000 {
			users, err := api.GetEncounters()
			if err != nil {
				log.Print(err)
				continue
			}
			if len(users) == 0 {
				log.Print("got 0 results")
				continue SearchLoop
			}
			for _, user := range users {
				enc.Encode(user)
				if err := api.Dislike(user.ID); err != nil {
					log.Print(err)
					time.Sleep(ErrBackoff)
					continue SearchLoop
				}
				numResults += 1
			}
		}
	}
}

func randomLocation() (lat, lon float64) {
	return rand.Float64()*180 - 90, rand.Float64()*360 - 180
}
