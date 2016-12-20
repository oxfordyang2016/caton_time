#!/bin/bash
echo "<========================full backup=============================>"
sqlsuffix=".sql"
#a=`date +'%m%d%Y'`
a=`date +'%Y%m%d%H%M'`
sdatefile=$a$sqlsuffix
echo $sdatefile
pd=`pwd`

#full-backup(a time/week)
echo "/tmp/goodday/$sdate"
/usr/bin/mysqldump -u root -p503951 --single-transaction --flush-logs --master-data=2  --all-databases >"/tmp/goodday/purebackdata/$sdatefile"
#Note that the sys.argv[1] is full path
/usr/bin/python /tmp/goodday/fullbackup.py "$sdatefile"