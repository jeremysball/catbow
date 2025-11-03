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
	"github.com/jeremysball/catbow/catbow/encoder/ansi"
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

	var shouldGenerateFlag bool
	var freqFlag float64
	var spreadFlag float64
	var seedFlag int
	var genLineLenFlag int
	var genNumLinesFlag int

	rainbowOpts := catbow.NewRainbowOptions()

	// these defaults SHOULD come from the Strategy itself
	flag.BoolVar(
		&shouldGenerateFlag,
		"gen",
		false,
		"Enable generating random text to colorize")
	flag.IntVar(&seedFlag,
		"seed",
		rainbowOpts.Seed,
		"Changes what color the rainbow starts on. 0 == random")
	flag.Float64Var(&spreadFlag,
		"spread",
		rainbowOpts.Spread,
		"Rotates the rainbow")
	flag.Float64Var(&freqFlag,
		"freq",
		rainbowOpts.Frequency,
		"Controls the horizontal width of each color band")
	flag.IntVar(&genLineLenFlag, "gen-line-width", 80, "")
	flag.IntVar(&genNumLinesFlag, "gen-num-lines", 256, "")

	flag.Parse()

	w := io.Writer(os.Stdout)
	if shouldGenerateFlag {
		r = newMockReader(genLineLenFlag, genNumLinesFlag)
	} else {
		r = bufio.NewReader(os.Stdin)
	}

	if seedFlag == 0 {
		// just picked a number here - the only thing that
		// matters it that it doesn't become MASSIVE and overflow
		// the color calculation
		rainbowOpts.Seed = rand.IntN(65535)
	} else {
		rainbowOpts.Seed = seedFlag
	}
	rainbowOpts.Spread = spreadFlag
	rainbowOpts.Frequency = freqFlag

	colorizer := catbow.NewColorizer(catbow.NewRainbowStrategy(rainbowOpts))
	err := colorizer.Colorize(r, w)

	fmt.Print(ansi.Reset)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
