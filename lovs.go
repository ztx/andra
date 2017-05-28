package andra

import (
	"fmt"

	"github.com/goadesign/goa/dslengine"
)

func NewLOVDefinition() *LOVDefinition {
	return &LOVDefinition{}
}

func NewLOVValuesDefinition() *LOVValueDefinition {
	return &LOVValueDefinition{}
}

func (f *LOVDefinition) Context() string {
	if f.Name != "" {
		return fmt.Sprintf("LOV %#v", f.Name)
	}
	return "unnamed LOV"
}

// DSL returns this object's DSL.
func (f *LOVDefinition) DSL() func() {
	return f.DefinitionDSL
}

// Children returns a slice of this objects children.
func (f LOVDefinition) Children() []dslengine.Definition {
	var stores []dslengine.Definition
	for _, s := range f.Values {
		stores = append(stores, s)
	}
	return stores
}

// DSLName is displayed to the user when the DSL executes.
func (sd *LOVDefinition) DSLName() string {
	return "Andra LOV for field Validation"
}

func (f *LOVValueDefinition) Context() string {
	if f.Name != "" {
		return fmt.Sprintf("LOV Value %#v", f.Name)
	}
	return "unnamed LOV Value"
}
