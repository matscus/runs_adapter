package adapter

import (
	"bytes"
	"html/template"
	"reflect"
)

var (
	templateFuncs = template.FuncMap{"rangeStruct": rangeStructer}
	htmlTemplate  = `
	<table border="1">
		<tr>
			<th>ID</th>
			<th>Test Time</th>
			<th>Project</th>
			<th>Scenarios</th>
			<th>Grafanalink</th>
			<th>Description</th>
			<th>Status</th>
		</tr>
		{{ range .}}
			<tr>
				<td>{{ .ID }}</td>
				<td><p><strong>Start: </strong>{{ .StartTime }}</p><p><strong>Stop: </strong>{{ .EndTime }}</p></td>
				<td>{{ .Data.Project }}</td>
				{{ range .Data.Scenarios }}
					<td colspan="1"><p><strong>Operation: </strong> {{ .Name }} </p><p><strong>TPS: </strong>{{ .TPS }}</p><p><strong>SLA: </strong>{{ .SLA }}</p><p><strong>Duration: </strong>{{ .Duration }}</p></td>
				{{ end}}
				<td colspan="1"><p><a href={{ .Data.Grafanalink }}>{{ .ID }}</a></p></td>
				<td>{{ .Data.Description }}</td>
				<td>{{ .Data.Status }}</td>
			</tr>
		{{ end}}
	</table>`
)

func rangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}
	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}
	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}
	return out
}

func getTable(runs []Run) (string, error) {
	var err error
	table := template.New("table").Funcs(templateFuncs)
	table, err = table.Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = table.Execute(buf, runs)
	if err != nil {
		panic(err)
	}
	return buf.String(), nil
}
