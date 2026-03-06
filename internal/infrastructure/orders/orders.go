package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/Neokrid/order_service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{conn: pool}
}

func (repo *Repository) Save(ctx context.Context, order entity.Order) error {
	query, args, err := squirrel.Insert("orders").
		Columns("id", "user_id", "items", "status", "created_at").
		Values(order.ID, order.UserID, pq.Array(order.Items), order.Status, order.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	fmt.Println("DEBUG SQL:", query)
	fmt.Println("DEBUG ARGS:", args)
	if err != nil {
		return err
	}

	_, err = repo.conn.Exec(ctx, query, args...)
	return err

}

func (repo *Repository) GetByID(ctx context.Context, orderId uuid.UUID) (*entity.Order, error) {
	query, args, err := squirrel.Select("id", "user_id", "items", "status", "created_at").
		From("orders").
		Where(squirrel.Eq{"id": orderId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var o entity.Order
	err = repo.conn.QueryRow(ctx, query, args...).
		Scan(&o.ID, &o.UserID, &o.Items, &o.Status, &o.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("заказ не найден")
		}
		return nil, err
	}
	return &o, nil
}

func (repo *Repository) GetByUserID(ctx context.Context, userId uuid.UUID) ([]*entity.Order, error) {
	query, args, err := squirrel.Select("id", "user_id", "items", "status", "created_at").
		From(`"orders"`).
		Where(squirrel.Eq{"user_id": userId}).
		OrderBy("created_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := repo.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Items, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

func (repo *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query, args, err := squirrel.Update(`"orders"`).
		Set("status", status).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	res, err := repo.conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("не найден заказ для обновления")
	}

	return nil
}
