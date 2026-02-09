package main

import (
	"fmt"
	"net"
	"strings"
)


const (
	Reset = "\x1b[0m"
	
	Blue       = "\x1b[38;5;33m"
	DarkBlue   = "\x1b[38;5;27m" 
	LightBlue  = "\x1b[38;5;39m" 
	Cyan       = "\x1b[38;5;51m" 
	DarkCyan   = "\x1b[38;5;45m" 
	LightCyan  = "\x1b[38;5;87m" 
	
	White      = "\x1b[37m"
	BrightWhite = "\x1b[97m"
	Gray       = "\x1b[90m"
	DarkGray   = "\x1b[38;5;240m"

	Green      = "\x1b[38;5;46m"
	Yellow     = "\x1b[38;5;226m"
	Red        = "\x1b[38;5;196m"
	Orange     = "\x1b[38;5;214m"
	Purple     = "\x1b[38;5;141m"
	
	Bold       = "\x1b[1m"
	Dim        = "\x1b[2m"
	Italic     = "\x1b[3m"
	Underline  = "\x1b[4m"
	Blink      = "\x1b[5m"
	Reverse    = "\x1b[7m"
	
	// Backgrounds
	BgBlue     = "\x1b[48;5;33m"
	BgCyan     = "\x1b[48;5;51m"
	BgGreen    = "\x1b[48;5;46m"
	BgRed      = "\x1b[48;5;196m"
	BgYellow   = "\x1b[48;5;226m"
	BgDark     = "\x1b[48;5;16m"
	BgWhite    = "\x1b[48;5;255m"
)

var GradientColors = []string{
	"\x1b[38;5;17m",
	"\x1b[38;5;18m",
	"\x1b[38;5;19m",
	"\x1b[38;5;20m",
	"\x1b[38;5;21m",
	"\x1b[38;5;27m",
	"\x1b[38;5;33m",
	"\x1b[38;5;39m",
	"\x1b[38;5;45m",
	"\x1b[38;5;51m",
	"\x1b[38;5;87m",
	"\x1b[38;5;123m",
}

var (
	ColorSuccess = Green
	ColorError   = Red
	ColorWarning = Yellow
	ColorInfo    = Cyan
	
	ColorPrompt  = Blue
	ColorInput   = LightCyan
	ColorBorder  = DarkBlue
	ColorTitle   = Bold + Cyan

	ColorHeader  = Bold + Blue
	ColorData    = White
	ColorHighlight = LightCyan
	ColorDim     = Gray
	
	ColorOnline  = Green
	ColorOffline = Red
	ColorIdle    = Yellow
)

func ApplyGradient(text string) string {
	if len(text) == 0 {
		return text
	}
	
	var result strings.Builder
	textLen := len(text)
	colorLen := len(GradientColors)
	
	for i, char := range text {
		if char == '\n' || char == '\r' {
			result.WriteRune(char)
			continue
		}
		
		colorIndex := (i * colorLen) / textLen
		if colorIndex >= colorLen {
			colorIndex = colorLen - 1
		}
		
		result.WriteString(GradientColors[colorIndex])
		result.WriteRune(char)
	}
	
	result.WriteString(Reset)
	return result.String()
}

func ApplyGradientToLines(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	
	result := make([]string, len(lines))
	colorLen := len(GradientColors)
	
	for i, line := range lines {
		colorIndex := (i * colorLen) / len(lines)
		if colorIndex >= colorLen {
			colorIndex = colorLen - 1
		}
		
		result[i] = GradientColors[colorIndex] + line + Reset
	}
	
	return result
}

func ColoredBanner() string {
	banner := []string{
		"           ▄▄▄ ▄▄▄▄▄▄ ▄▄▄  ▄▄   ▄▄ ▄▄  ▄▄▄▄ ",
		"          ██▀██  ██  ██▀██ ██▀▄▀██ ██ ██▀▀▀ ",
		"          ██▀██  ██  ▀███▀ ██   ██ ██ ▀████ ",
		"        Atomic C2 by CirqueiraDev | ATOM TEAM ",
	}
	
	coloredLines := ApplyGradientToLines(banner)
	return "\r\n" + strings.Join(coloredLines, "\r\n") + "\r\n\r\n"
}

func ShowColoredBanner(conn net.Conn) {
	conn.Write([]byte(ColoredBanner()))
}

func ColorText(text, color string) string {
	return color + text + Reset
}

func SuccessMsg(msg string) string {
	return fmt.Sprintf("%s%s Success %s %s%s%s", BgGreen, BgDark, Reset, ColorSuccess, msg, Reset)
}

func ErrorMsg(msg string) string {
	return fmt.Sprintf("%s%s Error %s %s%s%s", BgRed, BgDark, Reset, ColorError, msg, Reset)
}

func WarningMsg(msg string) string {
	return fmt.Sprintf("%s%s Warning %s %s%s%s", BgYellow, BgDark, Reset, ColorWarning, msg, Reset)
}

func InfoMsg(msg string) string {
	return fmt.Sprintf("%s%s Info %s %s%s%s", BgCyan, BgDark, Reset, ColorInfo, msg, Reset)
}

