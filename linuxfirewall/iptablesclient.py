import http.client, urllib.parse,json,time
params = urllib.parse.urlencode({'@number': 12524, '@type': 'issue', '@action': 'show'})
headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
conn = http.client.HTTPConnection("192.168.0.81",8000)
for num in range(0,2):
    data={'sn':'hallo','model':'model','version':'vesion','password':'password'}
    datat="helloyangming"+"this is "+str(num)+"   times"
    data1 = json.dumps(datat)
    #print('hello')
    conn.request("POST", "/", data1, headers)#the third argv is body
    response = conn.getresponse()
    print (response.status, response.reason)
    data = response.read()
    print(data)
    string1 = data.decode('utf-8')
    print(string1)
    print(str(string1))
    time.sleep(10)
    #print(type(string(data)))
    #json_obj = json.loads(string1)
    #print(json_obj['Token'])

'''
    if json_obj['Token']!="error":
        conn1 = http.client.HTTPConnection("localhost",8084)
        #Authorization: Bearer <token>
        #data1 = json.dumps(data)
        au="Bearer"+" "+json_obj['Token']
        headers = {"content-type": "application/x-www-form-urlencoded", "authorization":au}
        conn1.request("POST", "/api/v1/report", data1, headers)#the third argv is body
        response1 = conn1.getresponse()
        print (response1.status, response1.reason)
        data2 = response1.read()
        print(data2)
        print(num)
        conn1.close()
'''

conn.close()