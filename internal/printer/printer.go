package printer

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func Fatal(line string) {
	color.Red(line)
	os.Exit(1)
}

func Fatalf(line string, args ...interface{}) {
	Fatal(fmt.Sprintf(line, args...))
}

func Info(line string) {
	color.Green(line)
}

func Infof(line string, args ...interface{}) {
	Info(fmt.Sprintf(line, args...))
}

func NoColor(line string) {
	fmt.Println(line)
}

func NoColorf(line string, args ...interface{}) {
	NoColor(fmt.Sprintf(line, args...))
}

func Warn(line string) {
	color.Yellow(line)
}

func Warnf(line string, args ...interface{}) {
	Warn(fmt.Sprintf(line, args...))
}
