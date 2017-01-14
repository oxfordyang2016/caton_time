 def test():     
     for k in range(10):
         print k
         yield k
         print('i am test yied')
         print('i am function end')  


#it will output

'''
0
0
i am test yied
1
1
i am test yied
2
2
i am test yied
3
3
i am test yied
4
4
i am test yied
5
5
i am test yied
6
6
i am test yied
7
7
i am test yied
8
8
i am test yied
9
9
i am test yied
i am function end
'''