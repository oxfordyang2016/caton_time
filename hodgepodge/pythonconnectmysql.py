#!/usr/bin/python
import MySQLdb
#if all set is ok ,in mariadb ,please firewalld ,mariadb status,reboot method is good
import redis
import sys
# connect
print sys.argv[1:]
db = MySQLdb.connect(host="192.168.0.237", user="root", passwd="503951",db="cydex")
cursor = db.cursor()

#define a function to get table row info and write it to dict
def getrow():
    # commit your changes
    db.commit()
    tabledict={}
    numrows = int(cursor.rowcount)
    num_fields = len(cursor.description)
    field_names = [i[0] for i in cursor.description]
    for x in range(0,numrows):
        row = cursor.fetchone()
        #print(row)
        tmpdict={}
        for k in range(0,len(row)):
            #print str(field_names[k])+"                 |---------------------------->"+str(row[k]) 
            tmpdict[str(field_names[k])]=str(row[k])
        tabledict[str(x)]=tmpdict
    return tabledict

#--------set command line argvs-----------
if sys.argv[1]=="-sql":
	# execute SQL select statement
    #cursor.execute("SELECT * FROM package_file limit 2 offset 2")
    cursor.execute(sys.argv[2])
    numrows = int(cursor.rowcount)
    num_fields = len(cursor.description)
    field_names = [i[0] for i in cursor.description]
    for x in range(0,numrows):
        row = cursor.fetchone()
        print row
    sys.exit()
elif sys.argv[1]=="--task":
	cursor.execute("SELECT * FROM "+"transfer_task"+"  WHERE "+str(sys.argv[2])+" = "+" '"+str(sys.argv[3])+"'")
    


task=getrow()
cursor.execute("select * from package_job_detail where job_id= "+"'"+str(task["0"]["job_id"])+"'"+" AND "+"fid "+"= "+"'"+str(task["0"]["fid"])+"'")
job_detail=getrow()

cursor.execute("select * from package_job  where job_id= "+"'"+str(task["0"]["job_id"])+"'")
job=getrow()
cursor.execute("select * from package_file  where fid= "+"'"+str(task["0"]["fid"])+"'")
file=getrow()
cursor.execute("select * from stat_transfer  where node_id= "+"'"+str(task["0"]["node_id"])+"'")
n=getrow()["0"]#node stat_transfer
cursor.execute("select * from transfer_node  where nid= "+"'"+str(task["0"]["node_id"])+"'")
node=getrow()["0"]#node stat_transfer





