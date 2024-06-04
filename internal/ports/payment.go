package ports

import "github.com/3Davydov/ms-order/internal/application/core/domain"

type PaymentPort interface {
	Charge(*domain.Order) error
}
