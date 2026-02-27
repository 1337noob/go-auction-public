package create_lot

import (
	"context"
	"log"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/pkg"
)

type CreateLotHandler struct {
	repo      domain.LotRepository
	eventBus  application.EventBus
	txManager pkg.TransactionManager
}

func NewCreateLotHandler(repo domain.LotRepository, eventBus application.EventBus, txManager pkg.TransactionManager) *CreateLotHandler {
	return &CreateLotHandler{
		repo:      repo,
		eventBus:  eventBus,
		txManager: txManager,
	}
}

// must create transaction in handler
func (h *CreateLotHandler) Handle(ctx context.Context, cmd *CreateLot) error {
	tx, err := h.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	baseCtx := tx.Context()
	ctx = pkg.WithTransaction(baseCtx, tx)

	lot := domain.NewLot(
		cmd.AggregateID,
		cmd.Name,
		cmd.Description,
		cmd.OwnerID,
	)

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
	log.Println("commit from create_lot")

	return nil
}
