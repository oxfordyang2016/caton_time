def getinfo(func):
    def inner(*args,**kwargs):
        print('you are input args is the below')
        for k in  args:
            print(k)
        func(*args,**kwargs)
    return inner
'''
IMHO,the return inner  is a function name
when you use getinfo(test1),function getinfo(test1)  returns function name inner 
and you use a=getinfo(test1)---it is equivant to a=inner.
when you use a(3,4,5),it do inner(3,4,5)!!!!!!!!!

the decorator is equiant to the above idea!!!!!!!
'''
@getinfo
def test(a,b):
    print('i love you'+str(a+b))


'''
In [68]: test(3,4)
you are input args is the below
3
4
i love you7
'''        

def test1(a,b,c):
    print(a+b+c)

a=getinfo(test1)

'''
In [71]: a(3,4,5)
you are input args is the below
3
4
5
12
'''  

