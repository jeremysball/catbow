package catbow

import (
	"bufio"
	"io"
)

type ColorAlgorithm interface {
	ColorizeRune(r rune) string
}

type Colorizer struct {
	algo ColorAlgorithm
}

func NewColorizer(c ColorAlgorithm) *Colorizer {
	return &Colorizer{
		algo: c,
	}
}

// this function is concerned with reading input from r,
// running whatever APIs needed to get the data to write to w
func (c *Colorizer) Colorize(r io.Reader, w io.Writer) error {
	rw := bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))
	for {
		r, _, readErr := rw.Reader.ReadRune()

		if readErr != nil {
			flushErr := rw.Writer.Flush()
			if flushErr != nil {
				return flushErr
			}

			if readErr == io.EOF {
				return nil
			} else {
				return readErr
			}
		}

		_, err := rw.Writer.WriteString((c.algo.ColorizeRune(r)))
		if err != nil {
			return err
		}
	}

}
