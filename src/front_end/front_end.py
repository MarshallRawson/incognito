#!/usr/bin/env python3

import sys

# send message from python front end to c++ backend
print("hello world!")

# get a message from the c++ backend
a = input ('')

# print the message on the terminal
print(a, file=sys.stderr);
