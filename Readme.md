![Minion](http://octodex.github.com/images/minion.png)
# Happy Goalng Draft
```
      1.how to get fmt string output-->a:=fmt.sprint(object)
      2.defer func() {
        b:=json.marshal(&rsp)
        io.WS(string(b))
       }()#last bracket pass argvs of functon
       3.var TimeFunc = time.Now->donnot worry type
       4.func (m MapClaims) Valid() error->usage m.Valid(), m is a instance struct  
       https://gobyexample.com/interfaces
```
```
    5.func jsonResponse(response interface{}, w http.ResponseWriter) {
 
    json, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return<br/>
    }
 
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    w.Write(json)<br/>}
    write to header




    6.greate interface
    func main() {
    var k []interface{}//it is like list in python
    s := make([]interface{}, 3)
    s[0] = 5
    s[1] = false
    s[2] = "c"
    k = append(k, s...)
    fmt.Println(k) //[5 false c]--slice store different type
    }

    7.var k int
      fmt.Println(k) // will print 0
     type Stu struct {
     hout int
       }
     var stu Stu
     fmt.Println(stu)//will print{0} 

    8.if a function has two  or more return value ,i  donot receive any return_values,it will be 
       correct!if i  want to receive , i must receive all

    9.print manyline shcema in golang
    fmt.Println(`

                              i have monster
                              you know
                              i love
                              you
                                |
                                | 

      `,varname)

     10.change beego file output.go get back to client string info
         func (output *BeegoOutput) JSON(data interface{}, hasIndent bool, coding bool) error {
     output.Header("Content-Type", "application/json; charset=utf-8")
      var content []byte
       var err error
         if hasIndent {
    content, err = json.MarshalIndent(data, "", "  ")
           } else {
    content, err = json.Marshal(data)
           }
         fmt.Println(`
                   I
                  I
       i want to say this is
              return json
                       |
                       |


    `, fmt.Sprint(stringsToJSON(string(content))))
          if err != nil {
    http.Error(output.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
    return err
          }

         if coding {
    content = []byte(stringsToJSON(string(content)))
          }
            // i want to get json string in there

              return output.Body(content)
        }
       

    10 .interface also can initiate!!!
     //var unpacker Unpacker

     /*

     func GetUnpacker() Unpacker {
      return unpacker
     }


     */
     /*
      type Unpacker interface {
      // 根据用户和请求生成Pid
      GeneratePid(uid string, title, notes string) (pid string)
      // 根据pid生成fid
      GenerateFid(pid string, index int, file *cydex.SimpleFile) (fid string, err error)
      // 根据file生成segs
      GenerateSegs(fid string, f *cydex.SimpleFile) (segs []*cydex.Seg, err error)
      // GetPidFromFid, 快速获取pid
      GetPidFromFid(fid string) string
      // Unpack可能会设置参数,这里要保护
      Enter()
      Leave()
     }

     */





```
