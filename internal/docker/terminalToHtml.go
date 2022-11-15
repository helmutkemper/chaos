package docker

import (
	"bytes"
	"golang.org/x/text/encoding/charmap"
	"html"
	"regexp"
	"strconv"
	"strings"
)

const (
	KAnsiReset = `\u001b[0m`

	KAnsiBold      = `\u001b[1m`
	KAnsiUnderline = `\u001b[4m`
	KAnsiReversed  = `\u001b[7m`

	KAnsiBlack  = `\u001b[30m`
	KAnsiRed    = `\u001b[31m`
	KAnsiGreen  = `\u001b[32m`
	KAnsiYellow = `\u001b[33m`
	KAnsiBlue   = `\u001b[34m`
	KAnsiPurple = `\u001b[35m`
	KAnsiCyan   = `\u001b[36m`
	KAnsiWhite  = `\u001b[37m`

	KAnsiBrightBlack  = `\u001b[90m`
	KAnsiBrightRed    = `\u001b[91m`
	KAnsiBrightGreen  = `\u001b[92m`
	KAnsiBrightYellow = `\u001b[93m`
	KAnsiBrightBlue   = `\u001b[94m`
	KAnsiBrightPurple = `\u001b[95m`
	KAnsiBrightCyan   = `\u001b[96m`
	KAnsiBrightWhite  = `\u001b[97m`

	KAnsiBgBlack  = `\u001b[40m`
	KAnsiBgRed    = `\u001b[41m`
	KAnsiBgGreen  = `\u001b[42m`
	KAnsiBgYellow = `\u001b[43m`
	KAnsiBgBlue   = `\u001b[44m`
	KAnsiBgPurple = `\u001b[45m`
	KAnsiBgCyan   = `\u001b[46m`
	KAnsiBgWhite  = `\u001b[47m`

	KAnsiBrightBgBlack  = `\u001b[100m`
	KAnsiBrightBgRed    = `\u001b[101m`
	KAnsiBrightBgGreen  = `\u001b[102m`
	KAnsiBrightBgYellow = `\u001b[103m`
	KAnsiBrightBgBlue   = `\u001b[104m`
	KAnsiBrightBgPurple = `\u001b[105m`
	KAnsiBrightBgCyan   = `\u001b[106m`
	KAnsiBrightBgWhite  = `\u001b[107m`
)

type TerminalColor struct {
	resetCounter int
}

func (el *TerminalColor) AnsiColor8ToHtmlColor(ansiValue []byte) (htmlTag []byte) {
	if bytes.Compare([]byte(KAnsiReset), ansiValue) == 0 {
		for i := 0; i != el.resetCounter; i += 1 {
			htmlTag = append(htmlTag, []byte("</span>")...)
		}
		el.resetCounter = 0
		return
	}

	if bytes.Compare([]byte(KAnsiBold), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='font-weight: bold;'>")
		return
	}

	if bytes.Compare([]byte(KAnsiUnderline), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='text-decoration: underline;'>")
		return
	}

	if bytes.Compare([]byte(KAnsiReversed), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='filter: invert(100%);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBlack), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(0,0,0);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiRed), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(255,0,0);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiGreen), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(0,255,0);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiYellow), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(255,255,0);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBlue), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(0,0,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiPurple), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(128,0,128);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiCyan), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(0,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiWhite), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(245,245,245);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBlack), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(105,105,105);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightRed), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(205,51,51);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightGreen), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(127,255,0);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightYellow), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(255,255,127);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBlue), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(30,144,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightPurple), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(155,48,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightCyan), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(0,238,238);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightWhite), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='color: rgb(255,255,255);'>")
		return
	}

	// background
	if bytes.Compare([]byte(KAnsiBgBlack), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(0,0,0); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBgRed), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(255,0,0); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBgGreen), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(0,255,0); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBgYellow), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(255,255,0); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBgBlue), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(0,0,255); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBgPurple), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(128,0,128); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBgCyan), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(0,255,255); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBgWhite), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(245,245,245); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBgBlack), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(105,105,105); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBgRed), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(205,51,51); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBgGreen), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(127,255,0); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBgYellow), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(255,255,127); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBgBlue), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(30,144,255); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBgPurple), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(155,48,255); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBgCyan), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(0,238,238); color: rgb(255,255,255);'>")
		return
	}

	if bytes.Compare([]byte(KAnsiBrightBgWhite), ansiValue) == 0 {
		el.resetCounter += 1
		htmlTag = []byte("<span style='background-color: rgb(255,255,255); color: rgb(255,255,255);'>")
		return
	}

	return
}

