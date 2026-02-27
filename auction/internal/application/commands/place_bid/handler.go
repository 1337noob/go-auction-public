package place_bid

import (
	"context"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/pkg"
)

type PlaceBidHandler struct {
	repo      domain.AuctionRepository
	eventBus  application.EventBus
	txManager pkg.TransactionManager
}

func NewPlaceBidHandler(repo domain.AuctionRepository, eventBus application.EventBus, txManager pkg.TransactionManager) *PlaceBidHandler {
	return &PlaceBidHandler{
		repo:      repo,
		eventBus:  eventBus,
		txManager: txManager,
	}
}

// must create transaction in handler
func (h *PlaceBidHandler) Handle(ctx context.Context, cmd PlaceBid) error {
	tx, err := h.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	baseCtx := tx.Context()
	ctx = pkg.WithTransaction(baseCtx, tx)

	auction, err := h.repo.FindByID(ctx, cmd.AggregateID)
	if err != nil {
		return err
	}

	err = auction.PlaceBid(cmd.UserID, cmd.Amount)
	if err != nil {
		return err
	}

	events := auction.GetUncommitedEvents()

	err = h.repo.Save(ctx, auction)
	if err != nil {
		return err
	}

	for _, event := range events {
		err = h.eventBus.Publish(ctx, event)
		if err != nil {
			return err
		}
	}

	auction.ClearUncommitedEvents()

	err = tx.Commit()
	if err != nil {
		return err
	}
	//log.Println("commit from place_bid")

	return nil
}
