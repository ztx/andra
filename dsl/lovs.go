package dsl

import (
	"fmt"

	"github.com/goadesign/goa/goagen/codegen"
)
import "github.com/ztx/andra"
import "github.com/goadesign/goa/dslengine"

func LOV(name string, typ string, dsl func()) {
	fmt.Println("inside LOV DSL")
	name = codegen.Goify(name, true)
	name = SanitizeFieldName(name)
	if s, ok := isNoSqlStoreDefinition(true); ok {
		lov := andra.NewLOVDefinition()
		lov.Name = name
		lov.Type = typ
		lov.Parent = s
		lov.DefinitionDSL = dsl

		if len(s.LOVs) == 0 {
			s.LOVs = make(map[string]*andra.LOVDefinition)
		}
		s.LOVs[name] = lov
	}

}
func ValidateByLOV(lovName string) {
	if s, ok := isNoSqlFieldDefinition(true); ok {
		lov, ok := s.Parent.Parent.LOVs[lovName]
		if !ok {
			dslengine.ReportError("the LOV %s specified for field %s does not exits", lovName, dslengine.CurrentDefinition().Context())
		} else {
			s.LOV = lov
		}
	}
}
func Value(name string, typ string, val string) {
	l, ok := dslengine.CurrentDefinition().(*andra.LOVDefinition)
	if !ok {
		dslengine.IncompatibleDSL()
	} else {
		v := andra.NewLOVValuesDefinition()

		v.Name = name
		v.Type = typ
		v.Value = val
		v.Parent = l
		if v.Type == "" {
			v.Type = v.Parent.Type
		}
		v.Parent.Values = append(v.Parent.Values, v)
		fmt.Println("in Value DSL, parent is ", v.Parent)
	}
}
