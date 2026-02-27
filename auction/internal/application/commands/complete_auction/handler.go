package complete_auction

import (
	"context"
	"errors"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/pkg"
)

type CompleteAuctionHandler struct {
	repo     domain.AuctionRepository
	eventBus application.EventBus
}

func NewCompleteAuctionHandler(repo domain.AuctionRepository, eventBus application.EventBus) *CompleteAuctionHandler {
	return &CompleteAuctionHandler{
		repo:     repo,
		eventBus: eventBus,
	}
}

// must receive transaction in handler
func (h *CompleteAuctionHandler) Handle(ctx context.Context, cmd CompleteAuction) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}

	auction, err := h.repo.FindByID(ctx, cmd.AggregateID)
	if err != nil {
		return err
	}

	err = auction.CompleteAuction()
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

	return nil
}
