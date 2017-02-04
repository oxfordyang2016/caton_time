 
#case 1
 def smart_divide(func):
    def inner(a,b):
       print("I am going to divide",a,"and",b)
       if b == 0:
          print("Whoops! cannot divide")
          return

       return func#return function name
    return inner

#case 2
 def smart_divide(func):
    def inner(a,b):
       print("I am going to divide",a,"and",b)
       if b == 0:
          print("Whoops! cannot divide")
          return

       return func(a,b)#return result
    return inner


def divide(a,b):
    return a/b





# note:return funcion name and function result is different
'''
Case 1
In [42]: a=smart_divide(divide)
In [44]: a(6,2)
I am going to divide 6 and 2
Out[44]: <function __main__.divide>


In [45]: d=a(6,2)
I am going to divide 6 and 2


In [46]: d(6,2)
Out[46]: 3.0

'''


'''
In [49]: a(6,5)
I am going to divide 6 and 5
Out[49]: 1.2

'''


#you must note:func(a,b) is a value(it is also none),but func is a function