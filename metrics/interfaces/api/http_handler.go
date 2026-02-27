package api

import (
	"encoding/json"
	"main/metrics/internal/application/queries"
	"net/http"
	"time"
)

type HttpHandler struct {
	queryHandler *queries.GetMetricsHandler
}

func NewHttpHandler(queryHandler *queries.GetMetricsHandler) *HttpHandler {
	return &HttpHandler{
		queryHandler: queryHandler,
	}
}

type GlobalMetricsResponse struct {
	TotalAuctions          int     `json:"total_auctions"`
	CompletedAuctions      int     `json:"completed_auctions"`
	CancelledAuctions      int     `json:"cancelled_auctions"`
	TotalBids              int     `json:"total_bids"`
	AverageAuctionDuration string  `json:"average_auction_duration"`
	AverageFinalPrice      float64 `json:"average_final_price"`
	AverageBidsPerAuction  float64 `json:"average_bids_per_auction"`
	AverageBidAmount       float64 `json:"average_bid_amount"`
}

type AuctionMetricsResponse struct {
	AuctionID   string     `json:"auction_id"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Duration    *string    `json:"duration,omitempty"`
	FinalPrice  *int       `json:"final_price,omitempty"`
	BidCount    int        `json:"bid_count"`
	WinnerID    *string    `json:"winner_id,omitempty"`
	Status      string     `json:"status"`
	StartPrice  int        `json:"start_price"`
	MinBidStep  int        `json:"min_bid_step"`
}

type UserMetricsResponse struct {
	UserID           string     `json:"user_id"`
	TotalBids        int        `json:"total_bids"`
	WonAuctions      int        `json:"won_auctions"`
	TotalBidAmount   int        `json:"total_bid_amount"`
	AverageBidAmount float64    `json:"average_bid_amount"`
	LastBidAt        *time.Time `json:"last_bid_at,omitempty"`
}

func (h *HttpHandler) GetGlobalMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.queryHandler.GetGlobalMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GlobalMetricsResponse{
		TotalAuctions:          metrics.TotalAuctions,
		CompletedAuctions:      metrics.CompletedAuctions,
		CancelledAuctions:      metrics.CancelledAuctions,
		TotalBids:              metrics.TotalBids,
		AverageAuctionDuration: metrics.AverageAuctionDuration.String(),
		AverageFinalPrice:      metrics.AverageFinalPrice,
		AverageBidsPerAuction:  metrics.AverageBidsPerAuction,
		AverageBidAmount:       metrics.AverageBidAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *HttpHandler) GetAuctionMetrics(w http.ResponseWriter, r *http.Request) {
	auctionID := r.URL.Query().Get("auction_id")
	if auctionID == "" {
		http.Error(w, "auction_id is required", http.StatusBadRequest)
		return
	}

	metrics, err := h.queryHandler.GetAuctionMetrics(auctionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if metrics == nil {
		http.Error(w, "auction not found", http.StatusNotFound)
		return
	}

	var durationStr *string
	if metrics.Duration != nil {
		d := metrics.Duration.String()
		durationStr = &d
	}

	response := AuctionMetricsResponse{
		AuctionID:   metrics.AuctionID,
		CreatedAt:   metrics.CreatedAt,
		StartedAt:   metrics.StartedAt,
		CompletedAt: metrics.CompletedAt,
		Duration:    durationStr,
		FinalPrice:  metrics.FinalPrice,
		BidCount:    metrics.BidCount,
		WinnerID:    metrics.WinnerID,
		Status:      metrics.Status,
		StartPrice:  metrics.StartPrice,
		MinBidStep:  metrics.MinBidStep,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *HttpHandler) GetAllAuctionsMetrics(w http.ResponseWriter, r *http.Request) {
	allMetrics, err := h.queryHandler.GetAllAuctionMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]AuctionMetricsResponse, 0, len(allMetrics))
	for _, metrics := range allMetrics {
		var durationStr *string
		if metrics.Duration != nil {
			d := metrics.Duration.String()
			durationStr = &d
		}

		responses = append(responses, AuctionMetricsResponse{
			AuctionID:   metrics.AuctionID,
			CreatedAt:   metrics.CreatedAt,
			StartedAt:   metrics.StartedAt,
			CompletedAt: metrics.CompletedAt,
			Duration:    durationStr,
			FinalPrice:  metrics.FinalPrice,
			BidCount:    metrics.BidCount,
			WinnerID:    metrics.WinnerID,
			Status:      metrics.Status,
			StartPrice:  metrics.StartPrice,
			MinBidStep:  metrics.MinBidStep,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *HttpHandler) GetUserMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	metrics, err := h.queryHandler.GetUserMetrics(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if metrics == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	response := UserMetricsResponse{
		UserID:           metrics.UserID,
		TotalBids:        metrics.TotalBids,
		WonAuctions:      metrics.WonAuctions,
		TotalBidAmount:   metrics.TotalBidAmount,
		AverageBidAmount: metrics.AverageBidAmount,
		LastBidAt:        metrics.LastBidAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *HttpHandler) GetAllUsersMetrics(w http.ResponseWriter, r *http.Request) {
	allMetrics, err := h.queryHandler.GetAllUserMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]UserMetricsResponse, 0, len(allMetrics))
	for _, metrics := range allMetrics {
		responses = append(responses, UserMetricsResponse{
			UserID:           metrics.UserID,
			TotalBids:        metrics.TotalBids,
			WonAuctions:      metrics.WonAuctions,
			TotalBidAmount:   metrics.TotalBidAmount,
			AverageBidAmount: metrics.AverageBidAmount,
			LastBidAt:        metrics.LastBidAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
