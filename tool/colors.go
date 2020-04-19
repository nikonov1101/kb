package tool

const (
	RED = "\033[1;31m"
	YEL = "\033[1;33m"
	GRE = "\033[1;32m"
	NC  = "\033[0m" // No Color
)

func green(s string) string {
	return GRE + s + NC
}

func yellow(s string) string {
	return YEL + s + NC
}

func red(s string) string {
	return RED + s + NC
}