func TerminalToHtml(terminalString string) (htmlString string) {
	var terminalBytes = []byte(terminalString)
	var resetCounter = 0

	cursorUp := regexp.MustCompile("\\\\u001b\\[\\d+A")
	terminalBytes = cursorUp.ReplaceAllLiteral(terminalBytes, []byte(""))
	cursorDown := regexp.MustCompile("\\\\u001b\\[\\d+B")
	terminalBytes = cursorDown.ReplaceAllLiteral(terminalBytes, []byte(""))
	cursorRight := regexp.MustCompile("\\\\u001b\\[\\d+C")
	terminalBytes = cursorRight.ReplaceAllLiteral(terminalBytes, []byte(""))
	cursorLeft := regexp.MustCompile("\\\\u001b\\[\\d+D")
	terminalBytes = cursorLeft.ReplaceAllLiteral(terminalBytes, []byte(""))
	cursorColor8AndCommands := regexp.MustCompile("\\\\u001b\\[\\d+m")
	terminalBytes = cursorColor8AndCommands.ReplaceAllFunc(terminalBytes, func(in []byte) (out []byte) {
		inData := strings.TrimSpace(string(in))

		switch inData {
		case KAnsiReset:
			for i := 0; i != resetCounter; i += 1 {
				out = append(out, []byte("</span>")...)
			}
			resetCounter = 0

		case KAnsiBold:
			resetCounter += 1
			out = []byte("<span style='font-weight: bold;'>")
		case KAnsiUnderline:
			resetCounter += 1
			out = []byte("<span style='text-decoration: underline;'>")
		case KAnsiReversed:
			resetCounter += 1
			out = []byte("<span style='filter: invert(100%);'>")

		case KAnsiBlack:
			resetCounter += 1
			out = []byte("<span style='color: rgb(0,0,0);'>")
		case KAnsiRed:
			resetCounter += 1
			out = []byte("<span style='color: rgb(255,0,0);'>")
		case KAnsiGreen:
			resetCounter += 1
			out = []byte("<span style='color: rgb(0,255,0);'>")
		case KAnsiYellow:
			resetCounter += 1
			out = []byte("<span style='color: rgb(255,255,0);'>")
		case KAnsiBlue:
			resetCounter += 1
			out = []byte("<span style='color: rgb(0,0,255);'>")
		case KAnsiPurple:
			resetCounter += 1
			out = []byte("<span style='color: rgb(128,0,128);'>")
		case KAnsiCyan:
			resetCounter += 1
			out = []byte("<span style='color: rgb(0,255,255);'>")
		case KAnsiWhite:
			resetCounter += 1
			out = []byte("<span style='color: rgb(245,245,245);'>")

		case KAnsiBrightBlack:
			resetCounter += 1
			out = []byte("<span style='color: rgb(105,105,105);'>")
		case KAnsiBrightRed:
			resetCounter += 1
			out = []byte("<span style='color: rgb(205,51,51);'>")
		case KAnsiBrightGreen:
			resetCounter += 1
			out = []byte("<span style='color: rgb(127,255,0);'>")
		case KAnsiBrightYellow:
			resetCounter += 1
			out = []byte("<span style='color: rgb(255,255,127);'>")
		case KAnsiBrightBlue:
			resetCounter += 1
			out = []byte("<span style='color: rgb(30,144,255);'>")
		case KAnsiBrightPurple:
			resetCounter += 1
			out = []byte("<span style='color: rgb(155,48,255);'>")
		case KAnsiBrightCyan:
			resetCounter += 1
			out = []byte("<span style='color: rgb(0,238,238);'>")
		case KAnsiBrightWhite:
			resetCounter += 1
			out = []byte("<span style='color: rgb(255,255,255);'>")

		case KAnsiBgBlack:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(0,0,0); color: rgb(255,255,255);'>")
		case KAnsiBgRed:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(255,0,0); color: rgb(255,255,255);'>")
		case KAnsiBgGreen:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(0,255,0); color: rgb(255,255,255);'>")
		case KAnsiBgYellow:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(255,255,0); color: rgb(255,255,255);'>")
		case KAnsiBgBlue:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(0,0,255); color: rgb(255,255,255);'>")
		case KAnsiBgPurple:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(128,0,128); color: rgb(255,255,255);'>")
		case KAnsiBgCyan:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(0,255,255); color: rgb(255,255,255);'>")
		case KAnsiBgWhite:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(245,245,245); color: rgb(255,255,255);'>")

		case KAnsiBrightBgBlack:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(105,105,105); color: rgb(255,255,255);'>")
		case KAnsiBrightBgRed:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(205,51,51); color: rgb(255,255,255);'>")
		case KAnsiBrightBgGreen:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(127,255,0); color: rgb(255,255,255);'>")
		case KAnsiBrightBgYellow:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(255,255,127); color: rgb(255,255,255);'>")
		case KAnsiBrightBgBlue:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(30,144,255); color: rgb(255,255,255);'>")
		case KAnsiBrightBgPurple:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(155,48,255); color: rgb(255,255,255);'>")
		case KAnsiBrightBgCyan:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(0,238,238); color: rgb(255,255,255);'>")
		case KAnsiBrightBgWhite:
			resetCounter += 1
			out = []byte("<span style='background-color: rgb(255,255,255); color: rgb(255,255,255);'>")
		}
		return
	})
	//cursorColor16 := regexp.MustCompile("\\\\u001b\\[\\d+;\\d+m")
	//terminalBytes = cursorColor16.ReplaceAllLiteral(terminalBytes, []byte(""))
	//cursorColor256 := regexp.MustCompile("\\\\u001b\\[\\d+;\\d+;\\d+m")
	//terminalBytes = cursorColor256.ReplaceAllLiteral(terminalBytes, []byte(""))
	//cursorNav := regexp.MustCompile("\\\\u001b\\[\\d+D")
	//terminalBytes = cursorNav.ReplaceAllLiteral(terminalBytes, []byte(""))

	unicode := regexp.MustCompile("\\\\u.{4}")
	terminalBytes = unicode.ReplaceAllFunc(terminalBytes, func(in []byte) (out []byte) {
		//remove \u e deixa um número binário
		in = in[2:]
		number, _ := strconv.ParseInt(string(in), 16, 16)
		unicode := charmap.ISO8859_1.DecodeByte(byte(number))
		out = []byte(html.EscapeString(string(unicode)))
		return
	})

	terminalBytes = bytes.ReplaceAll(terminalBytes, []byte("\r\n"), []byte("<br>\n"))
	terminalString = string(terminalBytes)

	return terminalString
}
