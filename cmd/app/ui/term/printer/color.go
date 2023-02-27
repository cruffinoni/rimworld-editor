package printer

import "regexp"

const (
	Reset = iota
	Bold
	Faint
	Underlined = 4
	SlowBlink  = 5
)

const (
	ForegroundBlack = iota + 30
	ForegroundRed
	ForegroundGreen
	ForegroundYellow
	ForegroundBlue
	ForegroundMagenta
	ForegroundCyan
	ForegroundWhite
)

const (
	BackgroundBlack = iota + 40
	BackgroundRed
	BackgroundGreen
	BackgroundYellow
	BackgroundBlue
	BackgroundMagenta
	BackgroundCyan
	BackgroundWhite
)

var (
	colorFinderRegex = regexp.MustCompile(`\{-?([\w,_]*)\}`)
	colorValues      = map[string]int{
		"black":   0,
		"red":     1,
		"green":   2,
		"yellow":  3,
		"cyan":    4,
		"magenta": 5,
	}
	colorOptions = map[string]int{
		"reset":      Reset,
		"bold":       Bold,
		"faint":      Faint,
		"underlined": Underlined,
		"slowBlink":  SlowBlink,
	}
)
