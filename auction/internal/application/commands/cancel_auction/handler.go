package cancel_auction

import (
	"context"
	"log"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/pkg"
)

type CancelAuctionHandler struct {
	repo      domain.AuctionRepository
	eventBus  application.EventBus
	txManager pkg.TransactionManager
}

func NewCancelAuctionHandler(repo domain.AuctionRepository, eventBus application.EventBus, txManager pkg.TransactionManager) *CancelAuctionHandler {
	return &CancelAuctionHandler{
		repo:      repo,
		eventBus:  eventBus,
		txManager: txManager,
	}
}

// must create transaction in handler
func (h *CancelAuctionHandler) Handle(ctx context.Context, cmd CancelAuction) error {
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

	err = auction.CancelAuction(cmd.SellerID, cmd.Reason)
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
	log.Println("commit from cancel_auction")

	return nil
}
