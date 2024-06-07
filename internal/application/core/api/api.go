package api

import (
	"strings"

	"github.com/3Davydov/ms-order/internal/application/core/domain"
	"github.com/3Davydov/ms-order/internal/ports"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	PlaceOrder(order domain.Order) (domain.Order, error)
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

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}
	paymentErr := a.paymentStub.Charge(&order)
	if paymentErr != nil {
		st := status.Convert(paymentErr)
		var allErrors []string
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case errdetails.BadRequest:
				for _, violaton := range t.GetFieldViolations() {
					allErrors = append(allErrors, violaton.Description)
				}
			}
		}
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: strings.Join(allErrors, "\n"),
		}

		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")

		statusWithDetails, _ := orderStatus.WithDetails(badReq)
		return domain.Order{}, statusWithDetails.Err()
	}
	return order, nil
}
