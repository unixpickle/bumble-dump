// Command all_locations dumps a list of locations from
// the database, along with the number of occurrences of
// each location.
package main

import (
	"context"
	"fmt"

	"github.com/unixpickle/bumble-dump"
	"github.com/unixpickle/essentials"
)

func main() {
	locs := map[string]int{}
	db, err := bumble.OpenDatabase(bumble.GetConfig())
	essentials.Must(err)
	users, errCh := db.AllUsers(context.Background())
	for u := range users {
		locs[u.Location]++
	}
	if err := <-errCh; err != nil {
		essentials.Die(err)
	}

	var locations []string
	var counts []int
	for l, c := range locs {
		locations = append(locations, l)
		counts = append(counts, c)
	}
	essentials.VoodooSort(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	}, locations)
	for i, count := range counts {
		fmt.Println(count, locations[i])
	}
}
