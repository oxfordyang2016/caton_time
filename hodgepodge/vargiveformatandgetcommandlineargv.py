
import ftplib
import sys
#print sys.argv[1]##python ftpupload.py -u     will print -u

#coperate vsftpd!!!!
server = '123.206.232.46'
username = 'sk'
password = 'centos'
ftp_connection = ftplib.FTP(server, username, password)
remote_path = "/tmp/goodday/"
ftp_connection.cwd(remote_path)
string='"/tmp/goodday/%s"'%sys.argv[1]
print string
fh = open(string,'rb')
ftp_connection.storbinary('STOR yangmingawesomeday.py', fh)
fh.close()
~             