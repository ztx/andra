package andra

import (
	"text/template"

	"github.com/goadesign/goa/goagen/codegen"
)

type ModelWriter struct {
	*codegen.SourceFile
	ModelTmpl *template.Template
}

type ModelTmplData struct {
	Cluster  []string
	KeySpace string
}

func NewModelWriter(filename string) (*ModelWriter, error) {
	file, err := codegen.SourceFileFor(filename)
	if err != nil {
		return nil, err
	}
	return &ModelWriter{SourceFile: file}, nil
}

func (w *ModelWriter) Execute(data *ModelTmplData) error {
	fm := make(map[string]interface{})
	fm["wlov"] = writeLOV

	return w.ExecuteTemplate("ModelTmpl", modelDataStoreT, fm, data)
}

const (
	modelDataStoreT = `
	{{$ut := .UserType}}{{$ap := .AppPkg}}// {{if $ut.Description}}{{$ut.Description}}{{else}}{{$ut.ModelName}} NoSql Model{{end}}
{{$ut.StructDefinition}}
// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m {{$ut.ModelName}}) TableName() string {
{{ if ne $ut.Alias "" }}
return "{{ $ut.Alias}}" {{ else }} return "{{ $ut.TableName }}"
{{end}}
}

//ValueHolders return a collection of struct field pointers 
//specified by the param list
func (m *{{$ut.ModelName}}) ValueHolders(attribs ...string) []interface{} {
	var i interface{}
	var val []interface{}
	for _, v := range attribs {
		i = m.ValueHolder(v)
		if i != nil {
			val = append(val, v)
		} else {
			log.Println("Error: No such attribute ", v)
			panic("Error: No such attribute ")
		}
	}
	return val
}

//ValueHolder returns the pointer to struct field identified by name; can be used to
//store the scanned value from db
func (m *{{$ut.ModelName}})ValueHolder(attrib string) interface{} {
var out interface{}
	switch attrib {
		{{range $i, $field := $ut.Fields}}
		case "{{$field}}":
			out=&m.{{$field}}
		{{end}}
	}
return out	
}

//Validate will validate a model
func (m *{{$ut.ModelName}})Validate() bool{
	if {{ $ut.NotNullFieldNamesCheck "m"}}{
		log.Println(errors.New("Some null fields are found which should not be null"))
		return false
	}
	{{$lovValidation:=$ut.LovValidationCheck "m"}}
	{{if ne $lovValidation ""}}if {{$lovValidation}}{
		log.Println(errors.New("Invalid value for an attribute"))
		return false
	}
	{{end}}return true
}

//validates by comparing 2 models
func (im {{$ut.ModelName}}Model) Validate(m1, m2 {{$ut.ModelName}}) bool {
	
	return {{range $i, $pk := $ut.PrimaryKeys}}m1.{{$pk.FieldName}} == m2.{{$pk.FieldName}}&&
	{{end}}  m1.Validate() && m2.Validate() 
		
}

//Returns true for read only attributes
func (item {{$ut.ModelName}}) ReadOnly(attrib string) bool {
	switch attrib {
		{{range $i,$pk:=$ut.PrimaryKeys}}case "{{$pk.FieldName}}":
		return true{{end}}
	}
	return false
}

// {{$ut.ModelName}}Model is the implementation of the storage interface for
// {{$ut.ModelName}}.
type {{$ut.ModelName}}Model struct {
	
}
// New{{$ut.ModelName}}Model creates a new storage type.
func New{{$ut.ModelName}}Model() *{{$ut.ModelName}}Model {
	
	return &{{$ut.ModelName}}Model{}
}

// {{$ut.ModelName}}Storage represents the storage interface.
type {{$ut.ModelName}}Storage interface {
	
	List(ctx context.Context{{ if $ut.DynamicTableName}}, tableName string{{ end }}) ([]*{{$ut.ModelName}}, error)
	Get(ctx context.Context{{ if $ut.DynamicTableName }}, tableName string{{ end }}, {{$ut.PKAttributes}}) (*{{$ut.ModelName}}, error)
	Add(ctx context.Context{{ if $ut.DynamicTableName }}, tableName string{{ end }}, {{$ut.LowerName}} *{{$ut.ModelName}}) (error)
	Update(ctx context.Context{{ if $ut.DynamicTableName }}, tableName string{{ end }}, {{$ut.LowerName}} *{{$ut.ModelName}}) (error)
	Delete(ctx context.Context{{ if $ut.DynamicTableName }}, tableName string{{ end }}, {{ $ut.PKAttributes}}) (error)
{{range $rname, $rmt := $ut.RenderTo}}{{/*

*/}}{{range $vname, $view := $rmt.Views}}
	List{{goify $rmt.TypeName true}}{{if not (eq $vname "default")}}{{goify $vname true}}{{end}} (ctx context.Context{{ if $ut.DynamicTableName}}, tableName string{{ end }}{{/*
*/}}{{range $nm, $bt := $ut.BelongsTo}}, {{goify (printf "%s%s" $bt.ModelName "ID") false}} int{{end}}) []*app.{{goify $rmt.TypeName true}}{{if not (eq $vname "default")}}{{goify $vname true}}{{end}}
	One{{goify $rmt.TypeName true}}{{if not (eq $vname "default")}}{{goify $vname true}}{{end}} (ctx context.Context{{ if $ut.DynamicTableName}}, tableName string{{ end }}{{/*
*/}}, {{$ut.PKAttributes}}{{range $nm, $bt := $ut.BelongsTo}},{{goify (printf "%s%s" $bt.ModelName "ID") false}} int{{end}}){{/*
*/}} (*app.{{goify $rmt.TypeName true}}{{if not (eq $vname "default")}}{{goify $vname true}}{{end}}, error)
{{end}}{{/*

*/}}{{end}}
{{range $bfn, $bf := $ut.BuiltFrom}}
	UpdateFrom{{$bfn}}(ctx context.Context{{ if $ut.DynamicTableName}}, tableName string{{ end }},payload *app.{{goify $bfn true}}, {{$ut.PKAttributes}}) error
{{end}}
}

	`
)
