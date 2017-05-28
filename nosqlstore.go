package andra

import (
	"fmt"
	"sort"

	"github.com/goadesign/goa/dslengine"
)

// NewNoSqlStoreDefinition returns an initialized
// NoSqlStoreDefinition.
func NewNoSqlStoreDefinition() *NoSqlStoreDefinition {
	m := &NoSqlStoreDefinition{
		NoSqlModels: make(map[string]*NoSqlModelDefinition),
		LOVs:        make(map[string]*LOVDefinition),
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

// IterateModels runs an iterator function once per Model in the Store's model list.
func (sd *NoSqlStoreDefinition) IterateModels(it ModelIterator) error {
	names := make([]string, len(sd.NoSqlModels))
	i := 0
	for n := range sd.NoSqlModels {
		names[i] = n
		i++
	}
	sort.Strings(names)
	for _, n := range names {
		if err := it(sd.NoSqlModels[n]); err != nil {
			return err
		}
	}
	return nil
}

// IterateLovs runs an iterator function once per LOV in the Store's LOV list.
func (sd *NoSqlStoreDefinition) IterateLOVs(it LovIterator) error {
	names := make([]string, len(sd.LOVs))
	i := 0
	for n := range sd.LOVs {
		names[i] = n
		i++
	}
	sort.Strings(names)
	for _, n := range names {
		if err := it(sd.LOVs[n]); err != nil {
			return err
		}
	}
	return nil
}
