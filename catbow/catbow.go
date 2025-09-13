package catbow

// Library code

// pseudocode for RGB calculation
// red   = sin(frequency * line_number + phase_shift_R)
// green = sin(frequency * line_number + phase_shift_G)
// blue  = sin(frequency * line_number + phase_shift_B)

/*
multiplied with the current line number when calculating RGB values in sin.
Increasing it will stretch the rainbow vertically
spread float64

controls our step size through color space. Decreasing it will make the same
amount of characters transition through less colors
frequency float64

random number added to line number, allows for variation between runs
seed int64

turns off all processing. input == output
noColor bool
*/
type Options struct {
	// Controls the horizontal width of each color band
	Spread float64
	// Rotates the rainbow
	Frequency float64
	// An offset for the starting color allowing varied but deterministic output
	Seed int64
}
