import random 
from socket import *
serverSocket = socket(AF_INET, SOCK_DGRAM)
serverSocket.bind(('', 12000))#it is listenning any ip into my machine in port 12000

while True:
    rand = random.randint(0, 10)
    message, address = serverSocket.recvfrom(1024)
    message = message.upper()
    print(message.strip(),address)
    if rand >= 4:
        serverSocket.sendto(message, address)#senf info to client 
         


'''

import socket

# A UDP server

# Set up a UDP server
UDPSock = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)

# Listen on port 21567
# (to all IP addresses on this system)
listen_addr = ("",21567)
UDPSock.bind(listen_addr)

# Report on all data packets received and
# where they came from in each case (as this is
# UDP, each may be from a different source and it's
# up to the server to sort this out!)
while True:
        data,addr = UDPSock.recvfrom(1024)
        print data.strip(),addr

'''        