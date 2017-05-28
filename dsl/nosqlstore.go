package dsl

import (
	"github.com/goadesign/goa/dslengine"
	"github.com/ztx/andra"
)

// Store represents a database.  andra lets you define
// a database type
func Store(name string, storeType andra.NoSqlStorageType, dsl func()) {
	if name == "" || len(name) == 0 {
		dslengine.ReportError("NoSql Store requires a name.")
		return
	}
	if len(storeType) == 0 {
		dslengine.ReportError("NoSql Store requires a NoSqlStoreType.")
		return
	}
	if dsl == nil {
		dslengine.ReportError("NoSql Store requires a dsl.")
		return
	}
	if s, ok := storageGroupDefinition(true); ok {
		if s.NoSqlStores == nil {
			s.NoSqlStores = make(map[string]*andra.NoSqlStoreDefinition)
		}
		store, ok := s.NoSqlStores[name]
		if !ok {
			store = &andra.NoSqlStoreDefinition{
				Name:          name,
				DefinitionDSL: dsl,
				Parent:        s,
				Type:          storeType,
				NoSqlModels:   make(map[string]*andra.NoSqlModelDefinition),
			}
		} else {
			dslengine.ReportError("NoSql Store %s can only be declared once.", name)
		}
		s.NoSqlStores[name] = store
	}

}

func Cluster(nodes ...string) {
	if s, ok := isNoSqlStoreDefinition(true); ok {
		s.Cluster = nodes
	} else {
		dslengine.IncompatibleDSL()
		return
	}
}

//KeySpace is required to get a cassandra session
func KeySpace(keySpace string) {
	if s, ok := isNoSqlStoreDefinition(true); ok {
		if s.Type != andra.Cassandra {
			dslengine.ReportError("KeySpace is valid only for a Cassandra NoSQL store")
			return
		}
		s.KeySpace = keySpace
	} else {
		dslengine.IncompatibleDSL()
		return
	}
}

// // NoAutomaticIDFields applies to a `Store` or `Model` type.  It allows you
// // to turn off the default behavior that will automatically create
// // an ID/int Primary Key for each model.
// func NoAutomaticIDFields() {
// 	if s, ok := isNoSqlStoreDefinition(false); ok {
// 		s.NoAutoIDFields = true
// 	} else if m, ok := isNoSqlModelDefinition(false); ok {
// 		delete(m.NoSqlFields, "ID")
// 	}
// }

// // NoAutomaticTimestamps applies to a `Store` or `Model` type.  It allows you
// // to turn off the default behavior that will automatically create
// // an `CreatedAt` and `UpdatedAt` fields for each model.
// func NoAutomaticTimestamps() {
// 	if s, ok := isNoSqlStoreDefinition(false); ok {
// 		s.NoAutoTimestamps = true
// 	} else if m, ok := isNoSqlModelDefinition(false); ok {
// 		delete(m.NoSqlFields, "CreatedAt")
// 		delete(m.NoSqlFields, "UpdatedAt")
// 	}
// }

// // NoAutomaticSoftDelete applies to a `Store` or `Model` type.  It allows
// // you to turn off the default behavior that will automatically
// // create a `DeletedAt` field (*time.Time) that acts as a
// // soft-delete filter for your models.
// func NoAutomaticSoftDelete() {
// 	if s, ok := isNoSqlStoreDefinition(false); ok {
// 		s.NoAutoSoftDelete = true
// 	} else if m, ok := isNoSqlModelDefinition(false); ok {
// 		delete(m.NoSqlFields, "DeletedAt")
// 	}
// }
