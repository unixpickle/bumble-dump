// Command top_correlations computes word correlations.
package main

import (
	"context"
	"fmt"

	"github.com/unixpickle/bumble-dump"
	"github.com/unixpickle/essentials"
)

func main() {
	db, err := bumble.OpenDatabase(bumble.GetConfig())
	essentials.Must(err)

	doZodiacSigns(db)
}

func doZodiacSigns(db bumble.Database) {
	fmt.Println("Zodiac sign correlations:")
	correlations, err := bumble.WordCorrelations(context.Background(), db,
		func(u *bumble.User) bool {
			for _, field := range u.ProfileFields {
				if field.ID == "lifestyle_zodiak" {
					return true
				}
			}
			return false
		})
	essentials.Must(err)
	printTopCorrelations(correlations)
}

func printTopCorrelations(m map[string]float64, ignore ...string) {
	var words []string
	var corr []float64
	for w, c := range m {
		words = append(words, w)
		corr = append(corr, c)
	}
	essentials.VoodooSort(corr, func(i, j int) bool {
		return corr[i] < corr[j]
	}, words)
	count := 0
	for i := len(corr) - 1; i >= 0 && count < 20; i-- {
		if !essentials.Contains(ignore, words[i]) {
			fmt.Println(words[i], corr[i])
			count++
		}
	}
}

type correlationPair struct {
	Word string
	Corr float64
}
