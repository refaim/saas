#!/usr/bin/env python
import sys
import os
import time
import subprocess
import filecmp

args = sys.argv[1:]
aidx = args.index('-f')
archive, files = args[aidx+1], args[aidx+2:]

is_create = '-c' in args
if is_create:
    source_size = sum(os.path.getsize(file) for file in files)
else:
    source_size = os.path.getsize(archive)

start = time.time()
subprocess.call(args)
minutes = (time.time() - start) / 60
print '%.2f MB/min' % ((float(source_size) / (1024 ** 2)) / minutes)

if not is_create:
    target = archive + '.ex'
    match, mismatch, errors = filecmp.cmpfiles(os.getcwd(), target, os.listdir(target))
    if mismatch:
        print 'Mismatched files: %s' % ', '.join(mismatch)
    if errors:
        print 'Errors: %s' % ', '.join(errors)
