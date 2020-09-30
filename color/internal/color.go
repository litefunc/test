package internal

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

type Color struct {
	c string
}

func (rec SGRCode) String() string {
	return strconv.Itoa(int(rec))
}

func NewColor(codes ...SGRCode) *Color {
	var ss []string
	for _, v := range codes {
		ss = append(ss, v.String())
	}
	return &Color{strings.Join(ss, ";")}
}

func (rec Color) Println(s string) {
	c := fmt.Sprintf("\x1b[%sm%s\x1b[0m", rec.c, s)
	fmt.Println(c)
}

func Log(color int, s string) {
	c := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)

	fmt.Println(c)
}

func Bg(color int, s string) {
	c := fmt.Sprintf("\x1b[;%dm%s\x1b[0m", color, s)

	fmt.Println(c)
}

func B(color int, s string) {
	c := fmt.Sprintf("\x1b[0;0;%dm%s\x1b[0m", color, s)

	fmt.Println(c)
}
