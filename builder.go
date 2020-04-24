package hasura

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

const BaseQueryTpl = `query {{.Name}} {{ .Variables }} {
  {{ .Model.Name }} (where: {{ .Wheres }}) {
	{{ range $f := .Model.Fields }}{{ $f.Name }} {{ end }}
  }
}
`

const BaseUpdateTpl = `mutation {{.Name}} {{ .Variables }} {
  update_{{ .Model.Name }} (where: {{ .Wheres }}, _set: $changes) {
	affected_rows
	returning {
		{{ range $f := .Model.Fields }}{{ $f.Name }} {{ end }}
	}
  }
}
`

const BaseInsertTpl = `mutation {{.Name}} {{ .Variables }} {
  insert_{{ .Model.Name }} (objects: $objects) {
	affected_rows
	returning {
		{{ range $f := .Model.Fields }}{{ $f.Name }} {{ end }}
	}
  }
}
`
const BaseDeleteTpl = `mutation {{.Name}} {{ .Variables }} {
  delete_{{ .Model.Name }} (where: {{ .Wheres }}) {
	affected_rows
	returning {
		{{ range $f := .Model.Fields }}{{ $f.Name }} {{ end }}
	}
  }
}
`

type TplContent struct {
	Name      string
	Model     *Model
	Wheres    string
	Variables string
}

func QueryString(m *Model) string {
	return Base(BaseQueryTpl, &TplContent{
		Name:      "BaseQuery",
		Model:     m,
		Variables: Variables(m),
		Wheres:    Wheres(m),
	})
}

func UpdateString(m *Model) string {
	return Base(BaseUpdateTpl, &TplContent{
		Name:      "BaseUpdate",
		Model:     m,
		Variables: Variables(m),
		Wheres:    Wheres(m),
	})
}

func InsertString(m *Model) string {
	return Base(BaseInsertTpl, &TplContent{
		Name:      "BaseInsert",
		Model:     m,
		Variables: Variables(m),
		Wheres:    Wheres(m),
	})
}

func DeleteString(m *Model) string {
	return Base(BaseDeleteTpl, &TplContent{
		Name:      "BaseDelete",
		Model:     m,
		Variables: Variables(m),
		Wheres:    Wheres(m),
	})
}

func Wheres(m *Model) string {
	str := "{"
	for k, v := range m.Wheres {
		str += fmt.Sprintf("%s: {", k)
		parts := make([]string, len(v))
		var i = 0
		for ik, iv := range v {
			sb, _ := json.Marshal(iv)
			var val = string(sb)
			if strings.HasPrefix(val, "\"$") {
				parts[i] = fmt.Sprintf("%s: %v", ik, iv)
			} else {
				parts[i] = fmt.Sprintf("%s: %v", ik, string(sb))
			}
		}
		str += strings.Join(parts, ",")
		str += "}"
	}
	str += "}"
	return str
}

func Variables(m *Model) string {
	if len(m.Variables) == 0 {
		return ""
	}
	parts := make([]string, len(m.Variables))
	var i = 0
	for k, v := range m.Variables {
		parts[i] = fmt.Sprintf("$%s: %s", k, v.Type)
		i++
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
}

func Base(tmpl string, m *TplContent) string {
	t, err := template.New(m.Name).Parse(tmpl)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, m); err != nil {
		panic(err)
	}
	return tpl.String()
}
