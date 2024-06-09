package db

import (
	"context"
	"testing"

	"github.com/3Davydov/ms-order/config"
	"github.com/3Davydov/ms-order/internal/application/core/domain"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type OrderDatabaseTestSuite struct {
	suite.Suite
	DataSourceURL string
}

func (o *OrderDatabaseTestSuite) SetupSuite() {
	o.DataSourceURL = config.GetTestDataSourceURL()
}

func (o *OrderDatabaseTestSuite) TestSave() {
	adapter, err := NewAdapter(o.DataSourceURL)
	o.Nil(err)
	saveErr := adapter.Save(&domain.Order{})
	o.Nil(saveErr)
}

func (o *OrderDatabaseTestSuite) TestGet() {
	adapter, _ := NewAdapter(o.DataSourceURL)
	order := domain.NewOrder(2, []domain.OrderItem{
		{
			ProductCode: "CAM",
			Quantity:    5,
			UnitPrice:   1.32,
		},
	})
	adapter.Save(&order)
	ctx := context.Background()
	ord, _ := adapter.Get(ctx, order.ID)
	o.Equal(int64(2), ord.CustomerID)
}

func TestOrderDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDatabaseTestSuite))
}
