 def test(a):
     for k in  a:
         if k==3:
             return k
     return a

 test([1,2,3])
 #will print 3    

 test([0,1,2])
 #print [0, 1, 2]