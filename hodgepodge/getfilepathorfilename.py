import os
def listfile(dir):
    filename=os.listdir('')
    for file in filename:
        print dir+file
listdir('/tmp/goodday/')        