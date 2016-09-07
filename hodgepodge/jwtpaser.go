package main

import (
	//"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	mySigningKey = "WOW,MuchShibe,ToDogge"
)

func main() {
	createdToken, err := ExampleNew([]byte(mySigningKey))
	if err != nil {
		fmt.Println("Creating token failed,", err)
		return
	}
	ExampleParse(createdToken, mySigningKey)
}

func ExampleNew(mySigningKey []byte) (string, error) {
	/*
		token := jwt.New(jwt.SigningMethodRS512)
		claims := make(jwt.MapClaims)
		claims["foo"] = "halloaskkakslkkaskljaskjjashhjasjajshjas"
		claims["iat"] = time.Now().Unix()
		token.Claims = claims
	*/

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"exp": time.Now(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err

}

func ExampleParse(myToken string, myKey string) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})

	// if err1 != nil {
	//  panic(err1)
	// }
	//fmt.Println(err1)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//-------get claims's content------------------------------------
		fmt.Println(claims["exp"], claims["iat"])
	} else {
		fmt.Println(err)
	}
	//fmt.Println(token.Claims["exp"])

	if err == nil && token.Valid {
		fmt.Println("Your token is valid.  I like your style.")
	} else {
		fmt.Println("This token is terrible!  I cannot accept this.")
	}
	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.

}
