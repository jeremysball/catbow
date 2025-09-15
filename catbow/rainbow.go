package catbow

import (
	"fmt"
	"math"

	"github.com/jeremysball/catbow/catbow/encoder/ansi"
)

type rainbowOptions struct {
	// Spread is a divisor for the offset (position in line). In other words
	// it controls how fast the color advances through the rainbow as we move left
	// to right within the line
	Spread float64
	// This has a similar effect to spread. The frequency of the sine wave used to calculate the colors.
	Frequency float64
	// An offset for the starting color, allowing varied but deterministic output
	Seed int
	// Disables catbow, input will equal output
	NoColor    bool
	redShift   float64
	greenShift float64
	blueShift  float64
}

func NewRainbowOptions() *rainbowOptions {
	return &rainbowOptions{
		Spread:     3.0,
		Frequency:  0.1,
		Seed:       0,
		NoColor:    false,
		redShift:   0,
		greenShift: (2 * math.Pi) / 3,
		blueShift:  (4 * math.Pi) / 3,
	}
}

/*
	rainbowStrategy is stateful and therefore not re-usable. To colorize a new

stream, create another strategy.
*/
type rainbowStrategy struct {
	Opts                rainbowOptions
	offset              float64
	prevLineStartOffset float64
	wasLastRuneNewLine  bool
}

// extract to ColorEncoder
func (rb *rainbowStrategy) CleanupStr() string {
	return ansi.Reset
}

func NewRainbowStrategy(opts *rainbowOptions) *rainbowStrategy {
	s := &rainbowStrategy{
		Opts: *opts,
	}
	s.offset = float64(opts.Seed) + 1.0
	s.prevLineStartOffset = s.offset
	s.wasLastRuneNewLine = false

	return s
}

func (rb *rainbowStrategy) calculateRainbow(offset float64) rgbColor {

	freq := rb.Opts.Frequency

	scaledOffset := offset / rb.Opts.Spread

	// math.Sin(freq...)*127 + 128 maps sine (-1 to 1) to a number between 1 and 255
	red := math.Sin((freq*scaledOffset)+rb.Opts.redShift)*127 + 128
	green := math.Sin((freq*scaledOffset)+rb.Opts.greenShift)*127 + 128
	blue := math.Sin((freq*scaledOffset)+rb.Opts.blueShift)*127 + 128

	return rgbColor{
		R: uint8(math.Round(red)),
		G: uint8(math.Round(green)),
		B: uint8(math.Round(blue)),
	}
}

func (rb *rainbowStrategy) colorizeRune(r rune) string {
	/* TODO: Refactor match into a call to calculateRainbow
	and a call to the injected ColorFormatter which does what
	the fmt.Sprintf() call does but allows us to be agnostic as
	to what we're outputting to. Essentially this becomes the
	API for Colorizers to call

	*/
	if rb.Opts.NoColor {
		return string(r)
	}

	// this is again to deal with the prefix nature of the lolcat code

	// since we're doing essentially a postfix operation instead of
	// lolcat's prefix increment we add 1 to the Seed to derive the
	// starting offset when creating the strategy
	if r == '\n' {
		rb.wasLastRuneNewLine = true
	}

	// what about mutliple newlines in a row?
	if rb.wasLastRuneNewLine {
		rb.wasLastRuneNewLine = false
		rb.prevLineStartOffset += 1
		rb.offset = rb.prevLineStartOffset
	} else {
		/*
			A small deviation from lolcat. Offset
			accumulates the spread:

			off += (1 / Spread) instead of
			off = (charIndex / Spread).

			Importantly this allows for the offset to be testable
		*/

		rb.offset = rb.offset + (1 / rb.Opts.Spread)
	}

	rgb := rb.calculateRainbow(rb.offset)

	return fmt.Sprintf(
		ansi.Esc+"[38;2;%d;%d;%dm%c",
		rgb.R,
		rgb.G,
		rgb.B,
		r)
}
