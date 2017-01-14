import os
from datetime import datetime
binfilelist=[]
def listfile(dirname):
    filenamedir=os.listdir(dirname)
    for filename in filenamedir:
        if "sql" in filename:
            sqlfilelist.append(dirname+filename)
        if "bin" in filename:
            binfilelist.append(dirname+filename)
    return [sqlfilelist,binfilelist]
#print getdirfilelist[1]

#----------------------load database increlog--------------------
    #path is backup dir,in these case ,it is /tmp/goodday/purebackdata
    #there error-deal capacity is bad



#-----------------------full load database---------------------------
#mysql -h 192.168.1.211 -u root -p  --port 15076 <201612161425.sql
# mysql -u root -pyangmingtestmysql</testdata/testincre.sql
#-----------------------full load database---------------------------
#mysql -h 192.168.1.211 -u root -p  --port 15076 <201612161425.sql
# mysql -u root -pyangmingtestmysql</testdata/testincre.sql
'''
def fullfileload(path,host,port,user,passwd):
    #path is backup dir,in these case ,it is /tmp/goodday/purebackdata
    #there error-deal capacity is bad
    sqldir=(listfile(path))[0]
    loadfilepath=(sorted(sqldir))[-1]
    print 'i am loading %s'%loadfilepath
    os.system(commandline)
#print getdirfilelist[1]
'''
#----------------------load database increlog--------------------
#mysqlbinlog /testdata/yangming-bin.000716| mysql -u root -pyangmingtestmysql

def increfileload(path,host,port,user,passwd):
    #path is backup dir,in these case ,it is /tmp/goodday/purebackdata
    bindir=(listfile(path))[1]
    loadfilepath=(sorted(bindir))[-1]
    print 'i am loading %s'%loadfilepath
    commandline='mysqlbinlog '+loadfilepath+'| mysql -u '+str(user)+' -p'+str(passwd)
    os.system(commandline)
#mysql -h 192.168.1.211 -u root -p  --port 15076 <201612161425.sql
# mysql -u root -pyangmingtestmysql</testdata/testincre.sql
def fullfileload(path,host,port,user,passwd):
    sqldir=(listfile(path))[0]
    loadfilepath=(sorted(sqldir))[-1]
    os.system(commandline)

#------------------------load function local(you can work with crontab)-----------------
def run():
    if sys.argv[1]=="f":
        fullfileload('/tmp/goodday/purebackdata/',0,0)
    elif sys.argv[1]=="i":
        increfileload('/tmp/goodday/purebackdata/',0,0)
def mysql_load(dir,host,port,user,passwd):
    def mv_file(backdir,copyfile):
        mvcommandline='mv '+copyfile+' '+dir+backdir
        if os.path.exists(dir+backdir):
            os.system(mvcommandline)
        else:
            os.mkdir(dir+backdir)
            os.system(mvcommandline)
    def pure_file_list(x):
        if (os.path.isfile(x) and ('bin' or 'sql' in x)):
            retun(dir+x)
        def con(filenamelist):
            return dir+filenamelist
        file_path_list=pure_file_list
        print file_path_list
        def get_file_create_time(file):
            return os.path.getctime(file)
        file_create_time_list=map(get_file_create_time,file_path_list)
        print "---------file path list-----"+str(file_path_list)
        print "---------min time-------"+str(file_create_time_list)
        sorted_time_filetime=min(file_create_time_list)
        all_file_info={}
        for key, value in all_file_info.iteritems():
            if value==min(file_create_time_list):
               if "bin" in key:
                   os.system(commandline)
                   mv_file('backdir',key)
               elif "sql" in key:
                   os.system(commandline)
                   mv_file('backdir',key)
               break#break can terminate the loop


def mysql_load2(dir,host,port,user,passwd):
    def mv_file(backdir,copyfile):
        mvcommandline='mv '+copyfile+' '+dir+backdir
        if os.path.exists(dir+backdir):
            os.system(mvcommandline)
        else:
            os.mkdir(dir+backdir)
            os.system(mvcommandline)
    def con(filenamelist):
        return dir+filenamelist
    def pure_file_list_get(x):
        if os.path.isfile(x) and ('bin' in x or 'sql' in x):
            return(x)
    while len(filter(pure_file_list_get,map(con,os.listdir(dir))))!=0:

        pure_file_list=filter(pure_file_list_get,map(con,os.listdir(dir)))
        #print 'i am running '+str(i)+'times'
        #dir_file_name=os.listdir(pure_file_list)

        #file_path_list=map(con,pure_file_list)
        #print file_path_list
        file_path_list=pure_file_list
        print file_path_list
        def get_file_create_time(file):
            return os.path.getctime(file)
        file_create_time_list=map(get_file_create_time,file_path_list)
        print "---------file path list-----"+str(file_path_list)
        print "---------min time-------"+str(file_create_time_list)
        sorted_time_filetime=min(file_create_time_list)
        all_file_info={}
        for file_path in file_path_list:
            all_file_info[file_path]=os.path.getctime(file_path)
        for key, value in all_file_info.iteritems():
            if value==min(file_create_time_list):
               if "bin" in key:
                   commandline='mysqlbinlog '+key+'| mysql -h '+str(host)+' --port '+str(port)+' -u '+str(user)+' -p'+str(passwd)
                   os.system(commandline)
                   mv_file('backdir',key)
               elif "sql" in key:
                   commandline='mysql -h '+str(host)+' --port '+str(port)+ ' -u '+str(user)+' -p'+ str(passwd)+ ' < '+key
                   os.system(commandline)
                   mv_file('backdir',key)
               break#break can terminate the loop