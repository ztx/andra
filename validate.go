package andra

import (
	"fmt"

	"github.com/goadesign/goa/dslengine"
)

// Validate tests whether the StorageGroup definition is consistent.
func (a *StorageGroupDefinition) Validate() *dslengine.ValidationErrors {
	fmt.Println("Validating Group")
	verr := new(dslengine.ValidationErrors)
	if a.Name == "" {
		verr.Add(a, "storage group name not defined")
	}
	a.IterateStores(func(store *NoSqlStoreDefinition) error {
		verr.Merge(store.Validate())
		return nil
	})

	return verr.AsError()
}

// Validate tests whether the NoSqlStore definition is consistent.
func (a *NoSqlStoreDefinition) Validate() *dslengine.ValidationErrors {
	fmt.Println("Validating Store")
	verr := new(dslengine.ValidationErrors)
	if a.Name == "" {
		verr.Add(a, "store name not defined")
	}
	if a.Parent == nil {
		verr.Add(a, "missing storage group parent")
	}
	a.IterateModels(func(model *NoSqlModelDefinition) error {
		verr.Merge(model.Validate())
		return nil
	})

	return verr.AsError()
}

// Validate tests whether the NoSqlModel definition is consistent.
func (a *NoSqlModelDefinition) Validate() *dslengine.ValidationErrors {
	fmt.Println("Validating Model")
	verr := new(dslengine.ValidationErrors)
	if a.ModelName == "" {
		verr.Add(a, "model name not defined")
	}
	if a.Parent == nil {
		verr.Add(a, "missing NoSql store parent")
	}
	a.IterateFields(func(field *NoSqlFieldDefinition) error {
		verr.Merge(field.Validate())
		return nil
	})

	return verr.AsError()
}

// Validate tests whether the NoSqlField definition is consistent.
func (field *NoSqlFieldDefinition) Validate() *dslengine.ValidationErrors {
	fmt.Println("Validing Field")
	verr := new(dslengine.ValidationErrors)

	if field.Parent == nil {
		verr.Add(field, "missing NoSql model parent")
	}
	if field.FieldName == "" {
		verr.Add(field, "field name not defined")
	}
	return verr.AsError()
}
