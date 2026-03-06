package orders

import (
	"errors"

	"github.com/Neokrid/order_service/internal/entity"
	"github.com/Neokrid/order_service/pkg/constants"
	"github.com/Neokrid/order_service/pkg/logger"
	"github.com/Neokrid/order_service/pkg/trx"
)

type Order struct {
	tx        trx.TransactionManager
	logger    logger.Logger
	orderRepo OrdersRepository
	order     entity.Order
}

func NewService(tx trx.TransactionManager, logger logger.Logger, ordersRepository OrdersRepository, order entity.Order) *Order {
	return &Order{
		tx:        tx,
		logger:    logger,
		orderRepo: ordersRepository,
		order:     order,
	}
}

func (o *Order) CanChangeStatus(newStatus string) error {
	if o.order.Status == constants.StatusDone {
		return errors.New("недопустимый  статус")
	}
	if o.order.Status == constants.StatusCreated {
		return errors.New("недопустимый  статус")
	}
	return nil
}
