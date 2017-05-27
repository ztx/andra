package dsl

import (
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/dslengine"
	"github.com/ztx/andra"
)

func init() {
	andra.NoSqlDesign = andra.NewStorageGroupDefinition()
	if andra.NoSqlDesign == nil {
		panic("Root cannot be null andra-cqlgen")
	}
	dslengine.Register(andra.NoSqlDesign)
}

// storageDefinition returns true and current context if it is an StorageGroupDefinition,
// nil and false otherwise.
func storageGroupDefinition(failIfNotSD bool) (*andra.StorageGroupDefinition, bool) {
	a, ok := dslengine.CurrentDefinition().(*andra.StorageGroupDefinition)
	if !ok && failIfNotSD {
		dslengine.IncompatibleDSL()
	}
	return a, ok
}

// NoSqlStoreDefinition returns true and current context if it is an NoSqlStoreDefinition,
// nil and false otherwise.
func isNoSqlStoreDefinition(failIfNotSD bool) (*andra.NoSqlStoreDefinition, bool) {
	a, ok := dslengine.CurrentDefinition().(*andra.NoSqlStoreDefinition)
	if !ok && failIfNotSD {
		dslengine.IncompatibleDSL()
	}
	return a, ok
}

// NoSqlModelDefinition returns true and current context if it is an NoSqlModelDefinition,
// nil and false otherwise.
func isNoSqlModelDefinition(failIfNotSD bool) (*andra.NoSqlModelDefinition, bool) {
	a, ok := dslengine.CurrentDefinition().(*andra.NoSqlModelDefinition)
	if !ok && failIfNotSD {
		dslengine.IncompatibleDSL()
	}
	return a, ok
}

// NoSqlFieldDefinition returns true and current context if it is an NoSqlFieldDefinition,
// nil and false otherwise.
func isNoSqlFieldDefinition(failIfNotSD bool) (*andra.NoSqlFieldDefinition, bool) {
	a, ok := dslengine.CurrentDefinition().(*andra.NoSqlFieldDefinition)
	if !ok && failIfNotSD {
		dslengine.IncompatibleDSL()
	}
	return a, ok
}

// buildSourceDefinition returns true and current context if it is an BuildSource
// nil and false otherwise.
func buildSourceDefinition(failIfNotSD bool) (*andra.BuildSource, bool) {
	a, ok := dslengine.CurrentDefinition().(*andra.BuildSource)
	if !ok && failIfNotSD {
		dslengine.IncompatibleDSL()
	}
	return a, ok
}

// attributeDefinition returns true and current context if it is an AttributeDefinition
// nil and false otherwise.
func attributeDefinition(failIfNotSD bool) (*design.AttributeDefinition, bool) {
	a, ok := dslengine.CurrentDefinition().(*design.AttributeDefinition)
	if !ok && failIfNotSD {
		dslengine.IncompatibleDSL()
	}
	return a, ok
}
