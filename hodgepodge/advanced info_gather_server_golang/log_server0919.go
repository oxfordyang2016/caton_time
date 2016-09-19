package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type test_struct struct {
	Test string
}
type Response struct {
	Token  string
	Status int
}

var (
	db *sql.DB
)

const (
	mySigningKey = "caton"
	dbconnect    = "root:123456@tcp(192.168.0.81:3306)/hello"
	expire
)

func main() {
	var err1 error
	db, err1 = sql.Open("mysql", dbconnect) //first configure a database
	checkErr(err1)
	http.HandleFunc("/api/v1/login", login)
	http.HandleFunc("/api/v1/test", test2)
	http.HandleFunc("/api/v1/report", report)
	log.Fatal(http.ListenAndServe(":8087", nil))

}
func test2(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "hallo") //write to body
}
func login(rw http.ResponseWriter, req *http.Request) {

	//------------------------paser http request body's json-------------------------------------
	decoder := json.NewDecoder(req.Body)

	type Info struct {
		Sn      string
		Model   string
		Version string
		Pssword string
	}
	var rsp Response
	//-------------defer advanced usage--------------------------------------
	defer func() {
		fmt.Println("rsp.status   ", rsp.Token, "rsp.token  ", rsp.Status)
		b, _ := json.Marshal(&rsp)

		io.WriteString(rw, string(b))
	}()

	var t Info
	err1 := decoder.Decode(&t)
	if err1 != nil {
		rsp.Token = "error"
		rsp.Status = 2
		fmt.Println(err1)
		return
	}

	fmt.Println("json paser success-> ", t.Sn)

	//-----------------------------------retrive info from header------------------------------
	var sn, model, version, date, password string
	for name, headers := range req.Header {
		name := strings.ToLower(name)
		for _, h := range headers {

			if name == "sn" {
				sn = h
			}
			if name == "date" {
				date = h
			}
			if name == "model" {
				model = h
			}
			if name == "version" {
				version = h

			}
			if name == "password" {
				password = h
			}
			fmt.Println("name++++", name)
			fmt.Println("content+++", h)
		}
	}
	fmt.Println(sn, version, date, model, version)

	//-------------------------------------md5 verify login--------------------------------------------------------
	snk := []byte(t.Sn)
	md51 := fmt.Sprintf("%x", md5.Sum(snk))
	k := md51 + mySigningKey + date
	data1 := []byte(k)
	verify_data := fmt.Sprintf("%x", md5.Sum(data1))
	if verify_data == password {
		fmt.Println("login ok")
	} else if verify_data == t.Pssword {
		fmt.Println("login ok")
	} else {
		fmt.Println("when testing,comment below ")
		//rsp.Status = 4
		//rsp.Token = "error"
		//return
	}
	//-------------------------read http body----------------------------------------
	body, _ := ioutil.ReadAll(req.Body)

	fmt.Println("request body's string form", string(body))

	//-----------------------write who to databases table who----------------------------------------
	//lookup database
	rows, err_tmp := db.Query("SELECT * FROM who WHERE sn=?", t.Sn)

	fmt.Print("lookup database row    ", rows)
	if err_tmp != nil {
		rsp.Status = 3
		rsp.Token = "error"
		return
	}
	defer rows.Close()
	fmt.Println("row is  ", rows)
	checkErr(err_tmp)
	found := false
	fmt.Println("start to query database ")
	for rows.Next() {
		var id1 string
		var sn string
		var model string
		var version string
		var last_login string
		err := rows.Scan(&id1, &sn, &model, &version, &last_login)
		checkErr(err)

		fmt.Println(sn)

		if sn == t.Sn {
			found = true
			//-------------------------------update-------------------------------
			stmt, err := db.Prepare("update who set last_login=? where sn=?")

			checkErr(err)
			if err != nil {

				rsp.Status = 3
				rsp.Token = "error"
				return
			}

			last_login := time.Now()
			fmt.Println(last_login)
			res, err := stmt.Exec(last_login, t.Sn)
			if err != nil {
				rsp.Status = 3
				rsp.Token = "error"
				return
			}
			checkErr(err)
			defer stmt.Close()
			affect, err := res.RowsAffected()
			fmt.Println(affect)
			if err != nil {
				rsp.Token = "error"
				rsp.Status = 3
				return
			}
			//-------------break  and return usage---------------------
			break

		}

	}

	//if sn is not in table ,inser it
	if !found {

		stmt, err := db.Prepare("INSERT who SET sn=?,model=?,version=?,last_login=?") //in database create a table called who

		if err != nil {
			rsp.Token = "error"
			rsp.Status = 3
			return
		}

		defer stmt.Close()
		checkErr(err)
		time_login := time.Now()
		res, err := stmt.Exec(t.Sn, t.Model, t.Version, time_login)
		checkErr(err)
		if err != nil {
			rsp.Status = 3
			rsp.Token = "error"
			return
		}

		id, err := res.LastInsertId()

		fmt.Println(id)
		checkErr(err)
		if err != nil {
			rsp.Status = 3
			rsp.Token = "error"
			return
		}
	}
	//--------------return jwt in http response of login_ok------------------------------------->
	to, _ := Newjwt(t.Sn, t.Version, []byte(mySigningKey)) //give jwt token
	rsp.Token = to
	rsp.Status = 1
	return
}

