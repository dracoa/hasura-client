package hasura

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Model struct {
	Name      string
	Struct    interface{}
	Fields    map[string]Field
	Variables map[string]Variable
	Wheres    map[string]map[string]interface{}
	Operation string
	End       bool
}

type Field struct {
	Name         string
	SkipOnInsert bool
	GoField      reflect.Value
}

type Variable struct {
	Name  string
	Value interface{}
	Type  string
}

type Changes struct {
	val interface{}
}

type Inputs struct {
	val interface{}
}

func (c Changes) GetVal() interface{} {
	return c.val
}

func (c Inputs) GetVal() interface{} {
	return c.val
}

type Input interface {
	GetVal() interface{}
}

func Build(name string, st interface{}) *Model {
	m := &Model{
		Name:      name,
		Struct:    st,
		Variables: make(map[string]Variable),
		Wheres:    make(map[string]map[string]interface{}),
	}
	s := reflect.ValueOf(st).Elem()
	typeOfT := s.Type()
	m.Fields = make(map[string]Field)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fi := Field{
			Name:    typeOfT.Field(i).Tag.Get("json"),
			GoField: f,
		}
		m.Fields[fi.Name] = fi
	}
	return m
}

func (m *Model) Query(placeholder interface{}) error {
	query := QueryString(m)
	req := BuildRequest(query)
	for k, v := range m.Variables {
		req.Req.Var(k, v.Value)
	}
	resp := req.Query()
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
	req := BuildRequest(query)
	for k, v := range m.Variables {
		req.Req.Var(k, v.Value)
	}
	resp := req.Mutate()
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
	return m.UpdateOp(val, placeholder)
}

func (m *Model) Update(val interface{}, placeholder interface{}) (int, error) {
	if len(m.Wheres) == 0 {
		panic("no where clause for update operation is not permitted, use UpdateAll instead")
	}
	return m.UpdateOp(val, placeholder)
}

func (m *Model) UpdateOp(val interface{}, placeholder interface{}) (int, error) {
	m.SetOperation("update")
	m.SetVariable("changes", &Changes{val: val})
	return m.Mutation(UpdateString(m), placeholder)
}

func (m *Model) Insert(val interface{}, placeholder interface{}) (int, error) {
	m.SetOperation("insert")
	m.SetVariable("objects", &Inputs{val: val})
	return m.Mutation(InsertString(m), placeholder)
}

func (m *Model) Delete(placeholder interface{}) (int, error) {
	m.SetOperation("delete")
	return m.Mutation(DeleteString(m), placeholder)
}

func (m *Model) SetOperation(operation string) {
	if m.Operation != "" {
		panic("only one operation is allowed")
	}
	m.Operation = operation
}

func (m *Model) SetVariable(name string, val interface{}) *Model {
	key := fmt.Sprintf("%s", name)
	t := reflect.TypeOf(val)
	gqlType := m.getGQLType(t)
	var v = val
	if t == reflect.TypeOf(&Changes{}) || t == reflect.TypeOf(&Inputs{}) {
		v = val.(Input).GetVal()
	}
	m.Variables[key] = Variable{
		Name:  name,
		Value: v,
		Type:  gqlType,
	}
	return m
}

func (m *Model) SetWhere(name string, operator string, value interface{}) *Model {
	w := m.Wheres[name]
	if w == nil {
		m.Wheres[name] = make(map[string]interface{})
	}
	m.Wheres[name][operator] = value
	return m
}

func (m *Model) getGQLType(goType reflect.Type) string {
	name := fmt.Sprintf("%v", goType)
	switch name {
	case "string":
		return "String"
	case "int":
		return "Int"
	case "time.Time":
		fallthrough
	case "*time.Time":
		return "timestamptz"
	case "*dao.Changes":
		return fmt.Sprintf("%s_set_input", m.Name)
	case "*dao.Inputs":
		return fmt.Sprintf("[%s_insert_input!]!", m.Name)
	}
	return name
}
