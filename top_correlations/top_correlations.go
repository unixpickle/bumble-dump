// Command top_correlations computes word correlations.
package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/unixpickle/bumble-dump"
	"github.com/unixpickle/essentials"
)

func main() {
	db, err := bumble.OpenDatabase(bumble.GetConfig())
	essentials.Must(err)

	doCountry(db, "us")
	doGender(db, "Male", 1)
	doGender(db, "Female", 2)
	doUnder24(db)
	doOver40(db)
	doOverSixFoot(db)
	doZodiacSigns(db)
}

func doCountry(db bumble.Database, countryCode string) {
	fmt.Println("Country =", countryCode, "correlations:")
	countryLocs := map[string]bool{}
	locs, errCh := db.AllLocations(context.Background())
	for loc := range locs {
		if loc.CountryCode == countryCode {
			countryLocs[loc.Name] = true
		}
	}
	essentials.Must(<-errCh)
	correlations, err := bumble.WordCorrelations(context.Background(), db,
		func(u *bumble.User) bool {
			return countryLocs[u.Location]
		})
	essentials.Must(err)
	printTopCorrelations(correlations)
}

func doGender(db bumble.Database, genderStr string, genderNum int) {
	fmt.Println("Gender =", genderStr, "correlations:")
	correlations, err := bumble.WordCorrelations(context.Background(), db,
		func(u *bumble.User) bool {
			return u.Gender == genderNum
		})
	essentials.Must(err)
	printTopCorrelations(correlations)
}

func doUnder24(db bumble.Database) {
	fmt.Println("Age < 24 correlations:")
	correlations, err := bumble.WordCorrelations(context.Background(), db,
		func(u *bumble.User) bool {
			return u.Age < 24
		})
	essentials.Must(err)
	printTopCorrelations(correlations)
}

func doOver40(db bumble.Database) {
	fmt.Println("Age >= 40 correlations:")
	correlations, err := bumble.WordCorrelations(context.Background(), db,
		func(u *bumble.User) bool {
			return u.Age >= 40
		})
	essentials.Must(err)
	printTopCorrelations(correlations)
}

func doOverSixFoot(db bumble.Database) {
	fmt.Println("Height > 6ft correlations:")
	correlations, err := bumble.WordCorrelations(context.Background(), db,
		func(u *bumble.User) bool {
			for _, field := range u.ProfileFields {
				if field.ID == "lifestyle_height" {
					heightCm := 0
					fields := strings.Fields(strings.Replace(field.DisplayValue, "(", "", -1))
					for i, f := range fields[1:] {
						if f == "cm" {
							heightCm, _ = strconv.Atoi(fields[i])
							break
						}
					}
					return heightCm >= 183
				}
			}
			return false
		})
	essentials.Must(err)
	printTopCorrelations(correlations)
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
	fmt.Println()
}
