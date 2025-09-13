### Why
I included lolcat in a script that runs when I start my terminal `fortune | cowsay | lolcat`. It more than
doubled my startup time from 25ms to over 60ms (scales with number of lines). I want rainbow text but fast.
    
### MVP Requirements
Inputs: - stdin (default or when we see `-`)
        - files
        - if none are given we exit and print usage instructions  
        when stdin is given the application WILL accept input 
Outputs: escaped text to stdout 
        - [24 bit true color] (I will verify but I believe this is what Go prefers)
        - plain ascii (--no-color)
### High-Level Architecture
CLI:
       - if given an argument of a readable, accesible file will read and colorize
         contents
       - if no file provided will read from stdin
Application Layer (main package):
        - parses arguments and options
        - detect color mode 
        - orchestrates Colorizer
        - uses standard Go interfaces Reader and Writer to decouple API Layer
         from Application Layer
catbow package:
  ColorMath (internal):
        - Responsible for the actual math and production of colored strings
        - Pure functions
        - func rainbow(freq, i float64) (r, g, b uint)
  Colorizer:
        - Stateful 
        - Has an instance of Options struct (contains options required for the operation
          of the Colorizer)
        - func NewColorizer(Options) *Colorizer
        - func Colorize(w io.Writer, r io.Reader) error
        - Uses ColorMath to produce escaped strings
        - Can be extended create various types of colorizers beyond rainbow

