package main

import (
	"github.com/dracoa/hasura-client"
	"log"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	hasura.Init("http://127.0.0.1:8484/v1/graphql", "4HTSS9CWbnBuR49yAuhhaqSJUSQG8wjSb4bUkQTtgNrQ2RfvBmTLLe35V8vBxeE6")
	model := hasura.Build("user", &User{})
	users := make([]User, 0)
	_ = model.SetWhere("id", "_eq", "peter").
		SetWhere("password", "_eq", "iampeter").
		Query(&users)
	log.Println(users)
}
