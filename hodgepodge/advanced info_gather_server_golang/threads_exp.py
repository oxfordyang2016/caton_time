import http.client, urllib.parse,json
import threading
from random import randint
from time import sleep
def connect(q):
    params = urllib.parse.urlencode({'@number': 12524, '@type': 'issue', '@action': 'show'})
    headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
    conn = http.client.HTTPConnection("localhost",8084)
    for num in range(0,30):
        data={'sn':str(num)+str(q),'model':'model','version':'vesion','password':'password'}
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
    conn.close()


def print_number(number):
    # Sleeps a random 1 to 10 seconds
    rand_int_var = randint(1, 1000)
    sleep(rand_int_var)
    print ("Thread " + str(number) + " slept for " + str(rand_int_var) + " seconds"+"---------------------------------------------->")
    connect(number)

thread_list = []

for i in range(1, 100):
    # Instantiates the thread
    # (i) does not make a sequence, so (i,)
    t = threading.Thread(target=print_number, args=(i,))
    # Sticks the thread in a list so that it remains accessible
    thread_list.append(t)

# Starts threads
for thread in thread_list:
    thread.start()

# This blocks the calling thread until the thread whose join() method is called is terminated.
# From http://docs.python.org/2/library/threading.html#thread-objects
for thread in thread_list:
    thread.join()

# Demonstrates that the main process waited for threads to complete
print ("Done")
