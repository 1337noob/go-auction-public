package domain_event_handlers

import (
	"context"
	"main/auction/internal/application/read_model"
	"main/auction/internal/domain"
	"time"
)

type OnAuctionCreatedHandler struct {
	readModelRepo read_model.AuctionReadModelRepo
	lotRepo       domain.LotRepository
}

func NewOnAuctionCreatedHandler(readModelRepo read_model.AuctionReadModelRepo, lotRepo domain.LotRepository) *OnAuctionCreatedHandler {
	return &OnAuctionCreatedHandler{
		readModelRepo: readModelRepo,
		lotRepo:       lotRepo,
	}
}

func (h *OnAuctionCreatedHandler) Handle(ctx context.Context, event domain.Event) error {
	e, ok := event.(*domain.AuctionCreated)
	if !ok {
		return ErrInvalidEvent
	}

	lot, err := h.lotRepo.FindByID(ctx, e.LotID)
	if err != nil {
		return err
	}

	auction := &read_model.AuctionReadModel{
		ID:         e.AggregateID,
		LotID:      e.LotID,
		LotName:    lot.GetName(),
		StartPrice: e.StartPrice,
		MinBidStep: e.MinBidStep,
		SellerID:   e.SellerID,
		CurrentBid: nil,
		Bids:       nil,
		Status:     domain.AuctionStatusCreated,
		StartTime:  e.StartTime,
		EndTime:    e.EndTime,
		Timeout:    e.Timeout,
		CreatedAt:  time.Now(),
		Version:    e.Version,
		UpdatedAt:  e.Timestamp,
	}

	err = h.readModelRepo.Save(ctx, auction)
	if err != nil {
		return err
	}

	return nil
}
