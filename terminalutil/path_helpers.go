package terminalutil

import (
	"fmt"
)

var enableColors = false

func EnableColors(enable bool) {
	enableColors = enable
}

func FormatBold(str string) string {
	if enableColors {
		return "\033[1m" + str + "\033[0m"
	} else {
		return str
	}
}

func FormatRed(str string) string {
	if enableColors {
		return "\033[31m" + str + "\033[0m"
	} else {
		return str
	}
}

func FormatGreen(str string) string {
	if enableColors {
		return "\033[32m" + str + "\033[0m"
	} else {
		return str
	}
}

func FormatYellow(str string) string {
	if enableColors {
		return "\033[33m" + str + "\033[0m"
	} else {
		return str
	}
}

func FormatBlue(str string) string {
	if enableColors {
		return "\033[34m" + str + "\033[0m"
	} else {
		return str
	}
}

func PrintError(str string, err error) {
	msg := fmt.Sprintf("%s: %s", str, err)
	fmt.Println((FormatRed(msg)))
}
