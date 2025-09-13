# Catbow Design
### Why
I included lolcat in a script that runs when I start my terminal `fortune | cowsay | lolcat`. It more than
doubled my startup time from 25ms to over 60ms (scales with 
number of lines). I want rainbow text but fast.
### Core Contract
---
#### Inputs
- files
- stdin (when no files given, or when we see `-`)

##### Output (colored text to stdout)
- 24 bit true color escape codes

### Architecture
---
#### CLI
- if given an argument of a readable, accessible file CLI will read, colorize, and print contents
- if no file provided will read from stdin
##### flags:
    --spread (will stretch the rainbow vertically, default: 1.0)
    --freq (controls how quickly colors transition, default: 0.1)
    --seed (controls the random seed, zero is random default: 0)
    --no-color (prints text exactly as it came in)

#### main package (Application Layer):
- parses arguments and options
- detect color mode
- creates and passes `Options` to `Colorizer`
- calls `Colorizer`
- uses standard Go interfaces `Reader` and `Writer` to decouple library from main
#### catbow package:
##### `ColorMath`:
- Responsible for the actual math and production of colored strings
- Pure functions
- `func rainbow(freq, i float64) (r, g, b uint)`

Colorizer:
- Stateful 
- Uses ColorMath to produce escaped strings
- Reads and writes to injected `Reader` and `Writer`
- `Options` struct (contains options required for the operation of the Colorizer)
- `func NewColorizer(Options) *Colorizer`
- `func Colorize(w io.Writer, r io.Reader) error`

