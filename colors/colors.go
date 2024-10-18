package colors

const (
	tblack   = "\033[0;30m"
	tred     = "\033[0;31m"
	tgreen   = "\033[0;32m"
	tyellow  = "\033[0;33m"
	tblue    = "\033[0;34m"
	tmagenta = "\033[0;35m"
	tcyan    = "\033[0;36m"
	twhite   = "\033[0;37m"

	btblack   = "\033[1;30m"
	btred     = "\033[1;31m"
	btgreen   = "\033[1;32m"
	btyellow  = "\033[1;33m"
	btblue    = "\033[1;34m"
	btmagenta = "\033[1;35m"
	btcyan    = "\033[1;36m"
	btwhite   = "\033[1;37m"

	noColor = "\033[0m" // no color
)

func Black(s string) string {
	return tblack + s + noColor
}

func Red(s string) string {
	return tred + s + noColor
}

func Green(s string) string {
	return tgreen + s + noColor
}

func Yellow(s string) string {
	return tyellow + s + noColor
}

func Blue(s string) string {
	return tblue + s + noColor
}

func Magenta(s string) string {
	return tmagenta + s + noColor
}

func Cyan(s string) string {
	return tcyan + s + noColor
}

func White(s string) string {
	return twhite + s + noColor
}

func BBlack(s string) string {
	return btblack + s + noColor
}

func BRed(s string) string {
	return btred + s + noColor
}

func BGreen(s string) string {
	return btgreen + s + noColor
}

func BYellow(s string) string {
	return btyellow + s + noColor
}

func BBlue(s string) string {
	return btblue + s + noColor
}

func BMagenta(s string) string {
	return btmagenta + s + noColor
}

func BCyan(s string) string {
	return btcyan + s + noColor
}

func BWhite(s string) string {
	return btwhite + s + noColor
}
