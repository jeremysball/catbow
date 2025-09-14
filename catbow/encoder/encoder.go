package encoder

type Resetter interface {
	Reset() string
}

type ColorEncoder interface {
	FormatRgb(r, g, b int) string
}
