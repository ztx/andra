package andra

import (
	"text/template"

	"bitbucket.org/pkg/inflect"

	"github.com/goadesign/goa/goagen/codegen"
)

type CassandraModelWriter struct {
	*codegen.SourceFile
	UserTypeTmpl   *template.Template
	UserHelperTmpl *template.Template
}

func NewCassandraModelWriter(filename string) (*CassandraModelWriter, error) {
	file, err := codegen.SourceFileFor(filename)
	if err != nil {
		return nil, err
	}
	return &CassandraModelWriter{SourceFile: file}, nil
}

// Execute writes the code for the context types to the writer.
func (w *CassandraModelWriter) Execute(data *UserTypeTemplateData) error {
	fm := make(map[string]interface{})
	fm["famt"] = fieldAssignmentModelToType
	fm["fatm"] = fieldAssignmentTypeToModel
	fm["fapm"] = fieldAssignmentPayloadToModel
	fm["viewSelect"] = viewSelect
	fm["viewFields"] = viewFields
	fm["viewFieldNames"] = viewFieldNames
	fm["goDatatype"] = goDatatype
	fm["plural"] = inflect.Pluralize
	fm["gtt"] = codegen.GoTypeTransform
	fm["gttn"] = codegen.GoTypeTransformName
	fm["newMediaTemplate"] = newMediaTemplate
	return w.ExecuteTemplate("casModelT", casModelT, fm, data)
}

