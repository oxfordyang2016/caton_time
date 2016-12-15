import os
a=7
def func():
	a="string"
    #print(a)

func()
print(a)#will print 7



b=[]
def func1():
	b=["yangmingishere"]
	b.append(range(5))
	print(b)#will print ['yangmingishere', range(0, 5)]

func1()
print(b) #will print []


def func2():
	b.append(5)
	print("will print func2 --"+str(b))

func2()	

print(b)
def func3():
	b.append(8)
	print(b)#will print[5,8]

func3()








