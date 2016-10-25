![Minion](http://octodex.github.com/images/minion.png)
# Happy Goalng
```
      1.how to get fmt output-->a:=fmt.sprint(object)
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
    fmt.Println(k)

    }
```
