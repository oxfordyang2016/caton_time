#reference url:https://www.programiz.com/python-programming/namespace

import os
a=7
def func():
	a="string"
    #print(a)

func()
print(a)#will print 7



b=[]
def func1():
	b=["yangmingishere"]#this b in local namespace
	b.append(range(5))#modify localspace name b ,donot influcence global b
	print(b)#will print ['yangmingishere', range(0, 5)]

func1()
print(b) #will print []


def func2():
	b.append(5)#will modify global b
	print("will print func2 --"+str(b))#print [5]

func2()	

print(b)
def func3():
	b.append(8) #golbal b=[5], and then  modify it to [5,8]
	print(b)#will print[5,8]

func3()








