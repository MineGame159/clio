package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"strings"
	"unicode"
)

func Capitalize(input string) string {
	runes := []rune(input)
	var b strings.Builder

	for i, r := range runes {
		if i == 0 {
			b.WriteRune(unicode.ToTitle(r))
			continue
		}

		if unicode.IsUpper(r) {
			prev := runes[i-1]

			if !unicode.IsUpper(prev) {
				b.WriteRune(' ')
			} else {
				if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					b.WriteRune(' ')
				}
			}
		}

		b.WriteRune(r)
	}

	return b.String()
}

func Count[T any](it iter.Seq[T]) uint {
	count := uint(0)

	for range it {
		count++
	}

	return count
}

func GetJson[T any](url string) (T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		var empty T
		return empty, err
	}

	req.Header.Set("User-Agent", "clio")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		var empty T
		return empty, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var empty T
		return empty, errors.New(fmt.Sprintf("request failed with status code: %d '%s'", res.StatusCode, res.Status))
	}

	var body T
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		var empty T
		return empty, err
	}

	return body, nil
}
