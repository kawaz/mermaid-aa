package renderer

// BoxStyle defines the characters used for drawing boxes and lines
type BoxStyle struct {
	Name string

	// Box corners
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string

	// Box sides
	Horizontal string
	Vertical   string

	// T-junctions
	TeeDown  string // ┬
	TeeUp    string // ┴
	TeeRight string // ├
	TeeLeft  string // ┤
	Cross    string // ┼

	// Arrows
	ArrowDown  string
	ArrowUp    string
	ArrowRight string
	ArrowLeft  string

	// Rounded corners (for rounded rectangles)
	RoundTopLeft     string
	RoundTopRight    string
	RoundBottomLeft  string
	RoundBottomRight string

	// Diamond
	DiamondTop    string
	DiamondBottom string
	DiamondLeft   string
	DiamondRight  string
}

// ASCII style using basic ASCII characters
var ASCIIStyle = &BoxStyle{
	Name:             "ascii",
	TopLeft:          "+",
	TopRight:         "+",
	BottomLeft:       "+",
	BottomRight:      "+",
	Horizontal:       "-",
	Vertical:         "|",
	TeeDown:          "+",
	TeeUp:            "+",
	TeeRight:         "+",
	TeeLeft:          "+",
	Cross:            "+",
	ArrowDown:        "v",
	ArrowUp:          "^",
	ArrowRight:       ">",
	ArrowLeft:        "<",
	RoundTopLeft:     "/",
	RoundTopRight:    "\\",
	RoundBottomLeft:  "\\",
	RoundBottomRight: "/",
	DiamondTop:       "/\\",
	DiamondBottom:    "\\/",
	DiamondLeft:      "<",
	DiamondRight:     ">",
}

// Unicode light box drawing style
var UnicodeStyle = &BoxStyle{
	Name:             "unicode",
	TopLeft:          "┌",
	TopRight:         "┐",
	BottomLeft:       "└",
	BottomRight:      "┘",
	Horizontal:       "─",
	Vertical:         "│",
	TeeDown:          "┬",
	TeeUp:            "┴",
	TeeRight:         "├",
	TeeLeft:          "┤",
	Cross:            "┼",
	ArrowDown:        "▼",
	ArrowUp:          "▲",
	ArrowRight:       "▶",
	ArrowLeft:        "◀",
	RoundTopLeft:     "╭",
	RoundTopRight:    "╮",
	RoundBottomLeft:  "╰",
	RoundBottomRight: "╯",
	DiamondTop:       "◇",
	DiamondBottom:    "◇",
	DiamondLeft:      "◁",
	DiamondRight:     "▷",
}

// Unicode bold/heavy box drawing style
var UnicodeBoldStyle = &BoxStyle{
	Name:             "unicode-bold",
	TopLeft:          "┏",
	TopRight:         "┓",
	BottomLeft:       "┗",
	BottomRight:      "┛",
	Horizontal:       "━",
	Vertical:         "┃",
	TeeDown:          "┳",
	TeeUp:            "┻",
	TeeRight:         "┣",
	TeeLeft:          "┫",
	Cross:            "╋",
	ArrowDown:        "▼",
	ArrowUp:          "▲",
	ArrowRight:       "▶",
	ArrowLeft:        "◀",
	RoundTopLeft:     "┏",
	RoundTopRight:    "┓",
	RoundBottomLeft:  "┗",
	RoundBottomRight: "┛",
	DiamondTop:       "◆",
	DiamondBottom:    "◆",
	DiamondLeft:      "◀",
	DiamondRight:     "▶",
}

// Unicode double line style
var UnicodeDoubleStyle = &BoxStyle{
	Name:             "unicode-double",
	TopLeft:          "╔",
	TopRight:         "╗",
	BottomLeft:       "╚",
	BottomRight:      "╝",
	Horizontal:       "═",
	Vertical:         "║",
	TeeDown:          "╦",
	TeeUp:            "╩",
	TeeRight:         "╠",
	TeeLeft:          "╣",
	Cross:            "╬",
	ArrowDown:        "▼",
	ArrowUp:          "▲",
	ArrowRight:       "▶",
	ArrowLeft:        "◀",
	RoundTopLeft:     "╔",
	RoundTopRight:    "╗",
	RoundBottomLeft:  "╚",
	RoundBottomRight: "╝",
	DiamondTop:       "◆",
	DiamondBottom:    "◆",
	DiamondLeft:      "◀",
	DiamondRight:     "▶",
}

// Unicode rounded style with nice curves
var UnicodeRoundStyle = &BoxStyle{
	Name:             "unicode-round",
	TopLeft:          "╭",
	TopRight:         "╮",
	BottomLeft:       "╰",
	BottomRight:      "╯",
	Horizontal:       "─",
	Vertical:         "│",
	TeeDown:          "┬",
	TeeUp:            "┴",
	TeeRight:         "├",
	TeeLeft:          "┤",
	Cross:            "┼",
	ArrowDown:        "↓",
	ArrowUp:          "↑",
	ArrowRight:       "→",
	ArrowLeft:        "←",
	RoundTopLeft:     "╭",
	RoundTopRight:    "╮",
	RoundBottomLeft:  "╰",
	RoundBottomRight: "╯",
	DiamondTop:       "◇",
	DiamondBottom:    "◇",
	DiamondLeft:      "◁",
	DiamondRight:     "▷",
}

// GetStyle returns a BoxStyle by name
func GetStyle(name string) *BoxStyle {
	switch name {
	case "unicode", "u":
		return UnicodeStyle
	case "unicode-bold", "bold", "b":
		return UnicodeBoldStyle
	case "unicode-double", "double", "d":
		return UnicodeDoubleStyle
	case "unicode-round", "round", "r":
		return UnicodeRoundStyle
	default:
		return ASCIIStyle
	}
}

// AvailableStyles returns list of available style names
func AvailableStyles() []string {
	return []string{"ascii", "unicode", "unicode-bold", "unicode-double", "unicode-round"}
}
