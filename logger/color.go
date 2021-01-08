package logger

import (
	"fmt"
	"strconv"
	"strings"
)

type SGRCode int

// Base attributes
const (
	Reset SGRCode = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
	Primary
)

// Foreground text colors
const (
	FgBlack SGRCode = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground bright text colors
const (
	FgBrBlack SGRCode = iota + 90
	FgBrRed
	FgBrGreen
	FgBrYellow
	FgBrBlue
	FgBrMagenta
	FgBrCyan
	FgBrWhite
)

// Background text colors
const (
	BgBlack SGRCode = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background bright text colors
const (
	BgBrBlack SGRCode = iota + 100
	BgBrRed
	BgBrGreen
	BgBrYellow
	BgBrBlue
	BgBrMagenta
	BgBrCyan
	BgBrWhite
)

func (rec SGRCode) String() string {
	return strconv.Itoa(int(rec))
}

type color struct {
	c string
}

func newColor(codes ...SGRCode) *color {
	var ss []string
	for _, v := range codes {
		ss = append(ss, v.String())
	}
	return &color{strings.Join(ss, ";")}
}

func (rec color) println(s string) {
	c := fmt.Sprintf("\x1b[%sm%s\x1b[0m", rec.c, s)
	fmt.Println(c)
}
