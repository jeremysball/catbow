package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/lordxarus/catbow/catbow"
)

func newMockReader() *bufio.Reader {
	cmd := exec.Command("./generate_text.py")
	text, err := cmd.Output()
	println(len(text))
	if err != nil {
		fmt.Println(err)
	}
	r := bufio.NewReader(strings.NewReader(string(text)))
	return r
}

// simplest runner to test Cleanup()
func main() {

	var r *bufio.Reader

	if slices.Contains(os.Args, "-gen") {
		r = newMockReader()
	} else {
		r = bufio.NewReader(os.Stdin)
	}
	w := io.Writer(os.Stdout)
	colorizer := catbow.NewColorizer(catbow.NewRainbowStrategyDefaultOpts())
	err := colorizer.Colorize(r, w)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	/* this will be replaced by
	AnsiFormatter satisfies Formatter  interface
	colFmt := AnsiFormatter()
	if colFmt.(catbow.Cleanupper) {
		w.colFmt.Cleanup()
	}
	*/
	if cleaner, ok := colorizer.Strategy.(catbow.Cleanupper); ok {
		_, err := w.Write([]byte(cleaner.Cleanup()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}
}
