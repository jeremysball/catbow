package catbow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupRainbow() *RainbowStrategy {
	opts := *NewDefaultRainbowOptions()
	opts.Seed = 1
	return NewRainbow(opts)
}

func TestColorGeneration(t *testing.T) {

	rb := setupRainbow()
	outR := rb.ColorizeRune('r')
	outB := rb.ColorizeRune('b')

	assert.Contains(t, outR, "\033[38;2")
	assert.Contains(t, outR, "rm")

	assert.Contains(t, outB, "\033[38;2")
	assert.Contains(t, outB, "bm")

	assert.NotEqual(t, outR, rb.ColorizeRune('r'))
}

func TestNoColorGeneration(t *testing.T) {
	rb := setupRainbow()
	rb.opts.NoColor = true

	assert.Equal(t, string('r'), rb.ColorizeRune('r'))

}

func TestColorReset(t *testing.T) {

}
