#!/usr/bin/env python

import sys
import unittest
import random
import logging
import string
import optparse

import dijkstra
from dijkstra import main as run_dijkstra

def roll():
    return random.random() > 0.5

def number():
    return random.random() * 100 * (-1 if roll() else 1)

def operator():
    return random.choice(dijkstra.OPERATORS)

def process(expr):
    return run_dijkstra(expr.split())

def get_binary_expression(op1=None, op2=None):
    op1 = op1 or operator()
    op2 = op2 or operator()
    return '(%f) %s (%f) %s (%f)' % (number(), op1, number(), op2, number())

def get_unary_expression():
    return '(%s%f)' % (random.choice(dijkstra.UNARY_OPERATORS.keys()), number())

def evaluate(expr):
    try:
        return eval(expr)
    except (OverflowError, ZeroDivisionError, ValueError):
        return None


class DijkstraTestCase(unittest.TestCase):

    def assertError(self, expr):
        log.debug(expr)
        self.assertRaises(dijkstra.DijkstraError, process, expr)


    def assertErrors(self, *args):
        map(self.assertError, args)


    def assertRoughlyEqual(self, expr, value):
        log.debug(expr)
        self.assertAlmostEqual(process(expr), value, 6, expr)


    def testSyntaxErrors(self):
        self.assertError(random.choice(string.letters))

        n1, n2, n3 = number(), number(), number()

        self.assertError('%f %f' % (n1, abs(n2)))
        self.assertError('%d.%d.%d' % tuple(map(int, (n1, n2, n3))))

        oppair = ((op1, op2) for op1 in dijkstra.OPERATORS for op2 in dijkstra.OPERATORS)
        for op1, op2 in oppair:
            self.assertError(op1)
            self.assertError(op1 + op2)
            self.assertError('%f%s' % (n1, op1))
            if op1 not in dijkstra.UNARY_OPERATORS:
                self.assertError('%s%f' % (op1, abs(n1)))

        patterns = ('(', ')', '((', '))', '()', ')(', '(()', ')()')
        for item in patterns:
            self.assertError(item)
            self.assertErrors('%f%s' % (n1, item), '%s%f' % (item, n1))
            for op in dijkstra.OPERATORS:
                self.assertErrors(op + item, item + op)


    def testSimple(self):
        n = number()
        self.assertRoughlyEqual('%f' % n, n)
        self.assertRoughlyEqual('(%f)' % n, n)
        self.assertRoughlyEqual('((%f))' % n, n)


    def testArithmeticExceptions(self):
        self.assertError('5 / 0')
        self.assertError('0 / 0')
        self.assertError('0 ** -1')
        self.assertError('10 ** 10 ** 10')
        self.assertError('-1.1 ** - 1.1')


    def testUnaryOperators(self):
        n = int(abs(number()))
        for op in dijkstra.UNARY_OPERATORS:
            for i in (n, n+1): # odd and even
                expression = '%s%f' % (op * i, number())
                self.assertRoughlyEqual(expression, eval(expression))


    def testBinaryOperators(self):
        oppair = ((op1, op2) for op1 in dijkstra.OPERATORS for op2 in dijkstra.OPERATORS)
        for op1, op2 in oppair:
            while True:
                expression = get_binary_expression(op1, op2)
                result = evaluate(expression)
                if result is None: continue
                break
            log.debug(expression)
            self.assertRoughlyEqual(expression, result)


    def testCalculator(self):
        for i in range(10):
            while True:
                expression = []

                for j in range(3):
                    if roll():
                        expression.append(get_binary_expression())
                    else:
                        expression.append(get_unary_expression())
                    if roll():
                        expression.insert(0, '(')
                        expression.append(')')
                    expression.append(operator())

                expression.pop()
                expression = ''.join(expression)
                result = evaluate(expression)
                if result is None: continue
                break
            self.assertRoughlyEqual(expression, result)


if __name__ == '__main__':
    try:
        oparser = optparse.OptionParser()
        oparser.add_option('-d', '--debug', action='store_true', help='show expressions')

        options, args = oparser.parse_args()

        logging.basicConfig(level=logging.DEBUG if options.debug else logging.INFO)
        log = logging.getLogger()

        tests = [
            'testSyntaxErrors',
            'testSimple',
            'testArithmeticExceptions',
            'testUnaryOperators',
            'testBinaryOperators',
            'testCalculator',
        ]

        suite = unittest.TestSuite(map(DijkstraTestCase, tests))
        result = unittest.TestResult()
        suite.run(result)

        if not result.wasSuccessful():
            errors = result.errors + result.failures
            for fail in errors:
                print '\n'.join(map(str, fail))
        else:
            print '%d tests passed' % result.testsRun
    except KeyboardInterrupt:
        print 'Interrupted by user'
