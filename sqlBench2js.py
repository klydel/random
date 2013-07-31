#!/usr/bin/env python
#converts mysql bench mark results into js includes for highcharts
#yes this could be coded better
import os
import sys
categories = []
newData = []
oldData = []

if len(sys.argv) < 2:
    sys.exit('Usage: %s OLDSTATS NEWSTATS' % sys.argv[0])

if not os.path.exists(sys.argv[1]) and not os.path.exists(sys.argv[2]):
    sys.exit('ERROR: File was not found!')

oldStats = sys.argv[1]
newStats = sys.argv[2]

funcName = sys.argv[1].split('-')
funcName = funcName[0].split('_')
funcName = funcName[1]

if funcName == "big":
    new = os.popen("cat %s | grep 'Time' " % (newStats))
    old = os.popen("cat %s | grep 'Time' " % (oldStats))
else:
    new = os.popen("cat %s | grep 'Time' | awk '{print $3, $5}'" % (newStats))
    old = os.popen("cat %s | grep 'Time' | awk '{print $3, $5}'" % (oldStats))

for k in new:
    c, d = k.split()
    categories.append(c)
    newData.append(float(d))
for i in old:
    a, b = i.split()
    oldData.append(float(b))

fStats = open(funcName+'.js', 'w')
fTotal = open(funcName+'totals.js', 'w')

fStats.write(""" function get%sCategories(){
    return %s
        } function get%sNewData(){                                                                                                                                                
    return %s                                                                                                                                                                     
        } function get%sOldData(){                                                                                                                                                
    return %s                                                                                                                                                                     
        }   """ % (funcName, categories, funcName, newData, funcName, oldData))

fStats.close()

fTotal.write(""" function get%sStatsCategories(){
    return %s
        } function get%sStatsNewData(){                                                                                                                                                
    return %s                                                                                                                                                                     
        } function get%sStatsOldData(){                                                                                                                                                
    return %s                                                                                                                                                                     
        }   """ % (funcName, categories, funcName, newData, funcName, oldData))
fTotal.close()
