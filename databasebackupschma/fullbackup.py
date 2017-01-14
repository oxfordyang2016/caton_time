import ftplib
import sys
import os
#print sys.argv[1]##python ftpupload.py -u     will print -u
from ConfigParser import ConfigParser
from ctypes import *

sqlfilelist=[]
binfilelist=[]
barebinfilelist=[]
#----------------------get dir file-------------------
def listfile(filenamedir):
    for filename in filenamedir:
        if "sql" in filename:
            sqlfilelist.append(dirname+filename)
        if "bin" in filename:
            binfilelist.append(dirname+filename)
            barebinfilelist.append(filename)
    return [sqlfilelist,binfilelist]

#your path format is below
#getdirfilelist=listfile('/var/lib/mysql/')

#----------------via maccode to get username and pass----------
def userinfo(macget_lib_position):
    authlib = CDLL(macget_lib_position)
    mac=create_string_buffer(288)
    authlib.authGetFingerprint(pointer(mac),sizeof(mac))
    mac=mac.value#this is macode
    username='u'+mac[0:6]
    password='p'+mac[-6:]
    return [username,password]


#-------------------via ini to get upload info -------------------

def userinfo1(ini_file):
    config = ConfigParser()
    config.read(ini_file)
    #coperate vsftpd!!!!
    server = config.get('server', 'address')
    username = config.get('server', 'user')
    password = config.get('server', 'password')
    remote_path = config.get('server', 'remotelocation')
    return [server,username,password,remote_path]


def upload2remote(ip,username,password):
    ftp_connection = ftplib.FTP(ip,username, password)

fh = open('/tmp/goodday/purebackdata/'+sys.argv[1],'rb')
ftp_connection.storbinary('STOR %s'%sys.argv[1], fh)
fh.close


#ftp_connection.cwd(remote_path)
#-------------------copy files---------
#uploadbinfile_list=listfile('/var/lib/mysql/')[1]
#uploadbinfile_path=(sorted(uploadbinfile_list))[-3]

#print sorted(barebinfilelist)
#fh_bin = open(uploadbinfile_path,'rb')
#-----it is designed to upload bin file when full backup-------

'''

server = config.get('server', 'address')
username = config.get('server', 'user')
password = config.get('server','password')
ftp_connection = ftplib.FTP(server, username, password)
remote_path =config.get('server','remotelocation')
ftp_connection.cwd(remote_path)
uploadbinfile_list=listfile('/var/lib/mysql/')[1]
uploadbinfile_path=(sorted(uploadbinfile_list))[-3]

print sorted(barebinfilelist)
fh = open(uploadbinfile_path,'rb')
#--------------can torant error---------------
#centerpos=remote_path+currentpos
#if currentpos!=sorted(barebinfilelist)[-3]:


server = config.get('server', 'address')
username = config.get('server', 'user')
password = config.get('server','password')
ftp_connection = ftplib.FTP(server, username, password)
remote_path =config.get('server','remotelocation')
ftp_connection.cwd(remote_path)
uploadbinfile_list=listfile('/var/lib/mysql/')[1]
uploadbinfile_path=(sorted(uploadbinfile_list))[-3]

print sorted(barebinfilelist)
fh = open(uploadbinfile_path,'rb')
#--------------can torant error---------------
#centerpos=remote_path+currentpos
#if currentpos!=sorted(barebinfilelist)[-3]:




ftp_connection.storbinary('STOR /tmp/goodday/purebackdata/%s'%(sorted(barebinfilelist))[-3], fh)
fh.close()

'''



#------------------------------origin------------
'''
#coperate vsftpd!!!!
server = '192.168.201.139'
username = 'sk'
password = 'centos'
ftp_connection = ftplib.FTP(server, username, password)
remote_path = "/tmp/goodday/"

ftp_connection.cwd(remote_path)
string='/tmp/goodday/%s'%sys.argv[1]
passfilepos=sys.argv[1]
print string
fh = open(sys.argv[1],'rb')
ftp_connection.storbinary('STOR %s'%sys.argv[1], fh)
fh.close()
'''


#--------------Check if the file meets the upload criteria-------------

def filecheck():
    print('''
i need to upload pypi to show,if you see this,it is ok !!!!''')