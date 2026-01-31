#!/bin/python
# 'call' file can have any extension, however it must be an executable!
# any nested functions should be obtained by reading from stdin
# any output must be written to stdout.

import sys

print(f"function {sys.argv[1]} with storage minecraft:example input")

for line in sys.stdin:
    print(line, end="")
