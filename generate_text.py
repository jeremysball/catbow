#!/usr/bin/env python3
# generate test files
import argparse
from collections.abc import Callable
import math
import random
from pathlib import Path
from typing import Generator

DEFAULT_MIN_CHAR = 34
DEFAULT_MAX_CHAR = 127


def make_chars(
    n_char: int, min_char: int, max_char: int, rand_fn: Callable[[int, int], int]
) -> Generator[str]:
    for _ in range(n_char):
        code = int(random.random() * max_char) % max_char
        code += min_char
        yield chr(code)
        # yield chr(min_char) if code < min_char else chr(code)


def make_sentence(
    line_width: int, min_char: int, max_char: int, rand_fn: Callable[[int, int], int]
) -> str:
    return "".join(
        [
            c
            for c in make_chars(
                line_width, min_char=min_char, max_char=max_char, rand_fn=rand_fn
            )
        ]
    )


def make_sentences(
    num_lines: int,
    line_width: int,
    min_char: int,
    max_char: int,
    rand_fn: Callable[[int, int], int],
) -> Generator[str]:
    for _ in range(num_lines):
        yield make_sentence(line_width, min_char, max_char, rand_fn)


def main():
    parser = argparse.ArgumentParser(
        prog="generate",
        description="Fill files with characters for testing",
    )

    parser.add_argument("--file", required=False)
    parser.add_argument(
        "--line-width", type=int, default=80, required=False, help="default: 80"
    )
    parser.add_argument(
        "--num-lines", type=int, default=512, required=False, help="default: 512"
    )
    parser.add_argument(
        "--max",
        default=DEFAULT_MAX_CHAR,
        type=int,
        required=False,
        help="max ascii character code to write to file. default: DEFAULT_MAX_CHAR",
    )
    parser.add_argument(
        "--min",
        default=DEFAULT_MIN_CHAR,
        type=int,
        required=False,
        help="min ascii character code to write to file. default: DEFAULT_MIN_CHAR",
    )

    parser.add_argument(
        "--seed",
        default=0,
        type=int,
        required=False,
        help="seed for randomness. zero == random",
    )

    args = parser.parse_args()

    sentences = make_sentences(
        num_lines=args.num_lines,
        line_width=args.line_width,
        min_char=args.min,
        max_char=args.max,
        rand_fn=(
            random.Random().randint
            if args.seed == 0
            else random.Random(args.seed).randint
        ),
    )

    if args.file is not None:
        p = Path(args.file)
        if p.exists():
            p.unlink()
        with p.open("a") as f:
            percent = 0.0
            for i, s in enumerate(sentences):
                if (i + 1) % (args.num_lines // 10) == 0:
                    percent += 10.0
                    print(f"{percent}% done")

                f.write(s + "\n")
    else:
        for s in sentences:
            print(s)


if __name__ == "__main__":
    main()
