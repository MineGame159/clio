package core

import (
	"iter"
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
