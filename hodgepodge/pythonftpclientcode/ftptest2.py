import ftplib
session = ftplib.FTP('192.168.0.236','test','test')
file = open(r'C:\Users\dell\Desktop\goproject\caton_time\go.pdf','rb')                  # file to send
session.storbinary('STOR '+'hu.d', file)     # send the file
file.close()                                    # close file and FTP
session.quit()
