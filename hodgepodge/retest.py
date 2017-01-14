import re 
s="uu00000dd0033"
q="342dss3232rr323"
result=re.findall(r'@\w+.\w+','abc.test@gmail.com, xyz@test.in, test.first@analyticsvidhya.com, first.test@rest.biz')
print(result)
print(re.findall(r'dd|ss\w+rr|ss',q))
a="yangming-bin.2017"
print(re.split(r'.',a))
#will print ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '']
#in this case,you need to Escape
print(re.split(r'\.',a))
#wiil print ['yangming-bin', '2017']