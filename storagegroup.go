package andra

import (
	"fmt"
	"sort"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/dslengine"
)

// NewStorageGroupDefinition returns an initialized
// StorageGroupDefinition.
func NewStorageGroupDefinition() *StorageGroupDefinition {
	m := &StorageGroupDefinition{
		NoSqlStores: make(map[string]*NoSqlStoreDefinition),
	}
	return m
}

// IterateStores runs an iterator function once per NoSql Store in the
// StorageGroup's Store list.
func (sd *StorageGroupDefinition) IterateStores(it StoreIterator) error {
	if sd == nil {
		return nil
	}
	if sd.NoSqlStores != nil {
		names := make([]string, len(sd.NoSqlStores))
		i := 0
		for n := range sd.NoSqlStores {
			names[i] = n
			i++
		}
		sort.Strings(names)
		for _, n := range names {
			if err := it(sd.NoSqlStores[n]); err != nil {
				return err
			}
		}
	}
	return nil
}

// Context returns the generic definition name used in error messages.
func (sd StorageGroupDefinition) Context() string {
	if sd.Name != "" {
		return fmt.Sprintf("StorageGroup %#v", sd.Name)
	}
	return "unnamed Storage Group"
}

// DSL returns this object's DSL.
func (sd StorageGroupDefinition) DSL() func() {
	return sd.DefinitionDSL
}

// Children returns a slice of this objects children.
func (sd StorageGroupDefinition) Children() []dslengine.Definition {
	var stores []dslengine.Definition
	for _, s := range sd.NoSqlStores {
		stores = append(stores, s)
	}
	return stores
}

// DSLName is displayed to the user when the DSL executes.
func (sd *StorageGroupDefinition) DSLName() string {
	return "Gorma storage group"
}

// DependsOn return the DSL roots the Gorma DSL root depends on, that's the goa API DSL.
func (sd *StorageGroupDefinition) DependsOn() []dslengine.Root {
	return []dslengine.Root{design.Design, design.GeneratedMediaTypes}
}

// IterateSets goes over all the definition sets of the StorageGroup: the
// StorageGroup definition itself, each store definition, models and fields.
func (sd *StorageGroupDefinition) IterateSets(iterator dslengine.SetIterator) {
	// First run the top level StorageGroup

	iterator([]dslengine.Definition{sd})
	sd.IterateStores(func(store *NoSqlStoreDefinition) error {
		iterator([]dslengine.Definition{store})
		store.IterateModels(func(model *NoSqlModelDefinition) error {
			iterator([]dslengine.Definition{model})
			model.IterateFields(func(field *NoSqlFieldDefinition) error {
				iterator([]dslengine.Definition{field})
				return nil
			})
			model.IterateBuildSources(func(bs *BuildSource) error {
				iterator([]dslengine.Definition{bs})
				return nil
			})

			return nil
		})
		return nil
	})
}

// Reset resets the storage group to pre DSL execution state.
func (sd *StorageGroupDefinition) Reset() {
	n := NewStorageGroupDefinition()
	*sd = *n
}
