package main

import (
	"github.com/dracoa/hasura-client"
	"log"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	UserRole struct {
		Role string `json:"role"`
	} `json:"user_role"`
}

func main() {
	h := &hasura.HasuraClient{
		Url:    "http://127.0.0.1:8484/v1/graphql",
		Secret: "4HTSS9CWbnBuR49yAuhhaqSJUSQG8wjSb4bUkQTtgNrQ2RfvBmTLLe35V8vBxeE6",
	}
	model := h.Build("user", &User{}).
		SetWhere("id", "_eq", "mary").
		SetWhere("user_role", "", "")
	log.Println(hasura.QueryString(model))

	//.Query(&users)
	//if err != nil {
	//	log.Panic(err)
	//}
	//user := users[0]
	//user.Name = "Mary Man"
	//i, err := h.Build("user", &User{}).SetWhere("id", "_eq", "mary").
	//	Update(map[string]string{"name": "Mary Man"}, &user)
	//log.Println(i, user)
}
