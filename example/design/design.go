package design

import (
	. "github.com/goadesign/goa"
	"github.com/ztx/andra"
	. "github.com/ztx/andra/dsl"
	"github.com/ztx/entp/designsvc"
)

var _ = StorageGroup("entp", func() {
	Description("This is the global storage group")
	Store("cassandra", andra.Cassandra, func() {
		Model("Item", func() {
			Alias("item_master")

			Description("Model for item master")
			Field("id", andra.Integer, func() {
				PrimaryKey()

			})

		})
		Model("Pr", func() {
			Alias("pr_header")

			Description("Model for PR header")
			Field("id", andra.Integer, func() {
				PrimaryKey()

			})
			Field("pr_num", andra.String, func() {
				CQLTag("unique_index")

			})
			Field("approved_qty", andra.Integer)

		})
		Model("PrLine", func() {
			Alias("pr_line")
			RendersTo(designsvc.PrDetails)
			Description("Model for PR lines")

			Field("id", andra.Integer, func() {
				PrimaryKey()

			})

			BelongsTo("Pr")

		})

	})
})
