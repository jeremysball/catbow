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

func newMockReader(genLineLen, genNumLines, seed int) *bufio.Reader {
	cmd := exec.Command(
		"./generate_text.py",
		fmt.Sprintf("--line-width=%d", genLineLen),
		fmt.Sprintf("--num-lines=%d", genNumLines),
		fmt.Sprintf("--seed=%d", seed))
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

	rainbowOpts := catbow.NewRainbowOptions()

	var generatorOpts = struct {
		seed       int
		numLines   int
		lineLength int
	}{
		0,
		80,
		256,
	}
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
	flag.IntVar(&generatorOpts.lineLength, "gen-line-length", generatorOpts.lineLength, "")
	flag.IntVar(&generatorOpts.numLines, "gen-num-lines", generatorOpts.numLines, "")
	flag.IntVar(&generatorOpts.seed, "gen-seed", generatorOpts.seed, "")

	flag.Parse()

	w := io.Writer(os.Stdout)
	if shouldGenerateFlag {
		r = newMockReader(generatorOpts.lineLength, generatorOpts.numLines, generatorOpts.seed)
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
