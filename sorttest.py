#!/usr/bin/env python
import time
import timeit
times1=[
    "2010-04-20 10:07:30",
    "2010-04-20 10:07:38",
    "2010-04-20 10:07:52",
    "2010-04-20 10:08:22",
    "2010-04-20 10:08:22",
    "2010-04-20 10:09:46",
    "2010-04-20 10:10:37",
    "2010-04-20 10:10:58",
    "2010-04-20 10:11:50",
    "2010-04-20 10:12:13",
    "2010-04-20 10:12:13",
    "2010-04-20 10:25:38"
]

times2=[
    "2010-04-20 10:07:30",
    "2010-04-20 10:07:38",
    "2010-04-20 10:07:52",
    "2010-04-20 10:08:22",
    "2010-04-20 10:08:22",
    "2010-04-20 10:09:46",
    "2010-04-20 10:10:37",
    "2010-04-20 10:10:58",
    "2010-04-20 10:11:50",
    "2010-04-20 10:12:13",
    "2010-04-20 10:12:13",
    "2010-04-20 10:25:38"
]


def sorttime1():
    times1.sort(key=lambda x: time.strptime(x, '%Y-%m-%d %H:%M:%S')[0:6], reverse=True)



def sorttime2():
    times2.sort(key=lambda x: time.strptime(x, '%Y-%m-%d %H:%M:%S')[0:6])
    times2.reverse()


t1 = timeit.timeit('sorttime1()', setup="from __main__ import sorttime1; gc.enable()", number=100000)
t2 = timeit.timeit('sorttime2()', setup="from __main__ import sorttime2; gc.enable()", number=100000)

print "sort and reverse", t1
print "sort then reverse", t2

