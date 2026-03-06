package orders

import (
	"context"
	"errors"
	"log"

	"github.com/Neokrid/order_service/internal/domain/service/orders"
	"github.com/Neokrid/order_service/internal/entity"
	"github.com/Neokrid/order_service/pkg/constants"
	"github.com/Neokrid/order_service/pkg/trx"
	"github.com/google/uuid"
)

type OrderService struct {
	repo      orders.OrdersRepository
	publisher orders.OrdersPublisher
	order     OrdersService
	txManager trx.TransactionManager
}

func NewOrdersService(r orders.OrdersRepository, p orders.OrdersPublisher, o OrdersService, txManager trx.TransactionManager) *OrderService {
	return &OrderService{
		repo:      r,
		publisher: p,
		order:     o,
		txManager: txManager,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, items []string) (*entity.Order, error) {
	order := entity.Order{
		ID:     uuid.New(),
		UserID: userID,
		Items:  items,
		Status: constants.StatusCreated,
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return nil, err
	}
	if err := s.publisher.Send(ctx, &order, constants.OrderCreatedEvent); err != nil {

		log.Printf("Критическая ошибка Kafka: %v", err)
		return nil, nil
	}
	return &order, nil
}

func (s *OrderService) GetUserOrder(ctx context.Context, userID, orderID uuid.UUID) (*entity.Order, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("это не ваш заказ!")
	}

	return order, nil
}

func (s *OrderService) GetAllUserOrders(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error) {
	orders, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	return s.txManager.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.order.CanChangeStatus(status); err != nil {
			return err
		}

		if err := s.repo.UpdateStatus(txCtx, orderID, status); err != nil {
			return err
		}

		order, err := s.repo.GetByID(txCtx, orderID)

		if err != nil {
			log.Printf("Ошибка при получении заказа для Kafka: %v", err)
			return err
		}

		if err := s.publisher.Send(txCtx, order, constants.StatusUpdatedEvent); err != nil {

			log.Printf("Критическая ошибка Kafka: %v", err)
			return err
		}
		return nil

	})

}
