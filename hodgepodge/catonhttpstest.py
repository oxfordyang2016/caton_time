import urllib3
import timeit
from datetime import datetime
import json
encoded_body = json.dumps({
        "username": "yangming1",
        "password": "123456",
        
    })
start = timeit.default_timer()
ki=datetime.now()
l=[]
for k in range(1,2):
    try:
        c = urllib3.HTTPSConnectionPool('54.223.145.87', port=500, timeout=6,cert_reqs='CERT_NONE',assert_hostname=False)
        t=c.request('POST', '/login/',headers={'Content-Type': 'application/json'},body=encoded_body)
        #data=c.getresponse().read()
        print(t.status)
        print(t.data)
        l.append(1)
    except:
        pass