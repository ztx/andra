package andra

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"bitbucket.org/pkg/inflect"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/dslengine"
	"github.com/goadesign/goa/goagen/codegen"
)

// NewNoSqlModelDefinition returns an initialized
// NoSqlModelDefinition.
func NewNoSqlModelDefinition() *NoSqlModelDefinition {
	baseAttr := &design.AttributeDefinition{}
	utd := &design.UserTypeDefinition{
		AttributeDefinition: baseAttr,
	}
	utd.Type = design.Object{}
	m := &NoSqlModelDefinition{
		NoSqlFields: make(map[string]*NoSqlFieldDefinition),
		BuiltFrom:   make(map[string]*design.UserTypeDefinition),
		RenderTo:    make(map[string]*design.MediaTypeDefinition),
		BelongsTo:   make(map[string]*NoSqlModelDefinition),

		UserTypeDefinition: &design.UserTypeDefinition{
			AttributeDefinition: baseAttr,
		},
	}
	return m
}

// Context returns the generic definition name used in error messages.
func (f *NoSqlModelDefinition) Context() string {
	if f.ModelName != "" {
		return fmt.Sprintf("NoSqlModel %#v", f.Name())
	}
	return "unnamed NoSqlModel"
}

// DSL returns this object's DSL.
func (f *NoSqlModelDefinition) DSL() func() {
	return f.DefinitionDSL
}

// TableName returns the TableName of the struct.
func (f NoSqlModelDefinition) TableName() string {
	return inflect.Underscore(inflect.Pluralize(f.ModelName))
}

// Children returns a slice of this objects children.
func (f NoSqlModelDefinition) Children() []dslengine.Definition {
	var stores []dslengine.Definition
	for _, s := range f.NoSqlFields {
		stores = append(stores, s)
	}
	return stores
}

// PKAttributes constructs a pair of field + definition strings
// useful for method parameters.
func (f *NoSqlModelDefinition) PKAttributes() string {
	var attr []string
	for _, pk := range f.PrimaryKeys {
		attr = append(attr, fmt.Sprintf("%s %s", codegen.Goify(pk.DatabaseFieldName, false), goDatatype(pk, true)))
	}
	return strings.Join(attr, ",")
}

// PKWhere returns an array of strings representing the where clause
// of a retrieval by primary key(s) -- x = ? and y = ?.
func (f *NoSqlModelDefinition) PKWhere() string {
	var pkwhere []string
	for _, pk := range f.PrimaryKeys {
		def := fmt.Sprintf("%s = ?", pk.DatabaseFieldName)
		pkwhere = append(pkwhere, def)
	}
	return strings.Join(pkwhere, " and ")
}

// PKWhereFields returns the fields for a where clause for the primary
// keys of a model.
func (f *NoSqlModelDefinition) PKWhereFields() string {
	var pkwhere []string
	for _, pk := range f.PrimaryKeys {
		def := fmt.Sprintf("%s", codegen.Goify(pk.DatabaseFieldName, false))
		pkwhere = append(pkwhere, def)
	}
	return strings.Join(pkwhere, ",")
}

// PKUpdateFields returns something?  This function doesn't look useful in
// current form.  Perhaps it isn't.
func (f *NoSqlModelDefinition) PKUpdateFields(modelname string) string {
	var pkwhere []string
	for _, pk := range f.PrimaryKeys {
		def := fmt.Sprintf("%s.%s", modelname, codegen.Goify(pk.FieldName, true))
		pkwhere = append(pkwhere, def)
	}

	pkw := strings.Join(pkwhere, ",")
	return pkw
}

// StructDefinition returns the struct definition for the model.
func (f *NoSqlModelDefinition) StructDefinition() string {
	header := fmt.Sprintf("type %s struct {\n", f.ModelName)
	var output string
	f.IterateFields(func(field *NoSqlFieldDefinition) error {
		output = output + field.FieldDefinition()
		return nil
	})

	// Get a sortable slice of BelongsTo relationships
	var keys []string
	for k := range f.BelongsTo {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		output = output + f.BelongsTo[k].ModelName + "\t" + f.BelongsTo[k].ModelName + "\n"
	}
	footer := "}\n"
	return header + output + footer
}

