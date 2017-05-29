package designsvc

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var Item = MediaType("application/vnd.Item+json", func() {
	Description("item details")
	Attributes(func() {
		Attribute("id", Integer, "ID of the item")
		Attribute("code", String, "Item Code")
		Attribute("name", String, "Name of the item")
		Attribute("uom", String, "Unit Of Measure")
		Attribute("href", String, "Api href of item", func() {
			Example("/items/item_code1")
		})

	})
	View("default", func() {
		Attribute("id")
		Attribute("code")
		Attribute("name")

		Attribute("uom")
	})

})

var Pr = MediaType("application/vnd.prHeader+json", func() {
	Description("PR header details")
	Attributes(func() {
		Attribute("prNum", String, "pr number")
		Attribute("prDate", DateTime, "pr date")
		Attribute("href", String, "Api href of item", func() {
			Example("/pr/pr_numb1")
		})

	})
	View("default", func() {
		Attribute("prNum", String, "pr number")
		Attribute("prDate", DateTime, "pr date")
	})

})

var PrDetails = MediaType("application/vnd.PrDetails+json", func() {
	Description("PR details")
	Reference(PrLinePayload)
	Attributes(func() {
		Attribute("id", Integer, "ID of the item")
		Attribute("item", Item, "Item")
		Attribute("itemCode", String, "Item Code")
		Attribute("sl", Integer, "serial number")
		Attribute("qty", Integer, "Quantity")
		Attribute("price", Integer, "price")
		Attribute("href", String, "Api href of item", func() {
			Example("/pr/pr_numb1/prdetails/pr_line1")
		})

	})
	View("default", func() {
		Attribute("id", Integer, "ID of the item")
		Attribute("item", Item, "Item ")
		Attribute("sl", Integer, "serial number")
		Attribute("qty", Integer, "Quantity")
		Attribute("price", Integer, "price")
	})

})
