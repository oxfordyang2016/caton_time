import http.client, urllib.parse,json
import threading
from random import randint
from time import sleep
def connect():
    params = urllib.parse.urlencode({'@number': 12524, '@type': 'issue', '@action': 'show'})
    headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
    conn = http.client.HTTPConnection("localhost",8087)
    for num in range(0,2):
        data={'sn':str(num),'model':'model','version':'vesion','password':'password'}
        data1 = json.dumps(data)
        print('hello')
        conn.request("POST", "/api/v1/login", data1, headers)#the third argv is body
        response = conn.getresponse()
        print (response.status, response.reason)
        data = response.read()
        print(data)
        string1 = data.decode('utf-8')
        print(string1)
        print(str(string1))
        #print(type(string(data)))
        json_obj = json.loads(string1)
        print(json_obj['Token'])


        if json_obj['Token']!="error":
            conn1 = http.client.HTTPConnection("localhost",8087)
            #Authorization: Bearer <token>
            #data1 = json.dumps(data)
            au="Bearer"+" "+json_obj['Token']
            headers = {"content-type": "application/x-www-form-urlencoded", "authorization":au}
            #sleep(4)
            conn1.request("POST", "/api/v1/report", data1, headers)#the third argv is body
            response1 = conn1.getresponse()
            print (response1.status, response1.reason)
            data2 = response1.read()
            print(data2)    
            print(num)
            conn1.close()
    conn.close()
connect()
    
