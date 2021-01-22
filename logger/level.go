package logger

type Level int

const (
	LTrace Level = 1 << iota
	LDebug
	LInfo
	LWarn
	LError
	LPanic
	LFatal
	LHTTP
	LstdLevel = LTrace | LDebug | LInfo | LWarn | LError | LPanic | LFatal | LHTTP
)

func (rec Level) string() string {
	switch rec {
	case LTrace:
		return "Trace"

	case LDebug:
		return "Debug"

	case LInfo:
		return "Info"

	case LWarn:
		return "Warn"

	case LError:
		return "Error"

	case LPanic:
		return "Panic"

	case LFatal:
		return "Fatal"

	case LHTTP:
		return "HTTP"

	default:
		return ""
	}
}

func (rec Level) color() *color {
	switch rec {
	case LTrace:
		return newColor(FgBrBlue)

	case LDebug:
		return newColor(FgBrCyan)

	case LInfo:
		return newColor(FgBrGreen)

	case LWarn:
		return newColor(FgBrYellow)

	case LError:
		return newColor(FgBrRed)

	case LPanic:
		return newColor(FgBrMagenta)

	case LFatal:
		return newColor(FgRed, BgBrBlue)

	case LHTTP:
		return newColor(FgBrWhite)

	default:
		return newColor(FgBrWhite)
	}
}

func (rec Level) contains(level Level) bool {
	return rec&level != 0
}
