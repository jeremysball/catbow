#!/usr/bin/env python3
# generate test files
import argparse
import math
import random
from pathlib import Path
from typing import Generator


def make_sentences(
    num_lines: int, line_length: int, min_char: int, max_char: int
) -> list[str]:
    def make_chars(n_char: int) -> Generator[str]:
        for _ in range(n_char):
            code = math.floor(random.random() * max_char) % max_char
            yield chr(code + min_char) if code < min_char else chr(code)

    def make_sentences(n_sentences: int) -> list[str]:
        sentences = []
        for _ in range(n_sentences):
            sentence = "".join([c for c in make_chars(line_length)])
            sentences.append(sentence)
        return sentences

    return make_sentences(num_lines)


def main():
    parser = argparse.ArgumentParser(
        prog="generate",
        description="Fill files with characters for testing",
    )

    parser.add_argument("--file", required=False)
    parser.add_argument(
        "--line-length", type=int, default=80, required=False, help="default: 80"
    )
    parser.add_argument(
        "--num-lines", type=int, default=512, required=False, help="default: 512"
    )
    parser.add_argument(
        "--max",
        default=128,
        type=int,
        required=False,
        help="max ascii character code to write to file. default: 128",
    )
    parser.add_argument(
        "--min",
        default=33,
        type=int,
        required=False,
        help="min ascii character code to write to file. default: 33",
    )

    args = parser.parse_args()

    sentences = make_sentences(
        num_lines=args.num_lines,
        line_length=args.line_length,
        min_char=args.min,
        max_char=args.max,
    )
    if args.file is not None:
        p = Path(args.file)
        if p.exists():
            p.unlink()
        with p.open("a") as f:
            for s in sentences:
                f.write(f"{s}\n\r")
    else:
        print("\n".join(sentences))


if __name__ == "__main__":
    main()
