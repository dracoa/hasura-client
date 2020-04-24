package hasura

import (
	"context"
	"github.com/machinebox/graphql"
)

var client *graphql.Client
var secret string

func Init(url string, key string) {
	client = graphql.NewClient(url)
	secret = key
}

type Request struct {
	Raw string
	Req *graphql.Request
}

type MutationResult struct {
	AffectedRows int                      `json:"affected_rows"`
	Returning    []map[string]interface{} `json:"returning"`
}

func BuildRequest(body string) *Request {
	req := graphql.NewRequest(body)
	req.Header.Set("x-hasura-admin-secret", secret)
	r := &Request{Raw: body, Req: req}
	return r
}

func (r *Request) Query() map[string]interface{} {
	ctx := context.Background()
	var respData = make(map[string]interface{})
	if err := client.Run(ctx, r.Req, &respData); err != nil {
		panic(err)
	}
	return respData
}

func (r *Request) Mutate() map[string]MutationResult {
	ctx := context.Background()
	var respData = make(map[string]MutationResult)
	if err := client.Run(ctx, r.Req, &respData); err != nil {
		panic(err)
	}
	return respData
}