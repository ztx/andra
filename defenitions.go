package andra

import (
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/dslengine"
)

// NoSqlStorageType is the type of database.
type NoSqlStorageType string

// FieldType is the storage data type for a database field.
type FieldType string

// StorageGroupDefinition is the parent configuration structure for nosql definitions.
type StorageGroupDefinition struct {
	dslengine.Definition
	DefinitionDSL func()
	Name          string
	Description   string
	NoSqlStores   map[string]*NoSqlStoreDefinition
	NoSqlModels   map[string]*NoSqlModelDefinition
}

// NoSqlStoreDefinition is the parent configuration structure for andra NoSql model definitions.
type NoSqlStoreDefinition struct {
	dslengine.Definition
	DefinitionDSL func()
	Name          string
	Description   string
	Parent        *StorageGroupDefinition
	Type          NoSqlStorageType
	NoSqlModels   map[string]*NoSqlModelDefinition
	LOVs          map[string]*LOVDefinition
}

// NoSqlModelDefinition implements the storage of a domain model into a
// table in a NoSql database.
type NoSqlModelDefinition struct {
	dslengine.Definition
	*design.UserTypeDefinition
	DefinitionDSL    func()
	ModelName        string
	Description      string
	GoaType          *design.MediaTypeDefinition
	Parent           *StorageGroupDefinition
	BuiltFrom        map[string]*design.UserTypeDefinition
	BuildSources     []*BuildSource
	RenderTo         map[string]*design.MediaTypeDefinition
	BelongsTo        map[string]*NoSqlModelDefinition
	NoSqlStores      map[string]*NoSqlStoreDefinition
	Alias            string // gocql:tablename
	Cached           bool
	CacheDuration    int
	Roler            bool
	DynamicTableName bool
	CQLTag           string
	NoSqlFields      map[string]*NoSqlFieldDefinition
	PrimaryKeys      []*NoSqlFieldDefinition
	PartitionKeys    []*NoSqlFieldDefinition
	ReadOnlyFields   []*NoSqlFieldDefinition
	ClusterKeys      []*NoSqlFieldDefinition
}

// BuildSource stores the BuildsFrom sources
// for parsing.
type BuildSource struct {
	dslengine.Definition
	DefinitionDSL   func()
	Parent          *NoSqlModelDefinition
	BuildSourceName string
}

// MapDefinition represents field mapping to and from
// andra models.
type MapDefinition struct {
	RemoteType  *design.UserTypeDefinition
	RemoteField string
}

// MediaTypeAdapterDefinition represents the transformation of a
// Goa media type into a nosql Model.
//
// Unimplemented at this time.
type MediaTypeAdapterDefinition struct {
	dslengine.Definition
	DefinitionDSL func()
	Name          string
	Description   string
	Left          *design.MediaTypeDefinition
	Right         *NoSqlModelDefinition
}

// UserTypeAdapterDefinition represents the transformation of a Goa
// user type into a nosql Model.
//
// Unimplemented at this time.
type UserTypeAdapterDefinition struct {
	dslengine.Definition
	DefinitionDSL func()
	Name          string
	Description   string
	Left          *NoSqlModelDefinition
	Right         *NoSqlModelDefinition
}

// PayloadAdapterDefinition represents the transformation of a Goa
// Payload (which is really a UserTypeDefinition)
// into a nosql model.
//
// Unimplemented at this time.
type PayloadAdapterDefinition struct {
	dslengine.Definition
	DefinitionDSL func()
	Name          string
	Description   string
	Left          *design.UserTypeDefinition
	Right         *NoSqlModelDefinition
}

// NoSqlFieldDefinition represents
// a field in a NoSql database.
type NoSqlFieldDefinition struct {
	dslengine.Definition
	DefinitionDSL     func()
	Parent            *NoSqlModelDefinition
	a                 *design.AttributeDefinition
	FieldName         string
	TableName         string
	Datatype          FieldType
	CQLTag            string
	DatabaseFieldName string // gocql:column
	Description       string
	Nullable          bool
	PrimaryKey        bool
	Timestamp         bool
	Size              int // string field size
	BelongsTo         string
	PartitionKey      bool
	ClusterKey        bool
	ReadOnly          bool
	Mappings          map[string]*MapDefinition

	LOV *LOVDefinition
}

//LOV definition s
type LOVDefinition struct {
	dslengine.Definition
	DefinitionDSL func()
	Parent        *NoSqlStoreDefinition
	Name          string
	Type          string
	Values        []*LOVValueDefinition
}

type LOVValueDefinition struct {
	dslengine.Definition
	Parent *LOVDefinition
	Name   string
	Type   string
	Value  string
}

// StoreIterator is a function that iterates over NoSql Stores in a
// StorageGroup.
type StoreIterator func(m *NoSqlStoreDefinition) error

// ModelIterator is a function that iterates over Models in a
// NoSqlStore.
type ModelIterator func(m *NoSqlModelDefinition) error

// FieldIterator is a function that iterates over Fields
// in a NoSqlModel.
type FieldIterator func(m *NoSqlFieldDefinition) error

// BuildSourceIterator is a function that iterates over Fields
// in a NoSqlModel.
type BuildSourceIterator func(m *BuildSource) error

type LovIterator func(m *LOVDefinition) error
