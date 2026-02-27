package publish_lot

import (
	"context"
	"log"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/pkg"
)

type PublishLotHandler struct {
	repo      domain.LotRepository
	eventBus  application.EventBus
	txManager pkg.TransactionManager
}

func NewPublishLotHandler(repo domain.LotRepository, eventBus application.EventBus, txManager pkg.TransactionManager) *PublishLotHandler {
	return &PublishLotHandler{
		repo:      repo,
		eventBus:  eventBus,
		txManager: txManager,
	}
}

// must create transaction in handler
func (h *PublishLotHandler) Handle(ctx context.Context, cmd *PublishLot) error {
	tx, err := h.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	baseCtx := tx.Context()
	ctx = pkg.WithTransaction(baseCtx, tx)
	lot, err := h.repo.FindByID(ctx, cmd.AggregateID)
	if err != nil {
		return err
	}

	err = lot.Publish()
	if err != nil {
		return err
	}

	events := lot.GetUncommitedEvents()

	err = h.repo.Save(ctx, lot)
	if err != nil {
		return err
	}

	for _, event := range events {
		err = h.eventBus.Publish(ctx, event)
		if err != nil {
			return err
		}
	}

	lot.ClearUncommitedEvents()

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Println("commit from publish_lot")

	return nil
}
