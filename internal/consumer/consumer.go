// Package consumer
package consumer

import (
	"context"
	"encoding/json"
	"product-catalog-service/internal/entity"
	"product-catalog-service/internal/service"
	"strings"

	zlog "github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	productSvc *service.ProductService
}

func NewConsumer(productSvc *service.ProductService) *Consumer {
	return &Consumer{
		productSvc: productSvc,
	}
}

func (c *Consumer) StartKafkaConsumer() {
	// create kafka reader for order topic
	orderReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topic:    "order-topic",
		GroupID:  "product-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	for {
		// read message from order topic
		ctx := context.Background()
		msg, err := orderReader.ReadMessage(ctx)
		if err != nil {
			//logger.Error().Err(err).Msgf("error getting  product %d from cache", productID)
			zlog.Error().Err(err).Msgf("error reading message")

			continue
		}

		// Process Message
		c.processMessage(ctx, msg)

	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) {
	// unmarshal the message payload
	var orderevent entity.Order

	err := json.Unmarshal(msg.Value, &orderevent)
	if err != nil {
		zlog.Error().Msgf("error unmarshalling message: %v", err)
		return
	}

	// key -> "order.created.orderID" or "order.cancelled.orderID"
	key := string(msg.Key)
	listKey := strings.Split(key, ".")
	eventType := listKey[1]

	// process the order event bas on status
	switch eventType {
	case "created":
		// process order created event
		for _, item := range orderevent.Productrequests {
			err := c.productSvc.ReserveProductStock(ctx, item.ProductID, item.Quantity)
			if err != nil {

				zlog.Error().Msgf("error (created) updating stock for product %d: %v", item.ProductID, err)
			}
		}
	case "cancelled":
		// process order created event
		for _, item := range orderevent.Productrequests {
			err := c.productSvc.ReleaseProductStock(ctx, item.ProductID, item.Quantity)
			if err != nil {
				zlog.Error().Msgf("error (cancelled) updating stock for product %d: %v", item.ProductID, err)

			}
		}
	default:

		zlog.Error().Msgf("unknown order status: %d", orderevent.Status)

	}
}
