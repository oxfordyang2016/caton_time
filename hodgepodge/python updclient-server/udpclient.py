import time
from socket import *
clientSocket = socket(AF_INET, SOCK_DGRAM)
clientSocket.settimeout(6)
addr = ("54.223.102.182", 19345)
addr1 = ("127.0.0.1", 12000)

for pings in range(100):        
    message = bytes('yangming is here'+str(pings),"utf-8")    
    start = time.time()
    clientSocket.sendto(message, addr)
    clientSocket.sendto(message, addr1)
    time.sleep(2)
'''
    try:
        data, server = clientSocket.recvfrom(1024)
        end = time.time()
        elapsed = end - start
        print ('%s %d %d' % (data, pings, elapsed))
    except timeout:
        print ('REQUEST TIMED OUT')

    '''  
#when you do a task ,you can cancel to comment above