package internal

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func InitializeExpenseMCPTools(s *server.MCPServer) {
	calculateShareTool := mcp.NewTool(
		"calculate_share_per_item",
		mcp.WithDescription("Calculate the share amount per person per item based on a given receipt. Based on the receipt and the description given, be sure to take out the items and assign it to a given person. Once you have assigned the item to a given person, NOTE you might have to call this tool multiple times based on the number of items owned by the user, then take the total amount (including tip) and gross amount and pass it to the CalculateShare function. The function will return the share amount for each person per item. This would serve as a entry into the database using the write_query or update_query tool."),
		mcp.WithString("item_associated_with_person",
			mcp.Required(),
			mcp.Description("The item asscoiated with the person, in the format item_name:item_price"),
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
		itemsStr, ok := request.Params.Arguments["item_associated_with_person"].(string)
		if !ok {
			return mcp.NewToolResultError("missing or invalid 'item_associated_with_person' parameter"), nil
		}

		grossAmountStr, ok := request.Params.Arguments["gross_amount"].(string)
		if !ok {
			return mcp.NewToolResultError("missing or invalid 'gross_amount' parameter"), nil
		}
		grossAmountFloat, err := strconv.ParseFloat(grossAmountStr, 32)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid format for 'gross_amount': %v", err)), nil
		}
		grossAmount := float32(grossAmountFloat)

		totalAmountStr, ok := request.Params.Arguments["total_amount"].(string)
		if !ok {
			return mcp.NewToolResultError("missing or invalid 'total_amount' parameter"), nil
		}
		totalAmountFloat, err := strconv.ParseFloat(totalAmountStr, 32)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid format for 'total_amount': %v", err)), nil
		}
		totalAmount := float32(totalAmountFloat)

		itemsMap := make(map[string]float32)
		items := strings.Split(itemsStr, ",")
		var shareTotal float32 = 0

		for _, item := range items {
			parts := strings.Split(item, ":")
			if len(parts) == 2 {
				itemName := parts[0]
				itemPrice, err := strconv.ParseFloat(parts[1], 32)
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("invalid price format for item %s: %v", itemName, err)), nil
				}
				itemsMap[itemName] = float32(itemPrice)
				shareTotal += float32(itemPrice)
			} else {
				return mcp.NewToolResultError(fmt.Sprintf("invalid item format '%s'. Expected 'item_name:item_price'", item)), nil
			}
		}

		if totalAmount <= 0 {
			return mcp.NewToolResultError("total_amount must be greater than zero"), nil
		}

		share := CalculateShare(grossAmount, totalAmount, shareTotal)
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", share)), nil
	})
}

func CalculateShare(grossAmount float32, totalAmount float32, share float32) float32 {
	if totalAmount == 0 {
		return 0
	}
	shareAmount := grossAmount * (share / totalAmount)
	return shareAmount
}
