// Package renderer provides ASCII art rendering for flowcharts
package renderer

// Charset defines the characters used for drawing
type Charset struct {
	Name string

	// Box drawing
	HLine    string // horizontal line
	VLine    string // vertical line
	Corner   string // corner/intersection
	TopLeft  string
	TopRight string
	BotLeft  string
	BotRight string

	// Rounded corners
	RoundTopLeft  string
	RoundTopRight string
	RoundBotLeft  string
	RoundBotRight string

	// Arrows
	ArrowRight string
	ArrowLeft  string
	ArrowUp    string
	ArrowDown  string

	// Diamond
	DiamondLeft  string
	DiamondRight string

	// Hexagon
	HexLeft  string
	HexRight string

	// Circle/Stadium parentheses
	ParenOpen  string
	ParenClose string

	// Dotted line
	DottedH string
	DottedV string

	// Thick line
	ThickH string
	ThickV string
}

// Predefined character sets
var (
	ASCII = &Charset{
		Name:          "ascii",
		HLine:         "-",
		VLine:         "|",
		Corner:        "+",
		TopLeft:       "+",
		TopRight:      "+",
		BotLeft:       "+",
		BotRight:      "+",
		RoundTopLeft:  "+",
		RoundTopRight: "+",
		RoundBotLeft:  "+",
		RoundBotRight: "+",
		ArrowRight:    ">",
		ArrowLeft:     "<",
		ArrowUp:       "^",
		ArrowDown:     "v",
		DiamondLeft:   "<",
		DiamondRight:  ">",
		HexLeft:       "{",
		HexRight:      "}",
		ParenOpen:     "(",
		ParenClose:    ")",
		DottedH:       "-",
		DottedV:       ":",
		ThickH:        "=",
		ThickV:        "#",
	}

	Unicode = &Charset{
		Name:          "unicode",
		HLine:         "─",
		VLine:         "│",
		Corner:        "┼",
		TopLeft:       "┌",
		TopRight:      "┐",
		BotLeft:       "└",
		BotRight:      "┘",
		RoundTopLeft:  "┌",
		RoundTopRight: "┐",
		RoundBotLeft:  "└",
		RoundBotRight: "┘",
		ArrowRight:    "→",
		ArrowLeft:     "←",
		ArrowUp:       "↑",
		ArrowDown:     "↓",
		DiamondLeft:   "◇",
		DiamondRight:  "◇",
		HexLeft:       "⬡",
		HexRight:      "⬡",
		ParenOpen:     "(",
		ParenClose:    ")",
		DottedH:       "┄",
		DottedV:       "┆",
		ThickH:        "━",
		ThickV:        "┃",
	}

	UnicodeRound = &Charset{
		Name:          "unicode-round",
		HLine:         "─",
		VLine:         "│",
		Corner:        "┼",
		TopLeft:       "╭",
		TopRight:      "╮",
		BotLeft:       "╰",
		BotRight:      "╯",
		RoundTopLeft:  "╭",
		RoundTopRight: "╮",
		RoundBotLeft:  "╰",
		RoundBotRight: "╯",
		ArrowRight:    "→",
		ArrowLeft:     "←",
		ArrowUp:       "↑",
		ArrowDown:     "↓",
		DiamondLeft:   "◇",
		DiamondRight:  "◇",
		HexLeft:       "⬡",
		HexRight:      "⬡",
		ParenOpen:     "(",
		ParenClose:    ")",
		DottedH:       "┄",
		DottedV:       "┆",
		ThickH:        "━",
		ThickV:        "┃",
	}

	UnicodeBold = &Charset{
		Name:          "unicode-bold",
		HLine:         "━",
		VLine:         "┃",
		Corner:        "╋",
		TopLeft:       "┏",
		TopRight:      "┓",
		BotLeft:       "┗",
		BotRight:      "┛",
		RoundTopLeft:  "┏",
		RoundTopRight: "┓",
		RoundBotLeft:  "┗",
		RoundBotRight: "┛",
		ArrowRight:    "➤",
		ArrowLeft:     "◀",
		ArrowUp:       "▲",
		ArrowDown:     "▼",
		DiamondLeft:   "◆",
		DiamondRight:  "◆",
		HexLeft:       "⬢",
		HexRight:      "⬢",
		ParenOpen:     "(",
		ParenClose:    ")",
		DottedH:       "┅",
		DottedV:       "┇",
		ThickH:        "━",
		ThickV:        "┃",
	}

	UnicodeDouble = &Charset{
		Name:          "unicode-double",
		HLine:         "═",
		VLine:         "║",
		Corner:        "╬",
		TopLeft:       "╔",
		TopRight:      "╗",
		BotLeft:       "╚",
		BotRight:      "╝",
		RoundTopLeft:  "╔",
		RoundTopRight: "╗",
		RoundBotLeft:  "╚",
		RoundBotRight: "╝",
		ArrowRight:    "⇒",
		ArrowLeft:     "⇐",
		ArrowUp:       "⇑",
		ArrowDown:     "⇓",
		DiamondLeft:   "◆",
		DiamondRight:  "◆",
		HexLeft:       "⬢",
		HexRight:      "⬢",
		ParenOpen:     "(",
		ParenClose:    ")",
		DottedH:       "╌",
		DottedV:       "╎",
		ThickH:        "═",
		ThickV:        "║",
	}

	// CharsetMap maps charset names to charsets
	CharsetMap = map[string]*Charset{
		"ascii":          ASCII,
		"unicode":        Unicode,
		"unicode-round":  UnicodeRound,
		"unicode-bold":   UnicodeBold,
		"unicode-double": UnicodeDouble,
	}
)

// GetCharset returns a charset by name, defaults to Unicode
func GetCharset(name string) *Charset {
	if cs, ok := CharsetMap[name]; ok {
		return cs
	}
	return Unicode
}
