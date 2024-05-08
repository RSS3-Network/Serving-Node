package decentralized

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/node/internal/constant"
	"github.com/rss3-network/node/internal/database"
	"github.com/rss3-network/node/internal/node/component"
	"github.com/rss3-network/node/provider/ethereum/etherface"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Component struct {
	databaseClient  database.Client
	counter         metric.Int64Counter
	etherfaceClient etherface.Client
}

const Name = "decentralized"

func (h *Component) Name() string {
	return Name
}

var _ component.Component = (*Component)(nil)

func NewComponent(_ context.Context, apiServer *echo.Echo, databaseClient database.Client) component.Component {
	c := &Component{
		databaseClient: databaseClient,
	}

	group := apiServer.Group(fmt.Sprintf("/%s", Name))

	group.GET("/tx/:id", c.GetActivity)
	group.GET("/:account", c.GetAccountActivities)
	group.GET("/count", c.GetActivitiesCount)

	if err := c.InitMeter(); err != nil {
		panic(err)
	}

	// Initialize etherface client
	etherfaceClient, err := etherface.NewEtherfaceClient()
	if err != nil {
		panic(fmt.Errorf("failed to initialize etherface client, %w", err))
	}

	c.etherfaceClient = etherfaceClient

	return c
}

func (h *Component) InitMeter() (err error) {
	meter := otel.GetMeterProvider().Meter(constant.Name)

	if h.counter, err = meter.Int64Counter(h.Name()); err != nil {
		return fmt.Errorf("failed to init meter for component %s: %w", h.Name(), err)
	}

	return nil
}

func (h *Component) CollectMetric(ctx context.Context, value string) {
	measurementOption := metric.WithAttributes(
		attribute.String("component", h.Name()),
		attribute.String("path", value),
	)

	h.counter.Add(ctx, int64(1), measurementOption)
}
