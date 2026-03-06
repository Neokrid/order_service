package orders

type OrdersService interface {
	CanChangeStatus(newStatus string) error
}
