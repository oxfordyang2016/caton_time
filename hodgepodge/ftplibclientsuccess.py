import ftplib
#coperate vsftpd!!!!
server = '123.206.232.46'
username = 'sk'
password = 'centos'
ftp_connection = ftplib.FTP(server, username, password)
remote_path = "/tmp/"
ftp_connection.cwd(remote_path)
fh = open("/root/pyftpclient.py", 'rb')
ftp_connection.storbinary('STOR awesomeday.py', fh)
fh.close()