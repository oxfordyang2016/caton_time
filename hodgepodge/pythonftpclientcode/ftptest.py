import ftplib,os

ftp = ftplib.FTP("192.168.0.236")
ftp.login("test", "test")

data = []

ftp.dir(data.append)

ftp.quit()

for line in data:
    print ("-", line)

def upload(ftp, file):
    ext = os.path.splitext(file)[1]
    if ext in (".txt", ".htm", ".html"):
        ftp.storlines("STOR " + file, open(file))
    else:
        ftp.storbinary("STOR " + file, open(file, "rb"), 1024)

upload(ftp, r'C:\Users\dell\Desktop\goproject\caton_time\hodgepodge\gro.go')
  
