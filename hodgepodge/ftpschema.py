
A.set ftp server
1.install vsftpd and configure in linux https://www.liquidweb.com/kb/how-to-install-and-configure-vsftpd-on-centos-7/
2.when you change chroot_local_user=NO ,user will visit dir outside of /home/$user
3.useradd yangming     passwd yangming------>input password
4.launch service
5.allow user to visit dir outside of home
  http://www.ducea.com/2006/07/27/allowing-ftp-access-to-files-outside-the-home-directory-chroot/

B.use python pyftplibserver client
1.https://github.com/oxfordyang2016/caton_time/blob/master/hodgepodge/ftplibclientsuccess.py

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


2.chmod 777  /tmp/ in ftp server