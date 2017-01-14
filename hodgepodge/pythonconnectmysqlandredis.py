#!/usr/bin/python
import MySQLdb
import redis
# connect
db = MySQLdb.connect(host="192.168.0.68", user="root", passwd="503951",db="cydex")

cursor = db.cursor()

# execute SQL select statement
cursor.execute("SELECT * FROM package_file")

# commit your changes
db.commit()

# get the number of rows in the resultset
numrows = int(cursor.rowcount)

# get and display one row at a time.
for x in range(0,numrows):
    row = cursor.fetchone()
    print( row[0], "-->", row[1])

#--------------------------------------------------------------

POOL = redis.ConnectionPool(host='192.168.0.68', port=6379, db=0)

def getVariable(variable_name):
    my_server = redis.Redis(connection_pool=POOL)
    response = my_server.get(variable_name)
    return response

def setVariable(variable_name, variable_value):
    my_server = redis.Redis(connection_pool=POOL)
    my_server.set(variable_name, variable_value)
print(getVariable("hu"))


r = redis.StrictRedis(host='192.168.0.68', port=6379, db=0)
print(r.keys())

























