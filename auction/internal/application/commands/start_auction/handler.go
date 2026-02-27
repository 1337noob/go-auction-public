package start_auction

import (
	"context"
	"errors"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/pkg"
)

type StartAuctionHandler struct {
	repo     domain.AuctionRepository
	eventBus application.EventBus
}

func NewStartAuctionHandler(repo domain.AuctionRepository, eventBus application.EventBus) *StartAuctionHandler {
	return &StartAuctionHandler{
		repo:     repo,
		eventBus: eventBus,
	}
}

// must receive transaction in handler
func (h *StartAuctionHandler) Handle(ctx context.Context, cmd StartAuction) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}

	auction, err := h.repo.FindByID(ctx, cmd.AggregateID)
	if err != nil {
		return err
	}

	err = auction.StartAuction()
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
