package place_bid

type PlaceBid struct {
	AggregateID string `json:"aggregate_id"`
	UserID      string `json:"user_id"`
	Amount      int    `json:"amount"`
}