func Box(title, content string) string {
	lines := strings.Split(content, "\n")
	maxLen := len(title)
	
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	
	var box strings.Builder
	
	// Top border
	box.WriteString(ColorBorder + "╔" + strings.Repeat("═", maxLen+2) + "╗" + Reset + "\n")
	
	// Title
	if title != "" {
		padding := maxLen - len(title)
		box.WriteString(ColorBorder + "║ " + ColorTitle + title + strings.Repeat(" ", padding) + ColorBorder + " ║" + Reset + "\n")
		box.WriteString(ColorBorder + "╠" + strings.Repeat("═", maxLen+2) + "╣" + Reset + "\n")
	}
	
	// Content
	for _, line := range lines {
		padding := maxLen - len(line)
		box.WriteString(ColorBorder + "║ " + ColorData + line + strings.Repeat(" ", padding) + ColorBorder + " ║" + Reset + "\n")
	}
	
	box.WriteString(ColorBorder + "╚" + strings.Repeat("═", maxLen+2) + "╝" + Reset + "\n")
	
	return box.String()
}

func ProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}
	
	percentage := float64(current) / float64(total)
	filled := int(float64(width) * percentage)
	
	var bar strings.Builder
	bar.WriteString(ColorBorder + "[" + Reset)
	
	for i := 0; i < width; i++ {
		if i < filled {
			bar.WriteString(Green + "█" + Reset)
		} else {
			bar.WriteString(DarkGray + "░" + Reset)
		}
	}
	
	bar.WriteString(ColorBorder + "]" + Reset)
	bar.WriteString(fmt.Sprintf(" %s%.1f%%%s", Cyan, percentage*100, Reset))
	
	return bar.String()
}

func TableHeader(columns []string) string {
	var header strings.Builder
	header.WriteString(ColorBorder + "┌")
	for i, col := range columns {
		header.WriteString(strings.Repeat("─", len(col)+2))
		if i < len(columns)-1 {
			header.WriteString("┬")
		}
	}
	header.WriteString("┐" + Reset + "\n")
	
	header.WriteString(ColorBorder + "│ " + Reset)
	for i, col := range columns {
		header.WriteString(ColorHeader + col + Reset)
		if i < len(columns)-1 {
			header.WriteString(ColorBorder + " │ " + Reset)
		}
	}
	header.WriteString(ColorBorder + " │" + Reset + "\n")
	
	header.WriteString(ColorBorder + "├")
	for i, col := range columns {
		header.WriteString(strings.Repeat("─", len(col)+2))
		if i < len(columns)-1 {
			header.WriteString("┼")
		}
	}
	header.WriteString("┤" + Reset + "\n")
	
	return header.String()
}

func Prompt(username string) string {
	return fmt.Sprintf("\r\n%s(%s%s@%satomic%s) %s>%s ", 
		ColorBorder, 
		LightCyan, username, 
		Blue, 
		ColorBorder, 
		Cyan, 
		Reset)
}

func StatusBadge(status string) string {
	var color string
	switch strings.ToLower(status) {
	case "online", "active", "running":
		color = ColorOnline
	case "offline", "inactive", "stopped":
		color = ColorOffline
	case "idle", "pending", "waiting":
		color = ColorIdle
	default:
		color = ColorInfo
	}
	
	return fmt.Sprintf("%s[%s]%s", color, strings.ToUpper(status), Reset)
}

func AnimatedText(text string, frame int) string {
	colors := []string{Blue, LightBlue, Cyan, LightCyan}
	colorIndex := frame % len(colors)
	return colors[colorIndex] + text + Reset
}

func RainbowGradient(text string) string {
	rainbowColors := []string{
		"\x1b[38;5;27m", 
		"\x1b[38;5;33m", 
		"\x1b[38;5;39m", 
		"\x1b[38;5;45m",
		"\x1b[38;5;51m", 
		"\x1b[38;5;87m", 
		"\x1b[38;5;123m",
	}
	
	var result strings.Builder
	textLen := len(text)
	colorLen := len(rainbowColors)
	
	for i, char := range text {
		if char == '\n' || char == '\r' {
			result.WriteRune(char)
			continue
		}
		
		colorIndex := (i * colorLen) / textLen
		if colorIndex >= colorLen {
			colorIndex = colorLen - 1
		}
		
		result.WriteString(rainbowColors[colorIndex])
		result.WriteRune(char)
	}
	
	result.WriteString(Reset)
	return result.String()
}

func ClearScreen(conn net.Conn) {
	conn.Write([]byte("\033[2J\033[1H"))
}

func FormatBoolColored(b bool) string {
	if b {
		return ColorSuccess + "true" + Reset
	}
	return ColorError + "false" + Reset
}

func Separator(width int, char string) string {
	if char == "" {
		char = "─"
	}
	return ColorBorder + strings.Repeat(char, width) + Reset
}

func HighlightText(text, highlight string) string {
	return strings.ReplaceAll(text, highlight, ColorHighlight+Bold+highlight+Reset)
}

func LoadingSpinner(frame int) string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return Cyan + spinners[frame%len(spinners)] + Reset
}