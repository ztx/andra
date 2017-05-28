package andra

import (
	"fmt"

	"github.com/goadesign/goa/dslengine"
)

// NewNoSqlStoreDefinition returns an initialized
// NoSqlStoreDefinition.
func NewNoSqlStoreDefinition() *NoSqlStoreDefinition {
	m := &NoSqlStoreDefinition{
		NoSqlModels: make(map[string]*NoSqlModelDefinition),
	}
	return m
}

// Context returns the generic definition name used in error messages.
func (sd *NoSqlStoreDefinition) Context() string {
	if sd.Name != "" {
		return fmt.Sprintf("NoSqlStore %#v", sd.Name)
	}
	return "unnamed NoSqlStore"
}

// DSL returns this object's DSL.
func (sd *NoSqlStoreDefinition) DSL() func() {
	return sd.DefinitionDSL
}

// Children returns a slice of this objects children.
func (sd NoSqlStoreDefinition) Children() []dslengine.Definition {
	var stores []dslengine.Definition
	for _, s := range sd.NoSqlModels {
		stores = append(stores, s)
	}
	return stores
}
