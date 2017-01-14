import ftplib,os

ftp = ftplib.FTP("192.168.0.236")
ftp.login("test", "test")

data = []
print(ftp.pwd())
ftp.cwd("/ftp/test/")
'''
ftp.dir(data.append)

ftp.quit()

for line in data:
    print ("-", line)
'''
'''
def upload(ftp, file):
    ext = os.path.splitext(file)[1]
    if ext in (".txt", ".htm", ".html"):
        ftp.storlines("STOR " + file, open(file))
    else:
        ftp.storbinary("STOR " + file, open(file, "rb"), 1024)

upload(ftp, 'C:/Users/dell/Desktop/goproject/caton_time/hodgepodge/gro.go')

ftp.close()
'''
#

filename = "gro.go"


os.chdir(r"C:\Users\dell\Desktop\goproject\caton_time\hodgepodge")
myfile = open(filename, 'r')
ftp.storlines('STOR ' + filename, myfile)
myfile.close()

