#!/usr/bin/env python

import sys
import os
import re


OPERATOR_PRIORITIES = {
    '*': 1,
    '/': 1,
    '+': 0,
    '-': 0,
}

CHAR2FUNC = {
    '*': lambda a, b: a * b,
    '/': lambda a, b: a / b,
    '+': lambda a, b: a + b,
    '-': lambda a, b: a - b,
}

TOKEN_RE = re.compile(r'(\s*(\d+\.\d+|\d+|\+|-|\*|\\|\^))')

class DError(Exception): pass


def tokens(string):
    pos = 0
    while pos < len(string):
        match = TOKEN_RE.match(string[pos:])
        if match:
            whole, value = match.groups()
            pos += len(whole)
            yield value
        else:
            while string[pos].isspace():
                pos += 1
            raise DError('unknown char: %s' % string[pos])


def calculate(postfix):
    result = []
    while postfix:
        token = postfix.pop(0)
        if isop(token):
            args = [result.pop(), result.pop()][::-1]
            result.append(CHAR2FUNC[token](*args))
        else:
            result.append(token)
    return result[0]


def make_number(string):
    if string.isdigit():
        return int(string)
    return float(string)


def isop(token):
    return token in OPERATOR_PRIORITIES


def priority(operator):
    return OPERATOR_PRIORITIES[operator][0]


def main(argv):
    expression = (' '.join(argv) if argv else raw_input()).strip()
    operators, postfix = [], []
    for token in tokens(expression):
        if isop(token):
            while (operators and isop(operators[-1]) and
                   priority(token) <= priority(operators[-1])
            ):
                postfix.append(operators.pop())
            operators.append(token)
        else:
            postfix.append(make_number(token))
    postfix.extend(reversed(operators))
    result = calculate(postfix)
    correct = eval(expression.replace('^', '**'))
    print result
    if result != correct:
        print 'correct: {0}'.format(correct)


if __name__ == '__main__':
    try:
        sys.exit(main(sys.argv[1:]))
    except KeyboardInterrupt:
        print 'Interrupted by user'
    except DError, ex:
        print u'%s: error: %s' % (os.path.basename(__file__), ex.args[0])
    sys.exit(1)
