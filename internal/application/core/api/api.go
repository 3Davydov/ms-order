package api

import (
	"context"

	"github.com/3Davydov/ms-order/internal/application/core/domain"
	"github.com/3Davydov/ms-order/internal/ports"
)

type API interface {
	PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	GetOrder(ctx context.Context, id int64) (domain.Order, error)
}

type Application struct {
	db          ports.DBPort
	paymentStub ports.PaymentPort
}

func NewApplication(db ports.DBPort, paymentStub ports.PaymentPort) *Application {
	return &Application{
		db:          db,
		paymentStub: paymentStub,
	}
}

func (a Application) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}
	// paymentErr := a.paymentStub.Charge(&order)
	// if paymentErr != nil {
	// 	st := status.Convert(paymentErr)
	// 	var allErrors []string
	// 	for _, detail := range st.Details() {
	// 		switch t := detail.(type) {
	// 		case errdetails.BadRequest:
	// 			for _, violaton := range t.GetFieldViolations() {
	// 				allErrors = append(allErrors, violaton.Description)
	// 			}
	// 		}
	// 	}
	// 	fieldErr := &errdetails.BadRequest_FieldViolation{
	// 		Field:       "payment",
	// 		Description: strings.Join(allErrors, "\n"),
	// 	}

	// 	badReq := &errdetails.BadRequest{}
	// 	badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
	// 	orderStatus := status.New(codes.InvalidArgument, "order creation failed")

	// 	statusWithDetails, _ := orderStatus.WithDetails(badReq)
	// 	return domain.Order{}, statusWithDetails.Err()
	// }
	return order, nil
}

func (a Application) GetOrder(ctx context.Context, id int64) (domain.Order, error) {
	return a.db.Get(ctx, id)
}
