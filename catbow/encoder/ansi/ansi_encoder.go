package encoder

type AnsiColorMode uint

const (
	TrueColor AnsiColorMode = iota
	Pallete256
	Pallete16
	Pallete8
)

type AnsiEncoderOptions struct {
	ColorMode AnsiColorMode
}

type AnsiEncoder struct {
	opts AnsiEncoderOptions
}

func NewDefaultAnsiEncoder() *AnsiEncoder {
	return &AnsiEncoder{
		AnsiEncoderOptions{
			ColorMode: TrueColor
		}
	}
}

func NewAnsiEncoder(opts AnsiEncoderOptions) {

}
