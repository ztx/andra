package dsl

import (
	"log"

	"github.com/goadesign/goa/dslengine"
	"github.com/ztx/andra"
)

// StorageGroup implements the top level NoSql-andra DSL
// There should be one StorageGroup per Goa application.
func StorageGroup(name string, dsli func()) *andra.StorageGroupDefinition {
	if !dslengine.IsTopLevelDefinition() {
		return nil
	}
	if name == "" {
		dslengine.ReportError("Storage Group name cannot be empty")
	}

	if andra.NoSqlDesign != nil {
		if andra.NoSqlDesign.Name == name {
			dslengine.ReportError("Only one StorageGroup is allowed")
		}
	}
	andra.NoSqlDesign.Name = name
	andra.NoSqlDesign.NoSqlStores = make(map[string]*andra.NoSqlStoreDefinition)
	andra.NoSqlDesign.DefinitionDSL = dsli
	log.Println("in StorageGroup", andra.NoSqlDesign)
	return andra.NoSqlDesign
}

// Description sets the definition description.
// Description can be called inside StorageGroup, NoSqlStore, NoSqlModel, NoSqlField
func Description(d string) {
	if a, ok := storageGroupDefinition(false); ok {
		a.Description = d
	} else if v, ok := isNoSqlStoreDefinition(false); ok {
		v.Description = d
	} else if r, ok := isNoSqlModelDefinition(false); ok {
		r.Description = d
	} else if f, ok := isNoSqlFieldDefinition(false); ok {
		f.Description = d
	}
}
