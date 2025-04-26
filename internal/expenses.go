package internal

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func InitializeExpenseMCPTools(s *server.MCPServer) {
	calculateShareTool := mcp.NewTool(
		"calculate_share",
		mcp.WithDescription("Calculate the share amount per person based on a given receipt. Based on the receipt and the description given, be sure to take out the items and assign it to a given person. Once you have assigned the items to a given person, sum the prices associated with each of the items and then take the total amount (including tip) and gross amount and pass it to the CalculateShare function. The function will return the share amount for each person."),
		mcp.WithString("items_associated_with_person",
			mcp.Required(),
			mcp.Description("The set of items associated with the person and their associated amounts. If an item is shared by multiple people divide that amount by the number of people sharing it. The items are in the format of item_name:amount. For example, 'item1:10.00,item2:20.00' means that item1 is associated with 10.00 and item2 is associated with 20.00."),
		),
		mcp.WithString("gross_amount",
			mcp.Required(),
			mcp.Description("The gross amount of the receipt"),
		),
		mcp.WithString("total_amount",
			mcp.Required(),
			mcp.Description("The total amount of the receipt"),
		),
	)

	s.AddTool(calculateShareTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters
		// Extract parameters using the correct method
		itemsStr := request.Params.Arguments["items_associated_with_person"].(string)
		grossAmount := float32(request.Params.Arguments["gross_amount"].(float64))
		totalAmount := float32(request.Params.Arguments["total_amount"].(float64))

		// Parse the items string into a map
		itemsMap := make(map[string]float32)
		items := strings.Split(itemsStr, ",")
		var shareTotal float32 = 0

		for _, item := range items {
			parts := strings.Split(item, ":")
			if len(parts) == 2 {
				itemName := parts[0]
				itemPrice, err := strconv.ParseFloat(parts[1], 32)
				if err != nil {
					return nil, fmt.Errorf("invalid price for item %s: %v", itemName, err)
				}
				itemsMap[itemName] = float32(itemPrice)
				shareTotal += float32(itemPrice)
			}
		}
	})
}

func CalculateShare(grossAmount float32, totalAmount float32, share float32) float32 {
	shareAmount := grossAmount * (share / totalAmount)
	return shareAmount
}
