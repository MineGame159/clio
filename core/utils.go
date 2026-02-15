package core

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

const (
	Gray   = lipgloss.Color("0")
	Red    = lipgloss.Color("1")
	Green  = lipgloss.Color("2")
	Yellow = lipgloss.Color("3")
	Blue   = lipgloss.Color("4")
	Purple = lipgloss.Color("5")
	Cyan   = lipgloss.Color("6")
	White  = lipgloss.Color("7")

	LightGray   = lipgloss.Color("8")
	LightRed    = lipgloss.Color("9")
	LightGreen  = lipgloss.Color("10")
	LightYellow = lipgloss.Color("11")
	LightBlue   = lipgloss.Color("12")
	LightPurple = lipgloss.Color("13")
	LightCyan   = lipgloss.Color("14")
	LightWhite  = lipgloss.Color("15")
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