//-----------------------------------------------------------controller------------------------------------------
func report(rw http.ResponseWriter, req *http.Request) {
	var rsp Response
	defer func() {
		b, _ := json.Marshal(&rsp)

		io.WriteString(rw, string(b))
	}()
	fmt.Println("-----------------------------------report start----------------------------------------------------")
	var authorization, ip, content_length string
	time := time.Now()
	fmt.Println(time.Format("20060102150405"))
	ip = strings.Split(req.RemoteAddr, ":")[0] //get ip
	for name, headers := range req.Header {
		name := strings.ToLower(name)
		for _, h := range headers {
			//fmt.Println("name++++", name)

			if name == "authorization" {
				authorization = h

			}
			if name == "content-length" {
				content_length = h
			}
		}
	}
	word := strings.Fields(authorization)
	token := word[1]
	body, _ := ioutil.ReadAll(req.Body)
	//--------------------------------------------jwt verify---------------------------------
	report_sn, report_version := JwtParse(token, "caton")
	fmt.Println("----------------what is ok--------------------------------")
	fmt.Println("report_sn ", report_sn, "report_version  ", report_version)
	if report_sn == "error" {
		rsp.Token = "error"
		rsp.Status = 5
		return
	}

	fmt.Println("paser string success")

	//---------------------database access------------------------------
	var content_type string
	for name, headers := range req.Header {
		name := strings.ToLower(name)
		for _, h := range headers {
			//fmt.Println("name++++", name)
			if name == "content-type" {
				content_type = h
			}

		}
	}
	fmt.Println(content_type)

	//-----------------------------------------------write database-----------------------------------------------
	//https://github.com/go-sql-driver/mysql/wiki/Examples
	// insert
	stmt, err := db.Prepare("INSERT report SET sn=?,ip=?,time=?,version=?,content_type=?,content_length=?,msg_body=?")
	if err != nil {
		fmt.Println(err)
		rsp.Token = "error"
		rsp.Status = 6
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(report_sn, ip, time, report_version, content_type, content_length, body)

	if err != nil {
		rsp.Token = "error"
		rsp.Status = 6
		return
	}

	id, err := res.LastInsertId()
	fmt.Println(id)

	if err != nil {
		rsp.Token = "error"
		rsp.Status = 6
		return
	}
	//--------------------report error----------
	rsp.Token = "success"
	rsp.Status = 1
	return

}

//--------------------------------global error handing mechnism---------------------------
func checkErr(err error) {
	if err != nil {

		fmt.Println(err)

	}
}

//-------------------------------generate a  jwt token-------------------------------------
func Newjwt(sn string, version string, mySigningKey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"exp":     time.Now().Add(time.Second * 1).Unix(),
		"sn":      sn,
		"version": version,
		"iat":     time.Now().Unix(),
	}

	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err

}

//----------------------------------parser jwt token------------------------------------------

func JwtParse(myToken string, myKey string) (string, string) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})

	if err == nil && token.Valid {
		fmt.Println("jwt is valid")
	} else {
		fmt.Println(err)
		return "error", "error"
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		b := fmt.Sprint(claims["sn"])
		c := fmt.Sprint(claims["version"])
		return b, c
	} else {
		fmt.Println(err)
		return "error", "error"

	}

}
