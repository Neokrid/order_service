package app

import (
	"context"
	"fmt"
	"log"

	"net/http"

	"github.com/Neokrid/order_service/common"
	configs "github.com/Neokrid/order_service/config"
	orderApp "github.com/Neokrid/order_service/internal/application"
	orderService "github.com/Neokrid/order_service/internal/domain/service/orders"
	"github.com/Neokrid/order_service/internal/entity"
	"github.com/Neokrid/order_service/internal/infrastructure/kafka"
	"github.com/Neokrid/order_service/internal/infrastructure/orders"
	internalhttp "github.com/Neokrid/order_service/internal/ports/internal_http"
	publichttp "github.com/Neokrid/order_service/internal/ports/public_http"
	database "github.com/Neokrid/order_service/pkg/database/postgres"
	"github.com/Neokrid/order_service/pkg/logger"
	"github.com/Neokrid/order_service/pkg/trx"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	cfg            *configs.Config
	dbPool         *pgxpool.Pool
	kafkaPub       *kafka.KafkaPublisher
	publicServer   *http.Server
	internalServer *http.Server
}

func NewContainer(cfg *configs.Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) Start(ctx context.Context) error {
	var tx trx.TransactionManager
	var logger logger.Logger
	pool, err := database.NewPostgresPool(ctx, c.cfg.GetDBURL())
	if err != nil {
		return fmt.Errorf("подключение к бд: %w", err)
	}
	c.dbPool = pool

	c.kafkaPub = kafka.NewKafkaPublisher(c.cfg.Kafka.Brokers, "orders.v1.events")

	repo := orders.NewRepository(c.dbPool)
	orderService := orderService.NewService(tx, logger, repo, entity.Order{})
	app := orderApp.NewOrdersService(repo, c.kafkaPub, orderService, tx)
	publicHandler := publichttp.NewPublicOrderHandler(app)
	internalHandler := internalhttp.NewPublicOrderHandler(app)

	c.publicServer = c.setupPublicServer(publicHandler)
	c.internalServer = c.setupInternalServer(internalHandler)

	go func() {
		log.Printf("Публичный Orders API запущен на порту %s", c.cfg.HTTP.Port)
		if err := c.publicServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Public server failed: %v", err)
		}
	}()

	go func() {
		log.Printf("Внутренний Orders API запущен на порту %s", c.cfg.HTTP.InternalPort)
		if err := c.internalServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Internal server failed: %v", err)
		}
	}()

	return nil
}

func (c *Container) setupPublicServer(h *publichttp.PublicOrderHandler) *http.Server {
	r := gin.Default()
	auth := common.AuthMiddleware(c.cfg.Jwt.JwtSecret)

	h.Init(r.Group(""), auth)

	return &http.Server{Addr: ":" + c.cfg.HTTP.Port, Handler: r}
}

func (c *Container) setupInternalServer(h *internalhttp.InternalOrderHandler) *http.Server {
	r := gin.Default()

	h.Init(r.Group(""))

	return &http.Server{
		Addr:    ":" + c.cfg.HTTP.InternalPort,
		Handler: r,
	}
}

func (c *Container) Stop(ctx context.Context) error {
	log.Println("Остановка order сервиса")

	if c.dbPool != nil {
		c.dbPool.Close()
	}
	c.internalServer.Shutdown(ctx)
	c.publicServer.Shutdown(ctx)
	return nil
}
