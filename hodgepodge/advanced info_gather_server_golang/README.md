
# Info_gather protocol and api 
## login 

### login method
```
POST /api/v1/login
i.e post http://192.168.0.73/api/v1/login
```

### 1.login args
```
args         type          mustinput     illustration
-----------------------------------------------------
sn           string        yes           equipment_code
model        string        yes          
version      string        yes               
password     string        yes
date         string        yes

Note:
0.please use argvs above as json to http's body!!
1.these args pass to post request
2.about passwd's computing method:password=MD5(MD5(sn)+secret+date)
3.secret="caton"
4.cnvert md5(sn) to string,ie,string(md5("halloworld")) or a := fmt.Sprintf("%x", md5.Sum(sn)) //lowercase
5.lowercase all string
6.md5 computing_method example:
    snk := []byte("hallo")//convert string to byte
	md51 := fmt.Sprintf("%x", md5.Sum(snk))//convert md5 code to string
	fmt.Println("first md51 is ", md51)
	k := md51 + "caton" + "date"//string append
	fmt.Println("join part string ", k)
	data1 := []byte(k)
	verify_data := fmt.Sprintf("%x", md5.Sum(data1))
	fmt.Println("final md5 is ", verify_data)//convert byte to string

   output:
   first md51 is  598d4c200461b81522a3328565c25f7c
   join part string  598d4c200461b81522a3328565c25f 
   final md5 is  41d21f5dca8ce86e8b512fe0b46037af

```
### 2.login return  data
```
return       type         data                                                                      illustration
---------------------------------------------------------------------------------------------------------------------------------------------------------
Token        string        refer to below                                                           when you get this token,save it.it will be used agian                       
Status        int          1 login success  2 json arg error    3 server error    4 password error


Note:
0.the return_ data will be written to response body.
1.return_data's success form is 
  {"Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.ZMmY1db_M4sTzNSMjaJDRNK9rxc9kTxLzEo861NS4Vs","Status":1}
2.login password error {"Token":"error","Status":4}
3.json arg error {"Token":"error","Status":2}
3.servererror {"Token":"error","Status":3}

```
### 3.Storing data example_table who
```
+------------+--------------+------+-----+---------+----------------+
| Field      | Type         | Null | Key | Default | Extra          |
+------------+--------------+------+-----+---------+----------------+
| id         | int(11)      | NO   | PRI | NULL    | auto_increment |
| sn         | varchar(255) | YES  |     | NULL    |                |
| model      | varchar(255) | YES  |     | NULL    |                |
| version    | varchar(255) | YES  |     | NULL    |                |
| last_login | varchar(255) | YES  |     | NULL    |                |
+------------+--------------+------+-----+---------+----------------+
```

# report
### 1.report method
```
POST /api/v1/report
```
### 2.use token in http header 

```
---------------------------------------------------------
Note:use token  received above and put it in http header
The content of the header should look like the following:
---------------------------------------------------------
authorization: Bearer <token>

example 
au="Bearer"+" "+json_obj['Token']//get response's token
headers = {"content-type": "application/x-www-form-urlencoded", "authorization":au}

```
### 3.return json data
```
---------------------------------------------------
1 report  success
5 jwt token error
6 server error

Note:
0.json data in http body,it's form is that {"Token":"success","Status":1}
1.report ok {"Token":"success","Status":1}
1.jwt token error {"Token":"error","Status":5}
2.server error {"Token":"error","Status":6}
3.when returned status is 5,please login again 
4.Test_method is manual.txt


```
### 4.Storing data exmaple_table report
```
+----------------+--------------+------+-----+---------+----------------+
| Field          | Type         | Null | Key | Default | Extra          |
+----------------+--------------+------+-----+---------+----------------+
| id             | int(11)      | NO   | PRI | NULL    | auto_increment |
| sn             | varchar(255) | YES  |     | NULL    |                |
| ip             | varchar(255) | YES  |     | NULL    |                |
| time           | varchar(255) | YES  |     | NULL    |                |
| version        | varchar(255) | YES  |     | NULL    |                |
| content_type   | varchar(255) | YES  |     | NULL    |                |
| content_length | int(11)      | YES  |     | NULL    |                |
| msg_body       | longblob     | YES  |     | NULL    |                |
+----------------+--------------+------+-----+---------+----------------+
```
### 5.database create table
```
 create database hello

 create table who( id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, sn VARCHAR(255), model VARCHAR(255), version VARCHAR(255), last_login VARCHAR(255))

 create table report( id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, sn VARCHAR(255), ip VARCHAR(255),time VARCHAR(255), version VARCHAR(255),content_type VARCHAR(255), content_length INT(11), msg_body longblob)

 GRANT ALL ON hello.* TO root@'%' IDENTIFIED BY '123456';
```








