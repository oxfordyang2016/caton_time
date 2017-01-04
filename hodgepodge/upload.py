import ftplib
import sys
import re
import os
from ctypes import *
sqlfilelist=[]
binfilelist=[]
barebinfilelist=[]

#-----------via c lib to get  machinecode and generate usr pass--------------
def getuserinfo(macget_lib_position):
    authlib = CDLL(macget_lib_position)
    mac=create_string_buffer(288)
    authlib.authGetFingerprint(pointer(mac),sizeof(mac))
    mac=mac.value#this is macode
    username='u'+mac[0:6]
    password='p'+mac[-6:]
    return [username,password]

#-----------------via list to get bin and sql filename dir-----
def listfile(dirname):
    filenamedir=os.listdir(dirname)
    for filename in filenamedir:
        if "sql" in filename and os.path.isfile(dirname+filename):

            sqlfilelist.append(dirname+filename)
        if "bin" in filename and os.path.isfile(dirname+filename):
            binfilelist.append(dirname+filename)
            barebinfilelist.append(filename)
    return [sqlfilelist,binfilelist]


userinfo=getuserinfo('/opt/cydex/libauth.so')
server='192.168.1.211'
ftp_connection = ftplib.FTP(server, userinfo[0], userinfo[1])
def get_back_file_path():
#get sort dir including bin  and sql
    bin_sql_path_to_be_sorted_dir=listfile('/var/lib/mysql/')[1]+listfile('/opt/cydex/')[0]
    file_create_time_list=map(lambda x:os.path.getctime(x),bin_sql_path_to_be_sorted_dir)
    time_file_dict = dict(zip(file_create_time_list,bin_sql_path_to_be_sorted_dir))
    backup_dir_path=[time_file_dict[key] for key in sorted(time_file_dict.iterkeys())]
    print backup_dir_path
    return backup_dir_path
def upload(upload_path):
    fh = open(upload_path,'rb')
    ftp_connection.storbinary('STOR %s'%os.path.basename(upload_path), fh)
    os.remove(upload_path)
    fh.close()
    return upload_path+'====upload ok'
if len(backup_dir_path)>5:
    map(os.remove,get_back_file_path())
    os.system('/usr/bin/mysqldump -u root -p503951 --single-transaction --flush-logs --master-data=2  --database cydex cydex_user >"/opt/cydex/$sdatefile"')
map(upload,get_back_file_path()[:-4])











