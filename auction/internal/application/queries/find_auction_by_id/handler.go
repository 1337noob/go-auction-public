package find_auction_by_id

import (
	"context"
	"main/auction/internal/application/read_model"
)

type FindAuctionByIDHandler struct {
	readModelRepository read_model.AuctionReadModelRepo
}

func NewFindAuctionByIDHandler(readModelRepository read_model.AuctionReadModelRepo) *FindAuctionByIDHandler {
	return &FindAuctionByIDHandler{
		readModelRepository: readModelRepository,
	}
}

func (h *FindAuctionByIDHandler) Handle(ctx context.Context, query FindAuctionByID) (*read_model.AuctionReadModel, error) {
	auction, err := h.readModelRepository.FindByID(ctx, query.AuctionID)
	if err != nil {
		return nil, err
	}

	return auction, nil
}
