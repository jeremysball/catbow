package catbow

import (
	"github.com/lordxarus/catbow/catbow/encoder/ansi"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupRainbow() *RainbowStrategy {
	opts := *NewDefaultRainbowOptions()
	opts.Seed = 1
	return NewRainbowStrategy(opts)
}

func TestColorGeneration(t *testing.T) {

	rb := setupRainbow()
	outR := rb.colorizeRune('r')
	outB := rb.colorizeRune('b')

	assert.Contains(t, outR, ansi.Esc+"[38;2")
	assert.Contains(t, outR, "rm")

	assert.Contains(t, outB, ansi.Esc+"[38;2")
	assert.Contains(t, outB, "bm")

	assert.NotEqual(t, outR, rb.colorizeRune('r'))
}

func TestNoColorGeneration(t *testing.T) {
	rb := setupRainbow()
	rb.opts.NoColor = true

	assert.Equal(t, string('r'), rb.colorizeRune('r'))

}

func TestColorReset(t *testing.T) {
	rb := setupRainbow()
	out := rb.colorizeRune('a')
	assert.NotContains(t, out, ansi.Reset)
	out = rb.Cleanup()
	assert.Contains(t, out, ansi.Reset)
}
