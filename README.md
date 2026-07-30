<p align="center">
  <img width="845" height="352" alt="catbow rainbow text output" src="https://github.com/user-attachments/assets/462394ec-1da6-42a9-8bd5-e1cde4fdf9c6" />
</p>

<p align="center">
  <a href="https://github.com/jeremysball/catbow/releases/latest"><img src="https://img.shields.io/github/v/release/jeremysball/catbow?include_prereleases" alt="latest release"></a>
  <img src="https://img.shields.io/badge/go-1.25.1-00ADD8" alt="go version">
</p>

# catbow

A fast lolcat. I had `fortune | cowsay | lolcat` in my shell startup and it more
than doubled how long my terminal took to open, so I wrote one that does the
same thing in about 2ms instead of 50.

## Usage

```sh
fortune | cowsay | cb
```

## Install

```sh
go install github.com/jeremysball/catbow@latest
```

Or grab a binary from the [Releases](https://github.com/jeremysball/catbow/releases) tab.

## Flags

| Flag | Default | What it does |
|---|---|---|
| `-spread` | `1.05` | Stretches the rainbow vertically |
| `-freq` | `0.05` | Controls how quickly colors transition |
| `-seed` | `0` (random) | Picks what color the rainbow starts on |
| `-gen` | off | Generates sample text instead of reading stdin, for benchmarking |

## Benchmarks

Median of 5 runs, catbow vs. lolcat, same input:

| Case | catbow | lolcat | 
|---|---|---|
| shell startup (`echo asdf \| cb`) | 3.9ms | 108ms |
| 100 lines x 60 chars | 36ms | 206ms |
| 10k lines x 60 chars | 1.18s | 9.62s |
| 10k lines x 1000 chars | 18.9s | 169.8s |

Full numbers and how they were generated: [docs/BENCHMARKS.md](docs/BENCHMARKS.md).

## Build

```sh
make all      # clean, build, test
make build    # linux + windows binaries
make test     # go test ./catbow/
```

## Design

How the color math and the CLI/library split work: [docs/DESIGN.md](docs/DESIGN.md).

## License

No license file yet, so default copyright applies. Ask before reusing this
beyond personal use.
