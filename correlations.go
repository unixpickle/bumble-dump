package bumble

import (
	"context"
	"math"
	"strings"
	"unicode"
)

// WordCorrelations gets the correlation between each word
// and a binary variable determined by v.
func WordCorrelations(ctx context.Context, db Database,
	v func(u *User) bool) (map[string]float64, error) {
	users, errCh := db.AllUsers(ctx)
	occur := map[string]int{}
	cooccur := map[string]int{}
	numUsers := 0
	vCount := 0
	for user := range users {
		numUsers++
		isV := v(user)
		if isV {
			vCount++
		}
		for w := range WordsInBio(user) {
			occur[w]++
			if isV {
				cooccur[w]++
			}
		}
	}
	if err := <-errCh; err != nil {
		return nil, err
	}

	vMean := float64(vCount) / float64(numUsers)
	vNorm := math.Sqrt(float64(vCount)*math.Pow(1-vMean, 2) +
		float64(numUsers-vCount)*math.Pow(-vMean, 2))

	correlations := map[string]float64{}
	for word, count := range occur {
		cocount := cooccur[word]
		mean := float64(count) / float64(numUsers)
		norm := math.Sqrt(float64(count)*math.Pow(1-mean, 2) +
			float64(numUsers-count)*math.Pow(-mean, 2))
		dotProduct := float64(cocount)*(1-mean)*(1-vMean) +
			float64(count-cocount)*(1-mean)*-vMean +
			float64(vCount-cocount)*-mean*(1-vMean) +
			float64(numUsers-count-vCount+cocount)*mean*vMean
		correlations[word] = dotProduct / (norm * vNorm)
	}

	return correlations, nil
}

// WordsInBio cleans up and extracts words from a user's
// bio and puts them into a multi-set.
func WordsInBio(u *User) map[string]int {
	res := map[string]int{}
	var bio string
	for _, field := range u.ProfileFields {
		if field.ID == "aboutme_text" {
			bio = field.DisplayValue
		}
	}
	if bio == "" {
		return res
	}
	bio = strings.Replace(bio, "/", " ", -1)
	for _, field := range strings.Fields(bio) {
		runes := []rune(field)
		for len(runes) > 0 && !unicode.IsLetter(runes[0]) {
			runes = runes[1:]
		}
		for len(runes) > 0 && !unicode.IsLetter(runes[len(runes)-1]) {
			runes = runes[:len(runes)-1]
		}
		if len(runes) > 0 {
			res[strings.ToLower(string(runes))]++
		}
	}
	return res
}