// Attribute implements the Container interface of goa.
func (f *NoSqlModelDefinition) Attribute() *design.AttributeDefinition {
	return f.AttributeDefinition
}

// Project does something interesting, and I don't remember if I use it
// anywhere. used in template
func (f *NoSqlModelDefinition) Project(name, v string) *design.MediaTypeDefinition {
	p, _, _ := f.RenderTo[name].Project(v)
	return p
}

// LowerName returns the model name as a lowercase string.
func (f *NoSqlModelDefinition) LowerName() string {
	return codegen.Goify(strings.ToLower(f.ModelName), false)
}

// Underscore returns the model name as a lowercase string in snake case.
func (f *NoSqlModelDefinition) Underscore() string {
	runes := []rune(f.ModelName)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

// IterateBuildSources runs an iterator function once per Model in the Store's model list.
func (f *NoSqlModelDefinition) IterateBuildSources(it BuildSourceIterator) error {

	for _, bs := range f.BuildSources {
		if err := it(bs); err != nil {
			return err
		}
	}
	return nil
}

// IterateFields returns an iterator function useful for iterating through
// this model's field list.
func (f *NoSqlModelDefinition) IterateFields(it FieldIterator) error {
	// Break out each type of fields

	var pkkeys []string
	pks := make(map[string]string)
	for n := range f.NoSqlFields {
		if f.NoSqlFields[n].PrimaryKey {
			pks[n] = n
			pkkeys = append(pkkeys, n)
		}
	}
	sort.Strings(pkkeys)

	var namekeys []string
	names := make(map[string]string)
	for n := range f.NoSqlFields {
		if !f.NoSqlFields[n].PrimaryKey && !f.NoSqlFields[n].Timestamp {
			names[n] = n
			namekeys = append(namekeys, n)
		}
	}
	sort.Strings(namekeys)

	var datekeys []string
	dates := make(map[string]string)
	for n := range f.NoSqlFields {
		if f.NoSqlFields[n].Timestamp {
			dates[n] = n
			datekeys = append(datekeys, n)
		}
	}
	sort.Strings(datekeys)

	// Combine the sorted slices
	var fields []string
	fields = append(fields, pkkeys...)
	fields = append(fields, namekeys...)
	fields = append(fields, datekeys...)

	// Iterate them
	for _, n := range fields {
		if err := it(f.NoSqlFields[n]); err != nil {
			return err
		}
	}
	return nil
}

// PopulateFromModeledType creates fields for the model
// based on the goa UserTypeDefinition it models, which is
// set using BuildsFrom().
// This happens before fields are processed, so it's
// ok to just assign without testing.
func (f *NoSqlModelDefinition) PopulateFromModeledType() {
	if f.BuiltFrom == nil {
		return
	}
	for _, utd := range f.BuiltFrom {
		obj := utd.ToObject()
		obj.IterateAttributes(func(name string, att *design.AttributeDefinition) error {
			rf, ok := f.NoSqlFields[codegen.Goify(name, true)]
			if ok {
				// We already have a mapping for this field.  What to do?
				if rf.Datatype != "" {
					return nil
				}
				// we may have seen the field but don't know its type
				// TODO(BJK) refactor this into separate func later
				switch att.Type.Kind() {
				case design.BooleanKind:
					rf.Datatype = Boolean
				case design.IntegerKind:
					rf.Datatype = Integer
				case design.NumberKind:
					rf.Datatype = Decimal
				case design.StringKind:
					rf.Datatype = String
				case design.DateTimeKind:
					rf.Datatype = Timestamp
				case design.MediaTypeKind:
					// Embedded MediaType
					// Skip for now?
					return nil

				default:
					dslengine.ReportError("Unsupported type: %#v %s", att.Type.Kind(), att.Type.Name())
				}
				if !utd.IsRequired(name) {
					rf.Nullable = true
				}
			}

			rf = &NoSqlFieldDefinition{}
			rf.Parent = f
			rf.FieldName = codegen.Goify(name, true)

			if strings.HasSuffix(rf.FieldName, "Id") {
				rf.FieldName = strings.TrimSuffix(rf.FieldName, "Id")
				rf.FieldName = rf.FieldName + "ID"
			}
			switch att.Type.Kind() {
			case design.BooleanKind:
				rf.Datatype = Boolean
			case design.IntegerKind:
				rf.Datatype = Integer
			case design.NumberKind:
				rf.Datatype = Decimal
			case design.StringKind:
				rf.Datatype = String
			case design.DateTimeKind:
				rf.Datatype = Timestamp
			case design.MediaTypeKind:
				// Embedded MediaType
				// Skip for now?
				return nil

			default:
				dslengine.ReportError("Unsupported type: %#v %s", att.Type.Kind(), att.Type.Name())
			}
			if !utd.IsRequired(name) {
				rf.Nullable = true
			}
			// might need this later?
			rf.a = att
			f.NoSqlFields[rf.FieldName] = rf

			addAttributeToModel(name, att, f)

			return nil
		})
	}
	return
}

func addAttributeToModel(name string, att *design.AttributeDefinition, m *NoSqlModelDefinition) {
	var parent *design.AttributeDefinition
	parent = m.AttributeDefinition
	if parent != nil {
		if parent.Type == nil {
			parent.Type = design.Object{}
		}
		if _, ok := parent.Type.(design.Object); !ok {
			dslengine.ReportError("can't define child attributes on attribute of type %s", parent.Type.Name())
			return
		}

		parent.Type.(design.Object)[name] = att
	}

}

// copied from Goa
func parseAttributeArgs(baseAttr *design.AttributeDefinition, args ...interface{}) (design.DataType, string, func()) {
	var (
		dataType    design.DataType
		description string
		dsl         func()
		ok          bool
	)

	parseDataType := func(expected string, index int) {
		if name, ok := args[index].(string); ok {
			// Lookup type by name
			if dataType, ok = design.Design.Types[name]; !ok {
				if dataType = design.Design.MediaTypeWithIdentifier(name); dataType == nil {
					dslengine.InvalidArgError(expected, args[index])
				}
			}
			return
		}
		if dataType, ok = args[index].(design.DataType); !ok {
			dslengine.InvalidArgError(expected, args[index])
		}
	}
	parseDescription := func(expected string, index int) {
		if description, ok = args[index].(string); !ok {
			dslengine.InvalidArgError(expected, args[index])
		}
	}
	parseDSL := func(index int, success, failure func()) {
		if dsl, ok = args[index].(func()); ok {
			success()
		} else {
			failure()
		}
	}

	success := func() {}

	switch len(args) {
	case 0:
		if baseAttr != nil {
			dataType = baseAttr.Type
		} else {
			dataType = design.String
		}
	case 1:
		success = func() {
			if baseAttr != nil {
				dataType = baseAttr.Type
			}
		}
		parseDSL(0, success, func() { parseDataType("type, type name or func()", 0) })
	case 2:
		parseDataType("type or type name", 0)
		parseDSL(1, success, func() { parseDescription("string or func()", 1) })
	case 3:
		parseDataType("type or type name", 0)
		parseDescription("string", 1)
		parseDSL(2, success, func() { dslengine.InvalidArgError("func()", args[2]) })
	default:
		dslengine.ReportError("too many arguments in call to Attribute")
	}

	return dataType, description, dsl
}

//used to build the structure
func (model NoSqlModelDefinition) Fields() []string {
	// Get a sortable slice of field names
	var keys []string
	for k := range model.NoSqlFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

//used to build Validate() function

func (f *NoSqlModelDefinition) NotNullFieldNamesCheck(v string) string {
	var attr []string
	//pk are not null columns
	for _, pk := range f.PrimaryKeys {
		attr = append(
			attr, v+"."+pk.FieldName+" == nil",
		)
	}
	//slice.Contains()
	return strings.Join(attr, " || ")
}

func (f *NoSqlModelDefinition) LovValidationCheck(v string) string {
	var conds []string

	f.IterateFields(func(field *NoSqlFieldDefinition) error {
		if field.LOV != nil {
			for _, lovValue := range field.LOV.Values {

				conds = append(conds, v+"."+field.FieldName+" != "+lovValue.Name)
			}
			return nil
		}
		return nil
	})
	return strings.Join(conds, " && ")
}

func (f *NoSqlModelDefinition) ColumnName(v string) string {
	out := f.NoSqlFields[v].DatabaseFieldName
	if out == "" {
		out = f.NoSqlFields[v].Underscore()
	}
	return out
}
