package catbow

import (
	"fmt"
	"testing"

	"github.com/jeremysball/catbow/catbow/encoder/ansi"
	"github.com/stretchr/testify/assert"
)

func setupTestRainbow() *rainbowStrategy {
	opts := NewRainbowOptions()
	opts.Seed = 1
	opts.Spread = 3
	opts.Frequency = .1
	return NewRainbowStrategy(opts)
}

func TestColorControlCodes(t *testing.T) {

	rb := setupTestRainbow()
	outR := rb.colorizeRune('r')
	outB := rb.colorizeRune('b')

	assert.Contains(t, outR, ansi.Esc+"[38;2")
	assert.Contains(t, outR, "mr")

	assert.Contains(t, outB, ansi.Esc+"[38;2")
	assert.Contains(t, outB, "mb")

	// increase the offset to avoid collisions
	for range 20 {
		rb.colorizeRune('r')
	}
	assert.NotEqual(t, outR, rb.colorizeRune('r'))
}

func TestNoColorGeneration(t *testing.T) {
	rb := setupTestRainbow()
	rb.Opts.NoColor = true

	assert.Equal(t, string('r'), rb.colorizeRune('r'))

}

func TestColorReset(t *testing.T) {
	rb := setupTestRainbow()
	out := rb.colorizeRune('a')
	assert.NotContains(t, out, ansi.Esc+"[0m")
	out = rb.CleanupStr()
	assert.Equal(t, out, ansi.Esc+"[0m")
}

func TestRainbowAlgorithm(t *testing.T) {
	rb := setupTestRainbow()

	defer fmt.Println(rb.CleanupStr())

	assert.Equal(t, 2.0, rb.offset)
	rgb := rb.calculateRainbow(rb.offset)
	assert.Equal(t, rgbColor{136, 234, 14}, rgb)

	rb.offset += 10
	rgb = rb.calculateRainbow(rb.offset)
	assert.Equal(t, rgbColor{177, 205, 2}, rgb)
}

func TestOffsetProgression(t *testing.T) {
	/*
		Explanation of the lolcat offset:
	*/
	rb := setupTestRainbow()
	defer fmt.Println(rb.CleanupStr())

	// offset gets initialized to seed (1 if using setupRainbow() test setup)
	// in lolcat the input is loaded into memory and split on lines. before the first
	// iteration the offset is incremented so we start with seed + 1

	assert.Equal(t, 2.0, rb.offset)
	rb.colorizeRune('a')
	assert.InDelta(t, 7.0/3.0, rb.offset, .001)
	assert.Equal(t, 2.3333333333333335, rb.offset)
	rb.colorizeRune('s')
	assert.InDelta(t, 8.0/3.0, rb.offset, .001)
	rb.colorizeRune('d')
	assert.InDelta(t, 9.0/3.0, rb.offset, .001)
	rb.colorizeRune('f')
	assert.InDelta(t, 10.0/3.0, rb.offset, .001)
	rb.colorizeRune('\n')

	// rb.offset is reset to what it was when entering the
	// line loop in this case it gets set back to 2.0

	// then incremented again before the next iteration of the line loop
	assert.Equal(t, 3.0, rb.offset)

	// EOF

}
