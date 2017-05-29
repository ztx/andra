package design

import (
	. "github.com/goadesign/goa"
	"github.com/ztx/andra"
	. "github.com/ztx/andra/dsl"
	"github.com/ztx/andra/example/designsvc"
)

var _ = StorageGroup("entp", func() {
	Description("This is the global storage group")
	Model("Item", func() {
		Alias("item_master")
		RendersTo(designsvc.Item)
		Description("Model for item master")
		Field("id", andra.Integer, func() {
			PrimaryKey()

		})

		BuildsFrom(func() {
			Payload("item", "create")

		})
		Field("Name", func() {
			Nullable()
		})
	})
	Model("Pr", func() {
		Alias("pr_header")
		RendersTo(designsvc.Pr)
		Description("Model for PR header")
		Field("id", andra.Integer, func() {
			PrimaryKey()

		})
		Field("pr_num", andra.String, func() {
			CQLTag("unique_index")

		})
		Field("approved_qty", andra.Integer)

		BuildsFrom(func() {
			Payload("pr", "create")

		})
	})
	LOV("ItemType", "string", func() {
		Value("FinishedItem", "ItemType", "finishedItem")
		Value("MakeItem", "ItemType", "makeItem")

	})
	LOV("IType", "int", func() {
		Value("IFinishedItem", "", "")
		Value("IMakeItem", "", "")

	})
	Model("PrLine", func() {
		Alias("pr_line")
		RendersTo(designsvc.PrDetails)
		Description("Model for PR lines")

		Field("id", andra.Integer, func() {
			PrimaryKey()

		})

		BuildsFrom(func() {
			Payload("PrLine", "create")

		})
		BelongsTo("Pr")

	})

	//Defining store cassandra
	Store("cassandra", andra.Cassandra, func() {
		Cluster("127.0.0.1")
		KeySpace("storage")
	})
})
