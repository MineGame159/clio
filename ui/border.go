package ui

type BorderRunes struct {
	Horizontal string
	Vertical   string

	CornerTopLeft     string
	CornerTopRight    string
	CornerBottomRight string
	CornerBottomLeft  string

	TeeTop    string
	TeeRight  string
	TeeBottom string
	TeeLeft   string
}

var Sharp = &BorderRunes{
	Horizontal: "─",
	Vertical:   "│",

	CornerTopLeft:     "┌",
	CornerTopRight:    "┐",
	CornerBottomRight: "┘",
	CornerBottomLeft:  "└",

	TeeTop:    "┬",
	TeeRight:  "┤",
	TeeBottom: "┴",
	TeeLeft:   "├",
}

var Rounded = &BorderRunes{
	Horizontal: "─",
	Vertical:   "│",

	CornerTopLeft:     "╭",
	CornerTopRight:    "╮",
	CornerBottomRight: "╯",
	CornerBottomLeft:  "╰",

	TeeTop:    "┬",
	TeeRight:  "┤",
	TeeBottom: "┴",
	TeeLeft:   "├",
}