#print task,job,job_detail 
t=task["0"]
d=job_detail["0"]
j=job["0"]
f=file["0"]
cursor.execute("select * from package_pkg  where pid= "+"'"+str(f["pid"])+"'")
p=getrow()["0"]#package 
print"""
===================%s=============================================
task_id                 ----------------------------%s                           
task_create_at          ----------------------------%s
================corresponding package info============================= 
co_package title        -----------------------------%s
corresponding pid       ----------------------------%s
co_package create time  ----------------------------%s
co_package update time  ----------------------------%s
co_package notes        -----------------------------%s
co_package file nums    -----------------------------%s
co_package totalsize    -----------------------------%s
===================transfer stat info======================================
corresponding uid       ----------------------------%s                     
type                    ----------------------------%s   [from job detail] 
node_id                 ----------------------------%s
node_name               ----------------------------%s
zone_id                 ----------------------------%s
rx_total_bytes          ----------------------------%s
rx_max_bitrate          ----------------------------%s
rx_min_bitrate          ----------------------------%s
rx_avg_bitrate          ----------------------------%s
rx_tasks                ----------------------------%s
tx_total_bytes          ----------------------------%s
 tx_max_bitrate         ----------------------------%s
tx_min_bitrate          ----------------------------%s
tx_avg_bitrate          ----------------------------%s
tx_tasks                ----------------------------%s
stat_tran_create_at     ----------------------------%s
stat_tran_update_at     ----------------------------%s
==============corresponding job detail========================================
job_id                  ----------------------------%s
type                    ----------------------------%s   [from job]
corresponding fid       ----------------------------%s
state                   ----------------------------%s   [from job detail] 
checked                 ----------------------------%s                     
num_segs of the file    ----------------------------%s
num_finished_segs       ----------------------------%s
correspond filesize     ----------------------------%s   
finished_size           ----------------------------%s   
cof_start_time          ----------------------------%s
cof_finish_time         ----------------------------%s
d_create_at             ----------------------------%s
d_update_at             ----------------------------%s  
j_create_at             ----------------------------%s
J_update_at             ----------------------------%s
J_soft_del              ----------------------------%s    
================corresponding file  info==========================================
correspond filename     ----------------------------%s 
corresponding fid       ----------------------------%s   
correspond filesize     ----------------------------%s   
correspond file path    ----------------------------%s   
correspond file abs path ----------------------------%s   
===================node static info=================================================
correspond node_id      ----------------------------%s 
corresponding zone_id   ----------------------------%s   
correspond public_ip    ----------------------------%s   [from package_file]
correspond rx_bandwith  ----------------------------%s   
correspond tx_bandwith  ----------------------------%s   
correspond last_login   ----------------------------%s   
correspond last_logout  ----------------------------%s  
correspond create_time  ----------------------------%s   
correspond updatetime   ----------------------------%s 
========================END==================================================

"""%(t['task_id'],t['task_id'],t['create_at'],p['title'],j['pid'],
    p['create_at'],p['update_at'],p['notes'],p['num_files'],
    p['size'],j['uid'],t['type'],t['node_id'],node['name'],t['zone_id'],
    n['rx_total_bytes'],n['rx_max_bitrate'],n['rx_min_bitrate'],n['rx_avg_bitrate'],
   n['rx_tasks'],n['tx_total_bytes'],n['tx_max_bitrate'],
   n['tx_min_bitrate'],n['tx_avg_bitrate'],n['tx_tasks'],n['create_at'],n['update_at'],
    t['job_id'],j['type'],t['fid'],d['state'],
   d['checked'], t['num_segs'],d['num_finished_segs'],
   f['size'],d['finished_size'],d['start_time'],d['finish_time']
   ,d['create_at'],d['update_at'],j['create_at'],j['update_at'],j['soft_del']
   ,f['name'],t['fid'],f['size'],f['path'],f['path_abs']
    ,node['nid'],node['zone_id'],node['public_addr'],node['rx_bandwidth'],node['tx_bandwidth'],
    node['last_login_time'],node['last_logout_time'],node['create_at'],node['update_at']
  )


'''
            id: 1
    machine_code: 141c00767ee033454bc3218716ef46ec
             nid: 4038ecf3-fb49-42f6-9fa3-ab8d978ca43f
            name: 
         zone_id: 
     public_addr: 
    rx_bandwidth: 20971520
    tx_bandwidth: 20971520
 last_login_time: 2016-11-16 00:56:55
last_logout_time: 2016-11-14 11:09:10
       create_at: 2016-10-27 18:57:38
       update_at: 2016-11-16 00:56:55


'''



'''
            id: 1
       node_id: 4038ecf3-fb49-42f6-9fa3-ab8d978ca43f
rx_total_bytes: 2193314214
rx_max_bitrate: 91268352
rx_min_bitrate: 14820
rx_avg_bitrate: 469510
      rx_tasks: 285
tx_total_bytes: 17458522528
tx_max_bitrate: 32722560
tx_min_bitrate: 256
tx_avg_bitrate: 9854109
      tx_tasks: 270
     create_at: 2016-10-27 19:08:42
     update_at: 2016-11-07 10:30:42
1 row in set (0.00 sec)
'''
































'''
# get the number of rows in the resultset
numrows = int(cursor.rowcount)

#get mysql column name
num_fields = len(cursor.description)
field_names = [i[0] for i in cursor.description]
print("this table have "+str(num_fields)+" columns---->"+" their names are --->"+str(field_names))
# get and display one row at a time.
s="""              
                   ||
        -------power by yangming----
                   ||
"""
taskdic={}
for x in range(0,numrows):
    row = cursor.fetchone()
    print s
    for k in range(0,len(row)):
    	print str(field_names[k])+"                 |---------------------------->"+str(row[k])	
        #create a dictionary
        taskdic[str(field_names[k])]=str(row[k])
    print(row)
    print(taskdic)

print taskdic["job_id"]
print taskdic["fid"]
cursor.execute("select * from package_job_detail where job_id= "+"'"+str(taskdic["job_id"])+"'"+" AND "+"fid "+"= "+"'"+str(taskdic["fid"])+"'")
db.commit()







job=getrow()
print "this is corresponding package_job_detail--------------->",job
'''
















