import MySQLdb,sys,os,time
inputfilepath,user=sys.argv[1],sys.argv[2]#from command line get
fullfile=os.path.basename(inputfilepath)
filename=os.path.splitext(fullfile)[0]
#connect databse and prepare to write data to database
db = MySQLdb.connect("localhost","root","123456","" )
cursor = db.cursor()
userdatabase=str(user)+'words'
createdbsql="create database if not exists "+userdatabase
cursor.execute(createdbsql)
db=MySQLdb.connect("localhost","root","123456",userdatabase)
cursor=db.cursor()
tablename=str(filename)
#create the table correspond to the file
createtablesql='CREATE TABLE '+'if not exists '+str(userdatabase)+'.'+str(tablename)+'( word varchar(30),  amountofunderstand int(5), amountofnotunderstand int(5) ,createtime varchar(20))'
cursor.execute(createtablesql)


#write words from a file to databse
#read lines from a file and write it to db.
with open(str(inputfilepath)) as f:
    for line in f:
        now=time.time()
        print(type(line))
        #insert word to table
        insertsql='insert into '+str(userdatabase)+'.'+str(tablename)+ ' (word,amountofunderstand,amountofnotunderstand,createtime)'+ 'values('+"'"+line+"'"+','+'0'+','+'0'+','+str(now)+')'
        print insertsql
        cursor.execute(insertsql)


db.commit()
db.close()
