'''
from http.client import HTTPSConnection
import ssl
ctx = ssl.create_default_context()
ctx.check_hostname = False
ctx.verify_mode = ssl.CERT_NONE
 requests.get('https://kennethreitz.com', verify=False)
conn = HTTPSConnection('54.223.102.182', 6002)
headers = { 'Authorization' : 'Basic'}
#then connect
conn.request('GET', '/?query="alias"', headers=headers)
#get the response back
res = conn.getresponse()
# at this point you could check the status etc
# this gets the page text
data = res.read()  
'''
'''
import http.client
import socket
print(socket.ssl)
conn = http.client.HTTPSConnection("54.223.102.182",6002)
conn.request("GET", "/")
r1 = conn.getresponse()
print(r1.status, r1.reason)
'''
'''
import urllib

verify='/etc/ssl/certs/cacert.org.pem
'response =urllib.request.get('https://54.223.102.182:6002', verify=False)
'''
'''
from urllib import request
request.get('https://kennethreitz.com', verify=False)
'''
import urllib3
import timeit
from datetime import datetime
start = timeit.default_timer()
ki=datetime.now()
l=[]
while 1:
    try:
        c = urllib3.HTTPSConnectionPool('54.223.102.182', port=6002, timeout=6,cert_reqs='CERT_NONE',assert_hostname=False)
        t=c.request('GET', '/?query=alias')
        #data=c.getresponse().read()
        print(t.status)
        print(t.data)
        l.append(1)
    except :
        l.append(0)    
        print("error occur")
    finally:
        print(l)
        stop = timeit.default_timer()
        k=stop-start


        if k > 1800:
            
            break


print("start time "+str(ki))
print("send data  " +str(len(l)) +"   times")
print("except "+str(l.count(0))+" times")
print("end time "+str(datetime.now()))






