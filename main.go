package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"strings"

	"github.com/jeremysball/catbow/catbow"
)

func newMockReader(genLineLen, genNumLines int) *bufio.Reader {
	cmd := exec.Command(
		"./generate_text.py",
		fmt.Sprintf("--line-width=%d", genLineLen),
		fmt.Sprintf("--num-lines=%d", genNumLines))
	text, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	r := bufio.NewReader(strings.NewReader(string(text)))
	return r
}

// simplest runner to test Cleanup()
func main() {

	var r *bufio.Reader

	var shouldGenerate bool
	var freq float64
	var spread float64
	var seed int
	var genLineLen int
	var genNumLines int

	// these defaults SHOULD come from the Strategy itself
	flag.BoolVar(
		&shouldGenerate,
		"gen",
		false,
		"Enable generating random text to colorize")
	flag.IntVar(&seed,
		"seed",
		0,
		"Changes what color the rainbow starts on. 0 == random")
	flag.Float64Var(&spread,
		"spread",
		1.05,
		"Rotates the rainbow")
	flag.Float64Var(&freq,
		"freq",
		0.05,
		"Controls the horizontal width of each color band")
	flag.IntVar(&genLineLen, "gen-line-width", 80, "")
	flag.IntVar(&genNumLines, "gen-num-lines", 256, "")

	flag.Parse()

	w := io.Writer(os.Stdout)
	if shouldGenerate {
		r = newMockReader(genLineLen, genNumLines)
	} else {
		r = bufio.NewReader(os.Stdin)
	}

	opts := catbow.NewRainbowOptions()
	if seed == 0 {
		// just picked a number here - the only thing that
		// matters it that it doesn't become MASSIVE and overflow
		// the color calculation
		opts.Seed = rand.IntN(65535)
	} else {
		opts.Seed = seed
	}
	opts.Spread = spread
	opts.Frequency = freq

	colorizer := catbow.NewColorizer(catbow.NewRainbowStrategy(opts))
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
