package andra

// andraDesign is the root definition for andra
var NoSqlDesign *StorageGroupDefinition

const (
	// StorageGroup is the constant string used as the index in the
	// andraConstructs map
	StorageGroup = "storagegroup"
	// Cassandra is the StorageType for Cassandra databases
	Cassandra NoSqlStorageType = "cassanra"
	// ElasticSearch is the StorageType for ElasticSearch
	// ElasticSearch is not implemented
	ElasticSearch NoSqlStorageType = "elasticsearch"
	// MongoDb is the StorageType for MongoDb databases
	// MongoDb is not implemented
	MongoDb NoSqlStorageType = "mongodb"
	// None is For tests
	None NoSqlStorageType = ""
	// Boolean is a bool field type
	Boolean FieldType = "bool"
	// Integer is an integer field type
	Integer FieldType = "integer"
	// BigInteger is a large integer field type
	BigInteger FieldType = "biginteger"
	// AutoInteger is not implemented
	AutoInteger FieldType = "auto_integer"
	// AutoBigInteger is not implemented
	AutoBigInteger FieldType = "auto_biginteger"
	// Decimal is a float field type
	Decimal FieldType = "decimal"
	// BigDecimal is a large float field type
	BigDecimal FieldType = "bigdecimal"
	// String is a varchar field type
	String FieldType = "string"
	// Text is a large string field type
	Text FieldType = "text"
	// UUID is not implemented yet
	UUID FieldType = "uuid"
	// Timestamp is a date/time field in the database
	Timestamp FieldType = "timestamp"
	// NullableTimestamp is a timestamp that may not be
	// populated.  Fields with no value will be null in the database
	NullableTimestamp FieldType = "nulltimestamp"
	// NotFound is used internally
	NotFound FieldType = "notfound"
	// HasOne is used internally
	HasOne FieldType = "hasone"
	// HasOneKey is used internally
	HasOneKey FieldType = "hasonekey"
	// HasMany is used internally
	HasMany FieldType = "hasmany"
	// HasManyKey is used internally
	HasManyKey FieldType = "hasmanykey"
	// Many2Many is used internally
	Many2Many FieldType = "many2many"
	// Many2ManyKey is used internally
	Many2ManyKey FieldType = "many2manykey"
	// BelongsTo is used internally
	BelongsTo FieldType = "belongsto"
)
