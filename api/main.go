package main

import (
	"context"
	"goauth/auth/auth_module"
	"goauth/utils/datasources"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Init database
	dataSources, _ := datasources.Init()

	// gin router with default middlewares
	router := gin.Default()
	groupV1 := router.Group("/api/v1")

	auth_module.InitAuthModule(&auth_module.AuthModuleConfig{
		DataSources: dataSources,
		Group:       groupV1,
	})
	// Initializing the server
	srv := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: router,
	}

	// ============================================================= //
	//	Graceful server shutdown
	//Start the server in a separated goroutine
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to initialze server %v\n", err)
		}
	}()

	log.Printf("Running server on %v", srv.Addr)

	//	 Wait for kill signal
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// This blocks until a signal is passed to the channel
	<-quit

	//	The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//	Shutdown the server
	log.Println("Shutting down the server... ")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: \n%v\n", err)
	}
}
