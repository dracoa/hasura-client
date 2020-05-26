package hasura

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/machinebox/graphql"
)

func (m *Model) Raw(rawStr string) (map[string]interface{}, error) {
	ctx := context.Background()
	req := graphql.NewRequest(rawStr)
	for k, v := range m.Variables {
		req.Var(k, v.Value)
	}
	req.Header.Set("x-hasura-admin-secret", m.Secret)
	var resp = make(map[string]interface{})
	if err := m.Client.Run(ctx, req, &resp); err != nil {
		panic(err)
	}
	return resp, nil
}

func (m *Model) Query(placeholder interface{}) error {
	ctx := context.Background()
	req := graphql.NewRequest(QueryString(m))
	for k, v := range m.Variables {
		req.Var(k, v.Value)
	}
	req.Header.Set("x-hasura-admin-secret", m.Secret)

	var resp = make(map[string]interface{})
	if err := m.Client.Run(ctx, req, &resp); err != nil {
		panic(err)
	}

	data := resp[m.Name]
	jsondata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsondata, &placeholder); err != nil {
		return err
	}
	return nil
}

func (m *Model) Mutation(query string, placeholder interface{}) (int, error) {
	ctx := context.Background()
	req := graphql.NewRequest(query)
	for k, v := range m.Variables {
		req.Var(k, v.Value)
	}
	req.Header.Set("x-hasura-admin-secret", m.Secret)
	resp := make(map[string]MutationResult)
	if err := m.Client.Run(ctx, req, &resp); err != nil {
		panic(err)
	}

	data := resp[fmt.Sprintf("%s_%s", m.Operation, m.Name)]
	if placeholder != nil {
		jsondata, err := json.Marshal(data.Returning)
		if err != nil {
			return -1, err
		}
		if err = json.Unmarshal(jsondata, &placeholder); err != nil {
			return -1, err
		}
	}
	return data.AffectedRows, nil
}

func (m *Model) UpdateAll(val interface{}, placeholder interface{}) (int, error) {
	return m.updateOp(val, placeholder)
}

func (m *Model) Update(val interface{}, placeholder interface{}) (int, error) {
	if len(m.Wheres) == 0 {
		panic("no where clause for update operation is not permitted, use UpdateAll instead")
	}
	return m.updateOp(val, placeholder)
}

func (m *Model) Insert(val interface{}, placeholder interface{}) (int, error) {
	m.setOperation("insert")
	m.SetVariable("objects", &Inputs{val: val})
	return m.Mutation(InsertString(m), placeholder)
}

func (m *Model) Delete(placeholder interface{}) (int, error) {
	m.setOperation("delete")
	return m.Mutation(DeleteString(m), placeholder)
}

func (m *Model) setOperation(operation string) {
	if m.Operation != "" {
		panic("only one operation is allowed")
	}
	m.Operation = operation
}

func (m *Model) updateOp(val interface{}, placeholder interface{}) (int, error) {
	m.setOperation("update")
	m.SetVariable("changes", &Changes{val: val})
	return m.Mutation(UpdateString(m), placeholder)
}
