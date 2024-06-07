package payment

import (
	"context"
	"log"

	"github.com/3Davydov/ms-order/internal/application/core/domain"
	"github.com/3Davydov/ms-proto/golang/payment"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient
}

var cb *gobreaker.CircuitBreaker[interface{}]

func CircuitBreakerClientInterceptor(cb *gobreaker.CircuitBreaker[interface{}]) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		_, cbErr := cb.Execute(func() (interface{}, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return nil, err
			}

			return nil, nil

		})
		return cbErr
	}
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	cb = gobreaker.NewCircuitBreaker[interface{}](gobreaker.Settings{
		Name:        "demo",
		MaxRequests: 3,
		Timeout:     4,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.1
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(CircuitBreakerClientInterceptor(cb)))
	conn, err := grpc.NewClient(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	client := payment.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	_, err := a.payment.Create(context.Background(), &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	return err
}
