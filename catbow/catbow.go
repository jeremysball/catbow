package catbow

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

type Cleanupper interface {
	Cleanup() string
}

type ColorStrategy interface {
	colorizeRune(r rune) string
}

type Colorizer struct {
	// TODO: this can be unexported once encoder is implemented
	Strategy ColorStrategy
}

func NewColorizer(c ColorStrategy) *Colorizer {
	return &Colorizer{
		Strategy: c,
	}
}

// this function is concerned with reading input from r,
// running whatever APIs needed to get the data to write to w
func (c *Colorizer) Colorize(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanRunes)
	for {
		if !s.Scan() {
			return s.Err()
		}

		r, size := utf8.DecodeRune(s.Bytes())

		if r == utf8.RuneError {
			switch size {
			case 0:
				return errors.New("RuneError: DecodeRune got passed an empty buffer")
			case 1:
				return errors.New("colorize: invalid encoding for bytes % +q")
			}
		}

		coloredString := c.Strategy.colorizeRune(r)
		_, err := w.Write([]byte(coloredString))
		if err != nil {
			return fmt.Errorf("colorize: error while writing to stdout: %w", err)
		}
	}

}
