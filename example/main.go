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
	h := &hasura.HasuraClient{
		Url:    "http://127.0.0.1:8484/v1/graphql",
		Secret: "4HTSS9CWbnBuR49yAuhhaqSJUSQG8wjSb4bUkQTtgNrQ2RfvBmTLLe35V8vBxeE6",
	}
	users := make([]User, 0)
	result, err := h.Build("user", &User{}).Raw(`query BaseQuery  {
  user (where: {}) {
	id name 
  }
}`)
	if err != nil {
		log.Panic(err)
	}
	log.Println(result, users)
}
