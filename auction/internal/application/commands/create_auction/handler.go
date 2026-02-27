package create_auction

import (
	"context"
	"errors"
	"log"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/pkg"
)

type CreateAuctionHandler struct {
	repo      domain.AuctionRepository
	lotRepo   domain.LotRepository
	eventBus  application.EventBus
	txManager pkg.TransactionManager
}

func NewCreateAuctionHandler(repo domain.AuctionRepository, lotRepo domain.LotRepository, eventBus application.EventBus, txManager pkg.TransactionManager) *CreateAuctionHandler {
	return &CreateAuctionHandler{
		repo:      repo,
		lotRepo:   lotRepo,
		eventBus:  eventBus,
		txManager: txManager,
	}
}

// must create transaction in handler
func (h *CreateAuctionHandler) Handle(ctx context.Context, cmd CreateAuction) error {
	tx, err := h.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	baseCtx := tx.Context()
	ctx = pkg.WithTransaction(baseCtx, tx)

	log.Println("lot id:", cmd.LotID)
	lot, err := h.lotRepo.FindByID(ctx, cmd.LotID)
	if err != nil {
		return err
	}
	if lot.GetStatus() != domain.LotStatusPublished {
		return errors.New("lot not published")
	}

	auction, err := domain.NewAuction(
		cmd.AggregateID,
		cmd.LotID,
		cmd.StartPrice,
		cmd.MinBidStep,
		cmd.SellerID,
		cmd.StartTime,
		cmd.EndTime,
		cmd.Timeout,
	)
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
	log.Println("commit from create_auction")

	return nil
}
