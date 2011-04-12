#!/usr/bin/env python

import sys
import os
import re


UNARY_OPERATORS = {
    '+': 'p',
    '-': 'm',
}
RESTORE_UNARY = dict(zip(UNARY_OPERATORS.itervalues(), UNARY_OPERATORS.iterkeys()))

OPERATOR_PRIORITIES = {
    'm':  (3, 'r'),
    '**': (2, 'r'),
    '*':  (1, 'l'),
    '/':  (1, 'l'),
    '+':  (0, 'l'),
    '-':  (0, 'l'),
}

CHAR2FUNC = {
    '**': lambda a, b: a ** b,
    '*':  lambda a, b: a * b,
    '/':  lambda a, b: a / b,
    '+':  lambda a, b: a + b,
    '-':  lambda a, b: a - b,
}

TOKEN_RE = re.compile(
    '''
    (\s*(
         \d+\.\d+ |
         \d+      |
         \*\*     |
         [/*+-]   |
         \(       |
         \)
        )
    )''', re.VERBOSE)

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


def isop(token):
    return token in OPERATOR_PRIORITIES


def priority(operator):
    return OPERATOR_PRIORITIES[operator][0]


def isleft(operator):
    return OPERATOR_PRIORITIES[operator][1] == 'l'


def main(argv):
    expression = (' '.join(argv) if argv else raw_input()).strip()
    if not expression:
        return 0

    stack, postfix = [], []
    prev = '+' # dummy operator, only for first iteration
    for token in tokens(expression):

        if isop(token):
            if isop(prev) and token in UNARY_OPERATORS:
                if token == '+':
                    continue
                postfix.append(0)
                token = UNARY_OPERATORS[token]

            while (stack and isop(stack[-1]) and
                   (isleft(token) and priority(token) <= priority(stack[-1]) or
                    not isleft(token) and priority(token) < priority(stack[-1]))
            ):
                postfix.append(stack.pop())
            stack.append(token)

        elif token == '(':
            stack.append(token)

        elif token == ')':
            while stack[-1] != '(':
                postfix.append(stack.pop())
                if not stack:
                    raise DError('parentheses mismatch')
            stack.pop()

        else:
            postfix.append(float(token))

        prev = token

    while stack:
        operator = stack.pop()
        if operator == '(':
            raise DError('parentheses mismatch')
        postfix.append(operator)

    result = calculate(map(lambda x: RESTORE_UNARY.get(x, x), postfix))
    print int(result) if result.is_integer() else result
    return 0


if __name__ == '__main__':
    try:
        sys.exit(main(sys.argv[1:]))
    except KeyboardInterrupt:
        print 'Interrupted by user'
    except DError, ex:
        print u'%s: error: %s' % (os.path.basename(__file__), ex.args[0])
    sys.exit(1)
