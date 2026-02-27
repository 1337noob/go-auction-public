package cancel_auction

type CancelAuction struct {
	AggregateID string `json:"aggregate_id"`
	SellerID    string `json:"seller_id"`
	Reason      string `json:"reason"`
}
