package hasura

import (
	"fmt"
	"github.com/machinebox/graphql"
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
	Client    *graphql.Client
	Secret    string
}

type MutationResult struct {
	AffectedRows int                      `json:"affected_rows"`
	Returning    []map[string]interface{} `json:"returning"`
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

type HasuraClient struct {
	Url    string
	Secret string
}

func (h *HasuraClient) Build(name string, st interface{}) *Model {
	m := &Model{
		Name:      name,
		Struct:    st,
		Variables: make(map[string]Variable),
		Wheres:    make(map[string]map[string]interface{}),
		Client:    graphql.NewClient(h.Url),
		Secret:    h.Secret,
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
	case "*hasura.Changes":
		return fmt.Sprintf("%s_set_input", m.Name)
	case "*hasura.Inputs":
		return fmt.Sprintf("[%s_insert_input!]!", m.Name)
	}
	return name
}
