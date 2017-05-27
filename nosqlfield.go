package andra

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/dslengine"
)

// NewNoSqlFieldDefinition returns an initialized
// NoSqlFieldDefinition.
func NewNoSqlFieldDefinition() *NoSqlFieldDefinition {
	m := &NoSqlFieldDefinition{
		Mappings: make(map[string]*MapDefinition),
	}
	return m
}

// Context returns the generic definition name used in error messages.
func (f *NoSqlFieldDefinition) Context() string {
	if f.FieldName != "" {
		return fmt.Sprintf("NoSqlField %#v", f.FieldName)
	}
	return "unnamed NoSqlField"
}

// DSL returns this object's DSL.
func (f *NoSqlFieldDefinition) DSL() func() {
	return f.DefinitionDSL
}

// Children returns a slice of this objects children.
func (f NoSqlFieldDefinition) Children() []dslengine.Definition {
	// no children yet
	return []dslengine.Definition{}
}

// Attribute implements the Container interface of the goa Attribute
// model.
func (f *NoSqlFieldDefinition) Attribute() *design.AttributeDefinition {
	return f.a
}

// FieldDefinition returns the field's struct definition.
func (f *NoSqlFieldDefinition) FieldDefinition() string {
	var comment string
	if f.Description != "" {
		comment = "// " + f.Description
	}
	def := fmt.Sprintf("%s\t%s %s %s\n", f.FieldName, goDatatype(f, true), tags(f), comment)
	return def
}

// Tags returns the sql and gocql struct tags for the Definition.
func (f *NoSqlFieldDefinition) Tags() string {
	return tags(f)
}

// LowerName returns the field name as a lowercase string.
func (f *NoSqlFieldDefinition) LowerName() string {
	return strings.ToLower(f.FieldName)
}

// Underscore returns the field name as a lowercase string in snake case.
func (f *NoSqlFieldDefinition) Underscore() string {
	runes := []rune(f.FieldName)
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

func goDatatype(f *NoSqlFieldDefinition, includePtr bool) string {
	var ptr string
	if includePtr {
		ptr = "*"
	}
	switch f.Datatype {
	case Boolean:
		return ptr + "bool"
	case Integer:
		return ptr + "int"
	case BigInteger:
		return ptr + "int64"
	case AutoInteger, AutoBigInteger:
		return ptr + "int " // gocql tags later
	case Decimal:
		return ptr + "float32"
	case BigDecimal:
		return ptr + "float64"
	case String:
		return ptr + "string"
	case Text:
		return ptr + "string"
	case UUID:
		return ptr + "uuid.UUID"
	case Timestamp, NullableTimestamp:
		return ptr + "time.Time"
	case BelongsTo:
		return ptr + "int"

	}

	return "UNKNOWN TYPE"
}

func tags(f *NoSqlFieldDefinition) string {
	var cqltags []string
	if f.CQLTag != "" {
		cqltags = append(cqltags, f.CQLTag)
	}

	var gocqltags []string
	if f.DatabaseFieldName != "" && f.DatabaseFieldName != f.Underscore() {
		gocqltags = append(gocqltags, "column:"+f.DatabaseFieldName)
	}
	if f.PrimaryKey {
		gocqltags = append(gocqltags, "primary_key")
	}

	var tags []string
	if len(cqltags) > 0 {
		sqltag := "sql:\"" + strings.Join(cqltags, ";") + "\""
		tags = append(tags, sqltag)
	}
	if len(gocqltags) > 0 {
		gocqltag := "gocql:\"" + strings.Join(gocqltags, ";") + "\""
		tags = append(tags, gocqltag)
	}

	if len(tags) > 0 {
		return "`" + strings.Join(tags, " ") + "`"
	}
	return ""
}
