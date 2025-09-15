package catbow

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/lordxarus/catbow/catbow/encoder/ansi"
)

// I think these files that implement ColorAlgorithms should be self contained
// If we make Rainbow manage offset then I think that's the only way to maintain
// consistency. So that means this code will just move the cursor on a call to
// ColorizeString.
type RainbowOptions struct {
	// Controls the horizontal width of each color band
	Spread float64
	// Rotates the rainbow
	Frequency float64
	// An offset for the starting color allowing varied but deterministic output
	Seed int64
	// Disables catbow, input will equal output
	NoColor bool
}

func NewDefaultRainbowOptions() *RainbowOptions {
	return &RainbowOptions{
		Spread:    3.0,
		Frequency: 0.1,
		Seed:      int64(rand.Int()),
		NoColor:   false,
	}
}

type RainbowStrategy struct {
	opts       RainbowOptions
	cursor     int64
	redShift   float64
	greenShift float64
	blueShift  float64
}

func (rb *RainbowStrategy) Cleanup() string {
	return ansi.Reset
}

func NewRainbowStrategy(opts RainbowOptions) *RainbowStrategy {
	return &RainbowStrategy{
		opts:       opts,
		cursor:     0,
		redShift:   0,
		greenShift: 2 * math.Pi / 3,
		blueShift:  4 * math.Pi / 3,
	}
}

func NewRainbowStrategyDefaultOpts() *RainbowStrategy {
	return NewRainbowStrategy(*NewDefaultRainbowOptions())
}

/*
		 def self.rainbow(freq, i)
			red   = Math.sin(freq*i + 0) * 127 + 128
			green = Math.sin(freq*i + 2*Math::PI/3) * 127 + 128
			blue  = Math.sin(freq*i + 4*Math::PI/3) * 127 + 128
			"#%02X%02X%02X" % [ red, green, blue ]
	    end
*/
func (rb *RainbowStrategy) colorizeRune(r rune) string {
	if rb.opts.NoColor {
		return string(r)
	}

	freq := rb.opts.Spread

	// might want to store cursor and seed as floats
	seed := float64(rb.opts.Seed)
	cursor := float64(rb.cursor)

	red := math.Sin(freq*cursor+rb.redShift+seed)*127 + 128
	green := math.Sin(freq*cursor+rb.greenShift+seed)*127 + 128
	blue := math.Sin(freq*cursor+rb.blueShift+seed)*127 + 128

	rb.cursor += 1

	return fmt.Sprintf(
		ansi.Esc+"[38;2;%d;%d;%dm%c",
		int(math.Floor(red)),
		int(math.Floor(green)),
		int(math.Floor(blue)),
		r)
}
