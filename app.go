package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber"
)

type Tag struct {
	username string `json: "username"`
	dob      string `json:"dob"`
	age      int    `json:"age"`
	email    string `json:"email"`
	phone    int    `json:"phone"`
}

var (
	conn context.Context
)

var user Tag

func main() {
	app := fiber.New()

	//app.Use(middlewere.Logger())

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("helo welcome")
	})

	app.Post("/auth", Auth)
	app.Get("/user/profile", Profile)
	app.Get("/microservice/name", func(ctx *fiber.Ctx) {
		//securedStatus := ctx.Body()
		ctx.Send("user-microservice")
	})
	app.Post("/proxy", Proxy)
	app.Listen(": 3000")
	//if err != nil  err *

}

func Proxy(ctx *fiber.Ctx) {
	headers := ctx.Get("username")
	securedStatus := ctx.Body()
	check := []byte(securedStatus)
	mp := make(map[string]bool)
	err := json.Unmarshal(check, &mp)

	if err != nil {
		panic(err)
	}

	var url string
	if len(headers) > 0 && mp["Secured"] {
		url = "http://localhost:3000/auth"
		checkStatus := true
		Method := "POST"
		body := apiCall(url, checkStatus, Method)
		if body == "200" {
			url = "http://localhost:3000/user/profile"
			Method = "GET"
			body = apiCall(url, checkStatus, Method)
			ctx.Send(body)
		} else {
			url = "http://localhost:3000/microservice/name"
			checkStatus := true
			Method := "GET"
			body := apiCall(url, checkStatus, Method)
			ctx.Send(body)

		}

	} else {
		url = "http://localhost:3000/microservice/name"
		checkStatus := true
		Method := "GET"
		body := apiCall(url, checkStatus, Method)
		ctx.Send(body)
	}

}

func Auth(ctx *fiber.Ctx) {

	headers := ctx.Get("username")
	fmt.Println(headers)
	checkUserName := "loki"
	v := headers != checkUserName
	if !v {
		ctx.JSON(200)
	} else {
		ctx.JSON(401)
	}

}

func Profile(ctx *fiber.Ctx) {
	header := ctx.Get("username")
	fmt.Println(header)
	db, err := sql.Open("mysql", "root:Loki@123#@tcp(127.0.0.1:3306)/test")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("inside else")
	}

	user, err := getUserDetails(db, header)
	jsonString, _ := json.Marshal(user)

	if err != nil {
		fmt.Println("NO DATA FOR THE USER")
		ctx.Send(err)
	} else {
		ctx.Send(jsonString)
	}
}

func getUserDetails(db *sql.DB, username string) (user map[string]interface{}, err error) {
	mp := make(map[string]interface{}, 0)
	mps := make(map[string]interface{}, 0)
	sql := "select username, dob, age, email, phone from users where `username` = ?"
	res, err := db.Query(sql, username)
	defer res.Close()
	if err != nil {
		return mps, err
	}

	//

	var pUser Tag
	for res.Next() {
		err = res.Scan(&pUser.username, &pUser.dob, &pUser.age, &pUser.email, &pUser.phone)
		mp["username"] = pUser.username
		mp["dob"] = pUser.dob
		mp["age"] = pUser.age
		mp["email"] = pUser.email
		mp["phone"] = pUser.phone
		fmt.Println(pUser)
		if err != nil {
			return mps, err
		}
	}

	return mp, nil
}

func apiCall(url string, checkStatus bool, Method string) string {

	var jsonStr = []byte(`{"Secured":false}`)

	if checkStatus {
		jsonStr = []byte(`{"Secured":true}`)
	}

	req, err := http.NewRequest(""+Method+"", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("username", "loki ")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err condition true")
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)

}
