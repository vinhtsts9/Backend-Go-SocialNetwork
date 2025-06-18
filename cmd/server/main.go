package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	elasticsearch "go-ecommerce-backend-api/m/v2/elasticSearch"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	"go-ecommerce-backend-api/m/v2/internal/initialize"
	"go-ecommerce-backend-api/m/v2/worker"
	websocket "go-ecommerce-backend-api/m/v2/ws"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	r := initialize.Run()

	// Khởi tạo WebSocket manager
	cm := websocket.NewConnectionManager()
	cm.Run()

	// Đăng ký route WebSocket
	r.GET("/ws", func(c *gin.Context) {
		websocket.HandleConnections(c, cm)
	})
	r.GET("/search", elasticsearch.SearchUser)

	// Đăng ký route kiểm tra trạng thái
	r.GET("/checkStatus", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// Tạo context cho toàn bộ ứng dụng
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// Sử dụng errgroup để quản lý các thành phần
	waitGroup, ctx := errgroup.WithContext(ctx)

	// Khởi động task processor
	runTaskProcessor(ctx, waitGroup, global.RedisOpt, global.Store)

	// Khởi động HTTP server
	runGinServer(ctx, waitGroup, r)

	// Đợi tất cả các goroutines hoàn thành
	err := waitGroup.Wait()
	if err != nil {
		log.Fatalf("Error from wait group: %v", err)
	}

	log.Println("Server and task processor have stopped")
}

func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	redisOpt asynq.RedisClientOpt,
	store database.Store,
) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)

	log.Println("Starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatalf("Failed to start task processor: %v", err)
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Println("Graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Println("Task processor is stopped")

		return nil
	})
}

func runGinServer(ctx context.Context, waitGroup *errgroup.Group, r *gin.Engine) {
	waitGroup.Go(func() error {
		log.Println("Starting server on :8080")
		if err := r.Run(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Println("Graceful shutdown HTTP server")
		return nil
	})
}
