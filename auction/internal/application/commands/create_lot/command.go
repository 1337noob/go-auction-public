package create_lot

type CreateLot struct {
	AggregateID string `json:"aggregate_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
}
