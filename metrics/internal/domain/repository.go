package domain

type MetricsRepository interface {
	SaveAuctionMetrics(metrics *AuctionMetrics) error
	GetAuctionMetrics(auctionID string) (*AuctionMetrics, error)
	GetAllAuctionMetrics() ([]*AuctionMetrics, error)

	UpdateUserMetrics(userID string, updateFn func(*UserMetrics)) error
	GetUserMetrics(userID string) (*UserMetrics, error)
	GetAllUserMetrics() ([]*UserMetrics, error)

	GetGlobalMetrics() (*GlobalMetrics, error)
}
