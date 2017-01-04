#!/bin/bash

touch /opt/cydex/backlog
chmod 777 /opt/cydex/backlog
chmod 777 /opt/cydex/ftpupload.py
chmod 777 /opt/cydex/ftpuploadincre.py
chmod 777 /opt/cydex/fullbackup.sh
chmod 777 /opt/cydex/increbackup.sh
chmod 777 /opt/cydex/libauth.so

#this is used  to automatically add crontab item to crontab
(crontab -l ; echo "1 1 1 * * /bin/sh /opt/cydex/fullbackup.sh >> '/opt/cydex/backlog'  2>&1")| crontab -
(crontab -l ; echo "2 2 * * * /bin/sh /opt/cydex/increbackup.sh >> '/opt/cydex/backlog'  2>&1")| crontab -
(crontab -l ; echo "20 2 * * * /usr/bin/python /opt/cydex/upload.py >> '/opt/cydex/backlog'  2>&1")| crontab -