package internal

func CalculateShare(grossAmount float32, totalAmount float32, share float32) float32 {
	shareAmount := grossAmount * (share / totalAmount)
	return shareAmount
}
