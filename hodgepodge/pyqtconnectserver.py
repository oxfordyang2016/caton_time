# -*- coding: utf-8 -*-

# Form implementation generated from reading ui file 'pyqtconnectserver.ui'
#
# Created by: PyQt5 UI code generator 5.4.1
#
# WARNING! All changes made in this file will be lost!

from PyQt5 import QtCore, QtGui, QtWidgets
import http.client, urllib.parse,json
import threading
from random import randint
from time import sleep
class Ui_Dialog(object):
    def setupUi(self, Dialog):
        Dialog.setObjectName("Dialog")
        Dialog.resize(1025, 644)
        self.buttonBox = QtWidgets.QDialogButtonBox(Dialog)
        self.buttonBox.setGeometry(QtCore.QRect(220, 540, 341, 32))
        self.buttonBox.setOrientation(QtCore.Qt.Horizontal)
        self.buttonBox.setStandardButtons(QtWidgets.QDialogButtonBox.Cancel|QtWidgets.QDialogButtonBox.Ok)
        self.buttonBox.setObjectName("buttonBox")
        self.textBrowser = QtWidgets.QTextBrowser(Dialog)
        self.textBrowser.setGeometry(QtCore.QRect(490, 91, 501, 391))
        self.textBrowser.setObjectName("textBrowser")
        self.pushButton = QtWidgets.QPushButton(Dialog)
        self.pushButton.setGeometry(QtCore.QRect(670, 540, 75, 23))
        self.pushButton.setObjectName("pushButton")

        self.retranslateUi(Dialog)
        self.buttonBox.accepted.connect(Dialog.accept)
        self.buttonBox.rejected.connect(Dialog.reject)
        QtCore.QMetaObject.connectSlotsByName(Dialog)

    def retranslateUi(self, Dialog):
        _translate = QtCore.QCoreApplication.translate
        Dialog.setWindowTitle(_translate("Dialog", "Dialog"))
        self.pushButton.setText(_translate("Dialog", "connect"))
def serverconnect():
    params = urllib.parse.urlencode({'@number': 12524, '@type': 'issue', '@action': 'show'})
    headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
    conn = http.client.HTTPConnection("192.168.0.68",9000)
    for num in range(0,1):
        data={'sn':str(num)+str(9),'model':'model','version':'vesion','password':'password'}
        data1 = json.dumps(data)
        print('hello')
        conn.request("GET", "/api/v1/nodes", data1, headers)#the third argv is body
        response = conn.getresponse()
        print (response.status, response.reason)
        data = response.read()
        print(data)
        string1 = data.decode('utf-8')
        print(string1)
        print(str(string1))
        #print(type(string(data)))
        json_obj = json.loads(string1)  
        ui.textBrowser.setText(str(json_obj))
    return json_obj    

if __name__ == "__main__":
    import sys
    app = QtWidgets.QApplication(sys.argv)
    Dialog = QtWidgets.QDialog()
    ui = Ui_Dialog()
    ui.setupUi(Dialog)
    ui.pushButton.clicked.connect(serverconnect)#set signal and slot
   
    Dialog.show()
    sys.exit(app.exec_())