const (
	casModelT = `
	{{$ut := .UserType}}{{$ap := .AppPkg}}{{$dp := .DefaultPkg}}

	type {{$ut.ModelName}} struct {
	*{{$dp}}.{{$ut.ModelName}}
}

type SelectQuery struct {
	model           *{{$ut.ModelName}}Model
	selectedColumns []string
	table           string
	whereClause     []string
	values          []interface{}

	//
	cql    string
	cItem  *{{$ut.ModelName}}
	result []{{$ut.ModelName}}
}
{{range $ut.IndexedFields}}
//FilterBy{{.FieldName}} appends a where clause to CQL
func (q *SelectQuery) FilterBy{{.FieldName}}(val {{goDatatype . true}})  *SelectQuery {
	q.whereClause = append(q.whereClause, q.cItem.ColumnOf("{{.FieldName}}")+"=?")
	q.values = append(q.values, val)
	return q
} {{end}}

//Model implements {{$dp}}.{{$ut.ModelName}}Storage for cassandra db
type {{$ut.ModelName}}Model struct {
	{{$dp}}.{{$ut.ModelName}}Storage
}

//
func (m {{$ut.ModelName}}) ColumnOf(attrib string) string {
	out := ""
	switch attrib { {{range $ut.Fields}} 
		case "{{.}}":
			out = "{{$ut.ColumnName .}}"{{end}}
	}
	return out
}

func (m {{$ut.ModelName}}) ColumnsOf(attribs ...string) []string {
	out := []string{}
	for _,s := range attribs {
		out = append(out,m.ColumnOf(s))
	}
	return out
}

func (m *{{$ut.ModelName}}) SelectCQL(attribs ...string) (query string, values []interface{}, err error) {
	//leave the validation of attribute names to ValueHolders
	//as it will panic for invalid attribute names

	values = m.ValueHolders(attribs...)
	columns := make([]string, len(attribs))
	for i, a := range attribs {
		columns[i] = m.ColumnOf(a)
	}
	query = "SELECT " +
		strings.Join(columns, ",") +
		" FROM " +
		m.TableName() +
		" WHERE "

	whereCond := []string{}
	whereValues := []interface{}{}

	for _, k := range m.PrimaryKeys() {
		whereCond = append(whereCond, k+"=?")
		whereValues = append(whereValues, m.ValueHolder(k))
	}
	
	query = query + strings.Join(whereCond, " AND ")	
	values = append(values, whereValues...)
	return
}


//Model returns the containing model
func (m *{{$ut.ModelName}})Model()*{{$dp}}.{{$ut.ModelName}}{
	return m.{{$ut.ModelName}}
}

func (m {{$ut.ModelName}}) InsertCQL() (query string, values []interface{}, err error) {
	query = "INSERT INTO " + m.TableName()
	columns := []string{}
	columnValues := []interface{}{}
	{{range $ut.Fields}}
	if m.{{.}} != nil {
		columns = append(columns, m.ColumnOf("{{.}}"))
		columnValues = append(columnValues, &m.{{.}})
	}
	{{end}}

	query = query + strings.Join(columns, ",") + ")"
	query = query + " VALUES("
	for range columnValues {
		query = query + "?,"
	}
	//remove trailing ","
	query = query[0 : len(query)-1]
	query = query + ")"
	values = columnValues
	return
}

{{define "Media"}}` + mediaT2 + `{{end}}` + `{{$ut := .UserType}}{{$ap := .AppPkg}}
{{ if $ut.Roler }}
// GetRole returns the value of the role field and satisfies the Roler interface.
func (m {{$ut.ModelName}}) GetRole() string {
	return {{$f := $ut.Fields.role}}{{if $f.Nullable}}*{{end}}m.Role
}
{{end}}

{{ range $rname, $rmt := $ut.RenderTo }}
{{ range $vname, $view := $rmt.Views}}
{{ $mtd := $ut.Project $rname $vname }}

{{template "Media" (newMediaTemplate $rmt $vname $view $ut)}}
{{end}}{{end}}

{{ range $idx, $bt := $ut.BelongsTo}}
// Belongs To Relationships

// {{$ut.ModelName}}FilterBy{{$bt.ModelName}} is a gorm filter for a Belongs To relationship.
func {{$ut.ModelName}}FilterBy{{$bt.ModelName}}({{goify (printf "%s%s" $bt.ModelName "ID") false}} int, originaldb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if {{goify (printf "%s%s" $bt.ModelName "ID") false}} > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("{{if $bt.NoSqlFields.ID.DatabaseFieldName}}{{ if ne $bt.NoSqlFields.ID.DatabaseFieldName "id" }}{{$bt.NoSqlFields.ID.DatabaseFieldName}} = ?", {{goify (printf "%s%s" $bt.ModelName "ID") false}}){{else}}{{$bt.Underscore}}_id = ?", {{goify (printf "%s%s" $bt.ModelName "ID") false}}){{end}}
			{{ else }}{{$bt.Underscore}}_id = ?", {{goify (printf "%s%s" $bt.ModelName "ID") false}}){{ end }}
		}
	}
	return func(db *gorm.DB) *gorm.DB { return db }
}
{{end}}

// CRUD Functions

// Get returns a single {{$ut.ModelName}} as a Database Model
// This is more for use internally, and probably not what you want in  your controllers
func (m *{{$ut.ModelName}}Model) Get(ctx context.Context{{ if $ut.DynamicTableName}}, tableName string{{ end }}, {{$ut.PKAttributes}}) (*{{$ut.ModelName}}, error){
	defer goa.MeasureSince([]string{"goa","db","{{goify $ut.ModelName false}}", "get"}, time.Now())

	var native {{$ut.ModelName}}
	err := m.Db.Table({{ if $ut.DynamicTableName }}tableName{{else}}m.TableName(){{ end }}).Where("{{$ut.PKWhere}}",{{$ut.PKWhereFields}} ).Find(&native).Error
	if err ==  gorm.ErrRecordNotFound {
		return nil, err
	}
	{{ if $ut.Cached }}go m.cache.Set(strconv.Itoa(native.ID), &native, cache.DefaultExpiration)
	{{end}}
	return &native, err
}

// List returns an array of {{$ut.ModelName}}
func (m *{{$ut.ModelName}}Model) List(ctx context.Context{{ if $ut.DynamicTableName}}, tableName string{{ end }}) ([]*{{$ut.ModelName}}, error) {
	defer goa.MeasureSince([]string{"goa","db","{{goify $ut.ModelName false}}", "list"}, time.Now())

	var objs []*{{$ut.ModelName}}
	err := m.Db.Table({{ if $ut.DynamicTableName }}tableName{{else}}m.TableName(){{ end }}).Find(&objs).Error
	if err != nil && err !=  gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *{{$ut.ModelName}}Model) Add(ctx context.Context{{ if $ut.DynamicTableName }}, tableName string{{ end }}, model *{{$ut.ModelName}}) (error) {
	defer goa.MeasureSince([]string{"goa","db","{{goify $ut.ModelName false}}", "add"}, time.Now())

{{ range $l, $pk := $ut.PrimaryKeys }}
	{{ if eq $pk.Datatype "uuid" }}model.{{$pk.FieldName}} = uuid.NewV4(){{ end }}
{{ end }}
	err := m.Db{{ if $ut.DynamicTableName }}.Table(tableName){{ end }}.Create(model).Error
	if err != nil {
		goa.LogError(ctx, "error adding {{$ut.ModelName}}", "error", err.Error())
		return err
	}
	{{ if $ut.Cached }}
	go m.cache.Set(strconv.Itoa(model.ID), model, cache.DefaultExpiration) {{ end }}
	return nil
}

// Update modifies a single record.
func (m *{{$ut.ModelName}}Model) Update(ctx context.Context{{ if $ut.DynamicTableName }}, tableName string{{ end }}, model *{{$ut.ModelName}}) error {
	defer goa.MeasureSince([]string{"goa","db","{{goify $ut.ModelName false}}", "update"}, time.Now())

	obj, err := m.Get(ctx{{ if $ut.DynamicTableName }}, tableName{{ end }}, {{$ut.PKUpdateFields "model"}})
	if err != nil {
		goa.LogError(ctx, "error updating {{$ut.ModelName}}", "error", err.Error())
		return  err
	}
	err = m.Db{{ if $ut.DynamicTableName }}.Table(tableName){{ end }}.Model(obj).Updates(model).Error
	{{ if $ut.Cached }}go func(){
		m.cache.Set(strconv.Itoa(model.ID), obj, cache.DefaultExpiration)
	}()
	{{ end }}
	return err
}

// Delete removes a single record.
func (m *{{$ut.ModelName}}Model) Delete(ctx context.Context{{ if $ut.DynamicTableName }}, tableName string{{ end }}, {{$ut.PKAttributes}})  error {
	defer goa.MeasureSince([]string{"goa","db","{{goify $ut.ModelName false}}", "delete"}, time.Now())

	var obj {{$ut.ModelName}}{{ $l := len $ut.PrimaryKeys }}
	{{ if eq $l 1 }}
	err := m.Db{{ if $ut.DynamicTableName }}.Table(tableName){{ end }}.Delete(&obj, {{$ut.PKWhereFields}}).Error
	{{ else  }}err := m.Db{{ if $ut.DynamicTableName }}.Table(tableName){{ end }}.Delete(&obj).Where("{{$ut.PKWhere}}", {{$ut.PKWhereFields}}).Error
	{{ end }}
	if err != nil {
		goa.LogError(ctx, "error deleting {{$ut.ModelName}}", "error", err.Error())
		return  err
	}
	{{ if $ut.Cached }} go m.cache.Delete(strconv.Itoa(id)) {{ end }}
	return  nil
}

{{ range $bfn, $bf := $ut.BuiltFrom }}
// {{$ut.ModelName}}From{{$bfn}} Converts source {{goify $bfn true}} to target {{$ut.ModelName}} model
// only copying the non-nil fields from the source.
func {{$ut.ModelName}}From{{$bfn}}(payload *app.{{goify $bfn true}}) *{{$ut.ModelName}} {
	{{$ut.LowerName}} := &{{$ut.ModelName}}{}
 	{{ fapm $ut $bf "app" "payload" "payload" $ut.LowerName}}

 	 return {{$ut.LowerName}}
}

// UpdateFrom{{$bfn}} applies non-nil changes from {{goify $bfn true}} to the model and saves it
func (m *{{$ut.ModelName}}Model)UpdateFrom{{$bfn}}(ctx context.Context{{ if $ut.DynamicTableName}}, tableName string{{ end }},payload *app.{{goify $bfn true}}, {{$ut.PKAttributes}}) error {
	defer goa.MeasureSince([]string{"goa","db","{{goify $ut.ModelName false}}", "updatefrom{{goify $bfn false}}"}, time.Now())

	var obj {{$ut.ModelName}}
	 err := m.Db.Table({{ if $ut.DynamicTableName }}tableName{{else}}m.TableName(){{ end }}).Where("{{$ut.PKWhere}}",{{$ut.PKWhereFields}} ).Find(&obj).Error
	if err != nil {
		goa.LogError(ctx, "error retrieving {{$ut.ModelName}}", "error", err.Error())
		return  err
	}
 	{{ fapm $ut $bf "app" "payload" "payload" "obj"}}

	err = m.Db.Save(&obj).Error
 	 return err
}
{{ end  }}
	`

	mediaT2 = `// MediaType Retrieval Functions

// List{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}} returns an array of view: {{.ViewName}}.
func (m *{{.Model.ModelName}}Model) List{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}{{/*
*/}} (ctx context.Context{{ if .Model.DynamicTableName}}, tableName string{{ end }}{{/*
*/}} {{range $nm, $bt := .Model.BelongsTo}},{{goify (printf "%s%s" $bt.ModelName "ID") false}} int{{end}}){{/*
*/}} []*app.{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}{
	defer goa.MeasureSince([]string{"goa","db","{{goify .Media.TypeName false}}", "list{{goify .Media.TypeName false}}{{if eq .ViewName "default"}}{{else}}{{goify .ViewName false}}{{end}}"}, time.Now())

	var native []*{{goify .Model.ModelName true}}
	var objs []*app.{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}{{$ctx:= .}}
	err := m.Db.Scopes({{range $nm, $bt := .Model.BelongsTo}}{{/*
*/}}{{$ctx.Model.ModelName}}FilterBy{{goify $bt.ModelName true}}({{goify (printf "%s%s" $bt.ModelName "ID") false}}, m.Db), {{end}}){{/*
*/}}.Table({{ if .Model.DynamicTableName }}tableName{{else}}m.TableName(){{ end }}).{{ range $ln, $lv := .Media.Links }}Preload("{{goify $ln true}}").{{end}}Find(&native).Error
{{/* //	err := m.Db.Table({{ if .Model.DynamicTableName }}tableName{{else}}m.TableName(){{ end }}).{{ range $ln, $lv := .Media.Links }}Preload("{{goify $ln true}}").{{end}}Find(&objs).Error */}}
	if err != nil {
		goa.LogError(ctx, "error listing {{.Model.ModelName}}", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.{{.Model.ModelName}}To{{goify .Media.UserTypeDefinition.TypeName true}}{{if eq .ViewName "default"}}{{else}}{{goify .ViewName true}}{{end}}())
	}

	return objs
}

// {{$.Model.ModelName}}To{{goify .Media.UserTypeDefinition.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}{{/*
*/}} loads a {{.Model.ModelName}} and builds the {{.ViewName}} view of media type {{.Media.TypeName}}.
func (m *{{.Model.ModelName}}) {{$.Model.ModelName}}To{{goify .Media.UserTypeDefinition.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}(){{/*
*/}} *app.{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}} {
	{{.Model.LowerName}} := &app.{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}{}
 	{{ famt .Model .View "m" "m" .Model.LowerName}}

 	 return {{.Model.LowerName}}
}

// One{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}} loads a {{.Model.ModelName}} and builds the {{.ViewName}} view of media type {{.Media.TypeName}}.
func (m *{{.Model.ModelName}}DB) One{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}{{/*
*/}} (ctx context.Context{{ if .Model.DynamicTableName}}, tableName string{{ end }},{{.Model.PKAttributes}}{{/*
*/}}{{range $nm, $bt := .Model.BelongsTo}},{{goify (printf "%s%s" $bt.ModelName "ID") false}} int{{end}}){{/*
*/}} (*app.{{goify .Media.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}, error){
	defer goa.MeasureSince([]string{"goa","db","{{goify .Media.TypeName false}}", "one{{goify .Media.TypeName false}}{{if not (eq .ViewName "default")}}{{goify .ViewName false}}{{end}}"}, time.Now())

	var native {{.Model.ModelName}}
	
	if err != nil && err !=  gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting {{.Model.ModelName}}", "error", err.Error())
		return nil, err
	}
	{{ if .Model.Cached }} go func(){
		m.cache.Set(strconv.Itoa(native.ID), &native, cache.DefaultExpiration)
	}() {{ end }}
	view := *native.{{.Model.ModelName}}To{{goify .Media.UserTypeDefinition.TypeName true}}{{if not (eq .ViewName "default")}}{{goify .ViewName true}}{{end}}()
	return &view, err
}
`
)
