#!/usr/bin/env python3
# generate test files
import argparse
import math
import random
from pathlib import Path

MAX_ASCII_CODE = 90
MIN_ASCII_CODE = 65
PATH = "small-rand-short.txt"


def make_data(num_lines: int = 512, line_length: int = 80) -> list[str]:
    sentences = []
    for _ in range(0, num_lines):
        sentence = []
        for _ in range(0, line_length):
            code = math.floor(random.random() * MAX_ASCII_CODE) % MAX_ASCII_CODE
            if code < MIN_ASCII_CODE:
                code += MIN_ASCII_CODE
            sentence.append(chr(code))
        sentences.append("".join(sentence))
    return sentences


def main():
    parser = argparse.ArgumentParser(
        prog="generate",
        description="Fill files with characters for testing",
    )

    parser.add_argument("-f, --file-path", required=False)
    parser.add_argument(
        "--line-length", type=int, default=80, required=False, help="default: 80"
    )
    parser.add_argument(
        "--num-lines", type=int, default=512, required=False, help="default: 512"
    )
    parser.add_argument(
        "--max",
        default=33,
        type=int,
        required=False,
        help="max ascii character code to write to file. default: 33",
    )
    parser.add_argument(
        "--min",
        default=128,
        type=int,
        required=False,
        help="min ascii character code to write to file. default: 128",
    )

    args = parser.parse_args()
    ln_len = args.line_length
    n_lines = args.num_lines

    sentences = make_data(num_lines=n_lines, line_length=ln_len)
    output = "\n".join(sentences)

    if hasattr(args, "file_path"):
        p = Path(args.file_path)
        if p.exists():
            p.unlink()
        p.write_text(output)
    else:
        print(output)


if __name__ == "__main__":
    main()
