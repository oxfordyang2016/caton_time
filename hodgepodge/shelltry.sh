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
/usr/bin/mysqldump -u root -pyangmingtestmysql --single-transaction --flush-logs --master-data=2  --all-databases >"/tmp/goodday/$sdatefile"
/usr/bin/python /tmp/goodday/ftpupload.py "/tmp/goodday/$sdatefile"