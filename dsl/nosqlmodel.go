package dsl

import (
	"bitbucket.org/pkg/inflect"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/dslengine"
	"github.com/goadesign/goa/goagen/codegen"
	"github.com/ztx/andra"
)

// Model is the DSL that represents a NoSql Model.
// Model name should be Title cased.  Use BuildsFrom() and RendersTo() DSL
// to define the mapping between a Model and a Goa Type.
// Models may contain multiple instances of the `Field` DSL to
// add fields to the model.
//
// To control whether the ID field is auto-generated, use `NoAutomaticIDFields()`
// Similarly, use NoAutomaticTimestamps() and NoAutomaticSoftDelete() to
// prevent CreatedAt, UpdatedAt and DeletedAt fields from being created.
func Model(name string, dsl func()) {
	if s, ok := isNoSqlStoreDefinition(true); ok {
		var model *andra.NoSqlModelDefinition
		var ok bool
		model, ok = s.NoSqlModels[name]
		if !ok {
			model = andra.NewNoSqlModelDefinition()
			model.ModelName = name
			model.DefinitionDSL = dsl
			model.Parent = s
			model.NoSqlFields = make(map[string]*andra.NoSqlFieldDefinition)
		} else {
			dslengine.ReportError("Model %s already exists", name)
			return
		}
		s.NoSqlModels[name] = model
		model.UserTypeDefinition.TypeName = model.ModelName

	}
}

// RendersTo informs andra that this model will need to be
// rendered to a Goa type.  Conversion functions
// will be generated to convert to/from the model.
//
// Usage: RendersTo(MediaType)
func RendersTo(rt interface{}) {
	if m, ok := isNoSqlModelDefinition(false); ok {
		mts, ok := rt.(*design.MediaTypeDefinition)
		if ok {
			m.RenderTo[mts.TypeName] = mts
		}

	}
}

// BuildsFrom informs andra that this model will be populated
// from a Goa UserType.  Conversion functions
// will be generated to convert from the payload to the model.
//
// Usage:  BuildsFrom(YourType)
//
// Fields not in `YourType` that you want in your model must be
// added explicitly with the `Field` DSL.
func BuildsFrom(dsl func()) {
	if m, ok := isNoSqlModelDefinition(false); ok {
		/*		mts, ok := bf.(*design.UserTypeDefinition)
				if ok {
					m.BuiltFrom[mts.TypeName] = mts
				} else if mts, ok := bf.(*design.MediaTypeDefinition); ok {
					m.BuiltFrom[mts.TypeName] = mts.UserTypeDefinition
				}
				m.PopulateFromModeledType()
		*/
		bf := andra.NewBuildSource()
		bf.DefinitionDSL = dsl
		bf.Parent = m
		m.BuildSources = append(m.BuildSources, bf)
	}

}

// Payload specifies the Resource and Action containing
// a User Type (Payload).
// andra will generate a conversion function for the Payload to
// the Model.
func Payload(r interface{}, act string) {
	if bs, ok := buildSourceDefinition(true); ok {
		var res *design.ResourceDefinition
		var resName string
		if n, ok := r.(string); ok {
			res = design.Design.Resources[n]
			resName = n
		} else {
			res, _ = r.(*design.ResourceDefinition)
		}
		if res == nil {
			dslengine.ReportError("There is no resource %q", resName)
			return
		}
		a, ok := res.Actions[act]
		if !ok {
			dslengine.ReportError("There is no action")
			return
		}
		payload := a.Payload

		// Set UTD in BuildsFrom parent context

		bs.Parent.BuiltFrom[payload.TypeName] = payload
		bs.Parent.PopulateFromModeledType()
	}
}

// BelongsTo signifies a relationship between this model and a
// Parent.  The Parent has the child, and the Child belongs
// to the Parent.
//
// Usage:  BelongsTo("User")
func BelongsTo(parent string) {
	if r, ok := isNoSqlModelDefinition(false); ok {
		idfield := andra.NewNoSqlFieldDefinition()
		idfield.FieldName = codegen.Goify(inflect.Singularize(parent), true) + "ID"
		idfield.Description = "Belongs To " + codegen.Goify(inflect.Singularize(parent), true)
		idfield.Parent = r
		idfield.Datatype = andra.BelongsTo
		idfield.DatabaseFieldName = SanitizeDBFieldName(codegen.Goify(inflect.Singularize(parent), true) + "ID")
		r.NoSqlFields[idfield.FieldName] = idfield
		bt, ok := r.Parent.NoSqlModels[codegen.Goify(inflect.Singularize(parent), true)]
		if ok {
			r.BelongsTo[bt.ModelName] = bt
		} else {
			model := andra.NewNoSqlModelDefinition()
			model.ModelName = codegen.Goify(inflect.Singularize(parent), true)
			model.Parent = r.Parent
			r.BelongsTo[model.ModelName] = model
		}
	}
}

// Alias overrides the name of the SQL store's table or field.
func Alias(d string) {
	if r, ok := isNoSqlModelDefinition(false); ok {
		r.Alias = d
	} else if f, ok := isNoSqlFieldDefinition(false); ok {
		f.DatabaseFieldName = d
	}
}

// Roler sets a boolean flag that cause the generation of a
// Role() function that returns the model's Role value
// Creates a "Role" field in the table if it doesn't already exist
// as a string type
func Roler() {
	if r, ok := isNoSqlModelDefinition(false); ok {
		r.Roler = true
		if _, ok := r.NoSqlFields["Role"]; !ok {
			field := andra.NewNoSqlFieldDefinition()
			field.FieldName = "Role"
			field.Datatype = andra.String
			r.NoSqlFields["Role"] = field
		}
	}
}

// DynamicTableName sets a boolean flag that causes the generator to
// generate function definitions in the database models that specify
// the name of the database table.  Useful when using multiple tables
// with different names but same schema e.g. Users, AdminUsers.
func DynamicTableName() {
	if r, ok := isNoSqlModelDefinition(false); ok {
		r.DynamicTableName = true
	}
}

// CQLTag sets the model's struct tag `sql` value
// for indexing and other purposes.
func CQLTag(d string) {
	if r, ok := isNoSqlModelDefinition(false); ok {
		r.CQLTag = d
	} else if f, ok := isNoSqlFieldDefinition(false); ok {
		f.CQLTag = d
	}
}
