import http.client, urllib.parse,json
import threading
from random import randint
from time import sleep
def connect(q):
    params = urllib.parse.urlencode({'@number': 12524, '@type': 'issue', '@action': 'show'})
    headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
    conn = http.client.HTTPConnection("192.168.0.68",9000)
    for num in range(0,1):
        data={'sn':str(num)+str(q),'model':'model','version':'vesion','password':'password'}
        data1 = json.dumps(data)
        print('hello')
        conn.request("GET", "/api/v1/nodes", data1, headers)#the third argv is body
        response = conn.getresponse()
        print (response.status, response.reason)
        data = response.read()
        print(data)
        string1 = data.decode('utf-8')
        print(string1)
        print(str(string1))
        #print(type(string(data)))
        json_obj = json.loads(string1)
        


      
    conn.close()
connect(6)

