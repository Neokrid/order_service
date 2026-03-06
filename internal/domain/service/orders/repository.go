package orders

import (
	"context"

	"github.com/Neokrid/order_service/internal/entity"
	"github.com/google/uuid"
)

type OrdersRepository interface {
	Save(ctx context.Context, order entity.Order) error
	GetByID(ctx context.Context, orderId uuid.UUID) (*entity.Order, error)
	GetByUserID(ctx context.Context, userId uuid.UUID) ([]*entity.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type OrdersPublisher interface {
	Send(ctx context.Context, order *entity.Order, eventType string) error
}
