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
    'p':  (3, 'r'),
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

ALLOWED_SEQUENCES = {
    lambda t: isop(t) or t == '(' or t is None:
        lambda c: isnumber(c) or isunary(c) or c == '(',
    lambda t: isnumber(t) or t == ')':
        lambda c: isop(c) or c == ')',
}

class DijkstraError(Exception): pass


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
            raise DijkstraError('unexpected char: %s' % string[pos])


def calculate(postfix):
    result = []
    while postfix:
        token = postfix.pop(0)
        if isop(token):
            args = [result.pop(), result.pop()][::-1]
            try:
                result.append(CHAR2FUNC[token](*args))
            except (ValueError, OverflowError), ex:
                raise DijkstraError(ex.args[-1])
            except ZeroDivisionError, ex:
                msg = ex.args[0]
                if msg == 'float division':
                    msg = 'division by zero'
                raise DijkstraError(msg)
        else:
            result.append(token)
    return result[0]


def isop(token):
    return token in OPERATOR_PRIORITIES


def priority(operator):
    return OPERATOR_PRIORITIES[operator][0]


def isleft(operator):
    return OPERATOR_PRIORITIES[operator][1] == 'l'


def isunary(operator):
    return operator in UNARY_OPERATORS


def isnumber(token):
    try:
        float(token)
    except (TypeError, ValueError):
        return False
    return True


def main(argv):
    expression = (' '.join(argv) if argv else raw_input()).strip()
    if not expression:
        return None

    stack, postfix = [], []
    prev = None # dummy value, only for first iteration
    for token in tokens(expression):
        for pred, rule in ALLOWED_SEQUENCES.iteritems():
            if pred(prev) and not rule(token):
                raise DijkstraError('wrong syntax')

        if isop(token):
            if (isop(prev) or prev == None or prev == '(') and isunary(token):
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
            while stack and stack[-1] != '(':
                postfix.append(stack.pop())
            if not stack:
                raise DijkstraError('parentheses mismatch')
            stack.pop()

        else:
            postfix.append(float(token))

        prev = token

    if isop(prev):
        raise DijkstraError('wrong syntax')

    while stack:
        operator = stack.pop()
        if operator == '(':
            raise DijkstraError('parentheses mismatch')
        postfix.append(operator)

    result = calculate(map(lambda x: RESTORE_UNARY.get(x, x), postfix))
    print int(result) if result.is_integer() else result
    return 0


if __name__ == '__main__':
    try:
        sys.exit(main(sys.argv[1:]))
    except KeyboardInterrupt:
        print 'Interrupted by user'
    except DijkstraError, ex:
        print u'%s: error: %s' % (os.path.basename(__file__), ex.args[0].lower())
    sys.exit(1)
