package server

import (
	"OrderServer/models"
	"OrderServer/providers"
	"OrderServer/providers/dbProvider"
	"OrderServer/providers/redisProvider"
	"OrderServer/services/Order"
	"OrderServer/utils"
	"context"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	httpServer    *http.Server
	PSQL          providers.PSQLProvider
	RedisProvider providers.RedisProvider
	OrderHandler  *Order.Handler
}

func SrvInit() *Server {
	ctx := context.Background()
	var c models.DatabaseConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
	db := dbProvider.NewPSQLProvider(c, utils.SSLModeDisable)
	redis := redisProvider.NewRedisProvider()
	orderHandler := Order.NewHandler(Order.NewService(Order.NewRepository(db.DB())))

	return &Server{
		PSQL:          db,
		RedisProvider: redis,
		OrderHandler:  orderHandler,
	}
}

func (srv *Server) Start() {
	port := os.Getenv("serverPort")

	httpSrv := &http.Server{
		Addr:    port,
		Handler: srv.InjectRoutes(),
	}
	srv.httpServer = httpSrv
	logrus.Info("Server running at PORT ", port)
	go srv.ChannelHub()
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("Start %v", err)
		return
	}
}

func (srv *Server) Stop() {
	logrus.Info("closing Postgres...")
	_ = srv.PSQL.DB().Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logrus.Info("closing server...")
	_ = srv.httpServer.Shutdown(ctx)
	logrus.Info("Done")
}
