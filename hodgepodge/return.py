def test(a):
    for k in  a:
        if k==3:
            print(k) 
            return k#return finish,this iteration immediately end
    return a

 test([1,2,3])
 #will print 1,2,3 ,ouput 3    

 test([0,1,2])
 #print 0, 1, 2 output [0,1,2]

test([1,24,5,6,73,3,4,5,6])
#will print 1
#24
#5
#6
#73
#3
#output 3